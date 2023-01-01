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

	"github.com/GoMERA/proto"
	"google.golang.org/grpc"
)

// ----------------------------- //
// ---------- GLOBALS ---------- //
// ----------------------------- //
const (
    REQUEST = iota
    ALLOWED
    RELEASED
    WANTED
    HELD
)

type client struct {
    GoMERA.UnimplementedMERAServer
    id        int32
    state     int
    lamport   int64
    peers     map[int]GoMERA.MERAClient
    reqQueue  []GoMERA.Post
    replies   int
}

// ---------------------------- //
// ---------- METHODS --------- //
// ---------------------------- //
func (c *client) Send(id int, request int) {
    c.lamport++
    msg := GoMERA.Post{Id: c.id, Request: int32(request), Lamport: c.lamport}
    log.Printf("Sendt message %v to %v (lamport %v)", request, id, c.lamport)
    if _, err := c.peers[id].Recv(context.Background(), &msg, grpc.WaitForReady(true)); err != nil {
        log.Fatalf("Failed to send %v a message: %v", id, err)
    }
}

func (c *client) Broadcast(request int) {
    for i, _ := range c.peers {
        c.Send(i,request)
    }
}

func (c *client) Recv(context context.Context, resp *GoMERA.Post) (*GoMERA.Empty, error) {
    if resp.Lamport > c.lamport {
        c.lamport = resp.Lamport + 1
    } else {
        c.lamport++
    }
    
    log.Printf("Recieved message %v by %d (lamport %v)", resp.Request, resp.Id, c.lamport)
    
    if resp.Request == REQUEST {
        if c.state == HELD || (c.state == WANTED && c.id < resp.Id) {
            c.reqQueue = append(c.reqQueue,*resp)
        } else {
            c.Send(int(resp.Id),ALLOWED)
        }
    } else {
        c.replies++
    }
    
    return &GoMERA.Empty{}, nil
}

// -------------------------- //
// ---------- SETUP---------- //
// -------------------------- //
func main() {
	pid, _ := strconv.ParseInt(os.Args[1], 10, 32) // peer id
    pcount, _ := strconv.ParseInt(os.Args[2], 10, 32) // peer count
    
    c := &client {
        id: int32(pid),
        state: RELEASED,
        lamport: 0,
        peers: make(map[int]GoMERA.MERAClient,int(pcount)),
        reqQueue: make([]GoMERA.Post,0),
        replies: 0,
    }
    
    // Start Server
    lis, err := net.Listen("tcp",fmt.Sprintf("localhost:%d",int(c.id)+5000))
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }
    
    server := grpc.NewServer()
    GoMERA.RegisterMERAServer(server, c)
    log.Printf("Server listening at %v", lis.Addr())
    
    go func() {
        if err := server.Serve(lis); err != nil {
            log.Fatalf("Failed to serve: %v", err)
        }
    }()
    
    // Connect Peers
    for i := 0; i < int(pcount); i++ {
        if i != int(c.id) {
            conn, err := grpc.Dial(fmt.Sprintf("localhost:%d",i+5000), grpc.WithInsecure())
            if err != nil {
                log.Fatalf("Failed to connect %v", err)
            }

            c.peers[i] = GoMERA.NewMERAClient(conn)
            log.Printf("Client connected to peer %v", i)
            defer conn.Close()
        }
    }
    
    // Run Algorithm
    doneCritical := false
    for {
        if !doneCritical {
            if c.state != WANTED {
                c.state = WANTED
                c.Broadcast(REQUEST)
            } else if c.replies == int(pcount-1) {
                c.lamport++
                c.state = HELD
                log.Println("Inside of Critical section")
                time.Sleep(5 * time.Second)
                log.Println("Outside of critical section")
                c.state = RELEASED
                doneCritical = true
                for i, req := range c.reqQueue {
                    Send(req.Id,ALLOWED)
                }
                c.reqQueue = nil
            }
        }
    }
}
