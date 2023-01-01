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

	"github.com/GoMECS/proto"
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
    ANSWER
)

type client struct {
    GoMECS.UnimplementedMECSServer
    id          int32
    lamport     int64
    state       int
    votes       int
    coordID     int32
    coordAlive  bool
    occupiedId  int32
    doneCrit    bool
    inCrit      bool
    queue       []int32
    peers       map[int32]GoMECS.MECSClient
}

// ---------------------------- //
// ---------- METHODS --------- //
// ---------------------------- //
func (c *client) Request(context context.Context, query *GoMECS.Query) (*GoMECS.Empty, error) {
    log.Printf("Peer %v requested access to critical section",query.Id)
    c.queue = append(c.queue,query.Id)
    for {
        if c.occupiedId == -1 && c.queue[len(c.queue)-1] == query.Id {
            log.Printf("Granted peer %v access to critical section",query.Id)
            c.occupiedId = query.Id
            return &GoMECS.Empty{},nil
        }
    }
}

func (c *client) Release(context context.Context, query *GoMECS.Query) (*GoMECS.Empty, error) {
    log.Printf("occupied: %v, query: %v",c.occupiedId,query.Id)
    if c.occupiedId == query.Id {
        log.Printf("Peer %v released their access to critical section",query.Id)
        log.Println(c.queue)
        c.queue = c.queue[1:]
        log.Println(c.queue)
        c.occupiedId = -1
    }
    return &GoMECS.Empty{},nil
}

func (c *client) RunCritical(coordId int32) {
    if !c.inCrit {
        log.Printf("Requested access to critical section...")
        c.inCrit = true
        for {
            if !c.doneCrit {
                if _, err := c.peers[coordId].Request(context.Background(),&GoMECS.Query{Id: c.id}); err == nil {
                    c.lamport++
                    log.Println("Inside of Critical section...")
                    time.Sleep(5 * time.Second)
                    log.Println("Outside of critical section...")
                    c.doneCrit = true
                    c.peers[coordId].Release(context.Background(),&GoMECS.Query{Id: c.id})
                }
            }
        }
    }
}

// ---------------------------- //
// ----------- BULLY ---------- //
// ---------------------------- //
func (c *client) Send(id int32, query int) {
    c.lamport++
    msg := GoMECS.Task{Id: c.id, Query: int32(query)}
    log.Printf("Sendt query %v to %v (lamport %v)", query, id, c.lamport)
    if _, err := c.peers[id].Elect(context.Background(), &msg, grpc.WaitForReady(true)); err != nil {
        log.Fatalf("Failed to send %v a message: %v", id, err)
        delete(c.peers,id)
    }
}

func (c *client) Broadcast(query int, higher bool) {
    for i, _ := range c.peers {
        if higher && i > c.id || !higher {
            c.Send(i,query)
        }
    }
}

func (c *client) Elect(context context.Context, ans *GoMECS.Task) (*GoMECS.Empty, error) {
    if ans.Lamport > c.lamport {
        c.lamport = ans.Lamport + 1
    } else {
        c.lamport++
    }
    
    log.Printf("Recieved answer %v by %d (lamport %v)", ans.Query, ans.Id, c.lamport)
    
    switch ans.Query {
        case ANSWER:
            c.votes++
        case COORDINATOR:
            if c.state == CANDIDATE || c.state == COORDINATOR {
                log.Printf("Status: follower")
            }
            if c.coordID != ans.Id {
                log.Printf("Follower of new coordinator %v",ans.Id)
            }
            c.coordID = ans.Id
            c.state = FOLLOWER
            c.coordAlive = true
        case ELECTION:
            c.Send(ans.Id,ANSWER)
    }
    return &GoMECS.Empty{}, nil
}

func (c *client) Bully() {
    sendElection := false
    for {
        switch c.state {
            case CANDIDATE:
                if len(c.peers) == int(c.id) || sendElection && c.votes == 0 {
                    log.Printf("Status: coordinator")
                    c.Broadcast(COORDINATOR,false)
                    c.state = COORDINATOR
                    //c.CloneLog(context.Background(),&GoMECS.Log{Log: c.log})
                } else if sendElection && c.votes > 0 {
                    time.Sleep(time.Second*2)
                    if c.state != FOLLOWER {
                        c.votes = 0
                    }
                } else if !sendElection {
                    log.Printf("Status: candidate")
                    sendElection = true
                    c.Broadcast(ELECTION,true)
                    time.Sleep(time.Second*2)
                }
            case FOLLOWER:
                c.RunCritical(c.coordID)
                c.coordAlive = false
                time.Sleep(time.Second*8)
                if !c.coordAlive {
                    log.Printf("Coordinator not responding, starts new election...")
                    delete(c.peers,c.coordID)
                    c.state = CANDIDATE
                    c.coordID = -1
                    c.votes = 0
                }
            case COORDINATOR:
                log.Printf("Broadcasts heartbeat")
                c.Broadcast(COORDINATOR,false)
                time.Sleep(time.Second*5)
        }
    }
}

// -------------------------- //
// ---------- SETUP---------- //
// -------------------------- //
func main() {
	pid, _ := strconv.ParseInt(os.Args[1], 10, 32) // peer id
    pcount, _ := strconv.ParseInt(os.Args[2], 10, 32) // peer count
    
    c := &client {
        id: int32(pid),
        lamport: 0,
        state: CANDIDATE,
        votes: 0,
        coordID: -1,
        coordAlive: false,
        occupiedId: -1,
        peers: make(map[int32]GoMECS.MECSClient),
    }
    
    // Start Server
    lis, err := net.Listen("tcp",fmt.Sprintf("localhost:%d",int(pid)+5000))
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }
    
    s := grpc.NewServer()
    GoMECS.RegisterMECSServer(s, c)
    log.Printf("Server listening at %v", lis.Addr())
    
    go func() {
        if err := s.Serve(lis); err != nil {
            log.Fatalf("Failed to serve: %v", err)
        }
    }()
    
    // Connect Peers
    for i := 0; i < int(pcount); i++ {
        if i != int(pid) {
            conn, err := grpc.Dial(fmt.Sprintf("localhost:%d",i+5000), grpc.WithInsecure())
            if err != nil {
                log.Fatalf("Failed to connect %v", err)
            }

            c.peers[int32(i)] = GoMECS.NewMECSClient(conn)
            log.Printf("Client connected to peer %v", i)
            defer conn.Close()
        }
    }
    
    // Run Algorithm
    c.Bully()
}
