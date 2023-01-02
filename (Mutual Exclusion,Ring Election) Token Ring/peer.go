// ---------------------------- //
// ---------- IMPORT ---------- //
// ---------------------------- //
package main

import (
	"context"
	"log"
	"net"
	"os"
    "fmt"
    "time"
	"strconv"

	"github.com/GoMETR/proto"
	"google.golang.org/grpc"
)

// ----------------------------- //
// ---------- GLOBALS ---------- //
// ----------------------------- //
const (
    CANDIDATE = iota
    COORDINATOR
    FOLLOWER
    ELECTION
    ELECTED
)

type client struct {
    GoMETR.UnimplementedMETRServer
    id          int32
    lamport     int64
    state       int
    coordId     int32
    neighId     int32
    begunElect  bool
    neighbor    GoMETR.METRClient
}

// ---------------------------- //
// ---------- METHODS --------- //
// ---------------------------- //
func (c *client) Send(id int32, electstat int32) {
    c.lamport++
    msg := GoMETR.Token{Id: id, Electstat: electstat, Lamport: c.lamport}
    log.Printf("Sendt message '%v' to '%v' (lamport %v)", electstat, c.neighId, c.lamport)
    if _, err := c.neighbor.PassToken(context.Background(), &msg, grpc.WaitForReady(true)); err != nil {
        log.Fatalf("Failed to send %v a message: %v", id, err)
    }
}

func (c *client) PassToken(context context.Context, token *GoMETR.Token) (*GoMETR.Empty, error) {
    if token.Lamport > c.lamport {
        c.lamport = token.Lamport + 1
    } else {
        c.lamport++
    }
    c.Ring(token)
    return &GoMETR.Empty{}, nil
}

func (c *client) Ring(token *GoMETR.Token) {
    if c.state == CANDIDATE {
        switch token.Electstat {
            case ELECTION:
                if token.Id > c.id {
                    log.Printf("Recieved higher id message %v, forwarding...",token.Id)
                    c.coordId = token.Id
                    c.Send(token.Id, token.Electstat)
                } else if token.Id < c.id && !c.begunElect {
                    log.Printf("Recieved <id election message %v without starting election, forwarding own...",token.Id)
                    c.begunElect = true
                    c.Send(c.id, token.Electstat)
                } else if token.Id == c.id {
                    log.Println("Recieved own election message...")
                    c.Send(c.id, ELECTED)
                }
            case ELECTED:
                if token.Id != c.id {
                    log.Printf("Recieved elected message %v, now follower", token.Id)
                    c.state = FOLLOWER
                    c.Send(token.Id, token.Electstat)
                } else if token.Id == c.id {
                    log.Println("Recieved own elected message, now coordinator")
                    c.state = COORDINATOR
                    c.Send(c.id, -1)
                }
        }
    } else {
        log.Printf("Recieved token %v from neighbor, forwading...",token.Electstat)
        time.Sleep(time.Second*2)
        c.Send(c.id, token.Electstat)
    }
}

// -------------------------- //
// ---------- SETUP---------- //
// -------------------------- //
func main() {
	pid, _ := strconv.ParseInt(os.Args[1], 10, 32) // peer id
    pcount, _ := strconv.ParseInt(os.Args[2], 10, 32) // peer count
    pneigh := (int(pid)+1)%int(pcount)
    
    c := &client {
        id: int32(pid),
        lamport: 0,
        state: CANDIDATE,
        neighId: int32(pneigh),
        coordId: -1,
    }
    
    // Start Server
    lis, err := net.Listen("tcp",fmt.Sprintf("localhost:%d",int(pid)+5000))
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }
    
    server := grpc.NewServer()
    GoMETR.RegisterMETRServer(server, c)
    log.Printf("Server listening at %v", lis.Addr())
    
    go func() {
        if err := server.Serve(lis); err != nil {
            log.Fatalf("Failed to serve: %v", err)
        }
    }()
    
    // Connect Neighbor
    log.Println(pneigh)
    conn, err := grpc.Dial(fmt.Sprintf("localhost:%d",pneigh+5000), grpc.WithInsecure())
    if err != nil {
        log.Fatalf("Failed to connect %v", err)
    }

    c.neighbor = GoMETR.NewMETRClient(conn)
    log.Printf("Client connected to peer %v", pneigh)
    defer conn.Close()
    
    // Ring Token
    for {
        if !c.begunElect {
            c.begunElect = true
            log.Println("Begins new election...")
            c.Send(int32(pid), ELECTION)
        }
    }
}
