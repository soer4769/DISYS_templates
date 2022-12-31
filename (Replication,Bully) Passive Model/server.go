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

	"github.com/GoPrm/proto"
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

type server struct {
    GoPrm.UnimplementedPrmServer
    id          int32
    lamport     int64
    state       int
    votes       int
    coordID     int32
    coordAlive  bool
    number      int32
    log         map[int]int32
    peers       map[int32]GoPrm.PrmClient
}

// ---------------------------- //
// ---------- METHODS --------- //
// ---------------------------- //
func (s *server) GetPrimary(context context.Context, _ *GoPrm.Empty) (*GoPrm.Primary, error) {
    return &GoPrm.Primary{Id: s.coordID}, nil
}

func (s *server) SetVal(context context.Context, task *GoPrm.Task) (*GoPrm.Empty, error) {
    log.Printf("Set number to value %v", task.Query)
    s.log[len(s.log)] = task.Query
    s.number = task.Query
    if s.state == COORDINATOR {
        log.Println("Broadcasting value to backup servers...")
        for _, p := range s.peers {
            p.SetVal(context, task)
        }
    }
    return &GoPrm.Empty{},nil
}

func (s *server) GetVal(context context.Context, _ *GoPrm.Empty) (*GoPrm.Value, error) {
    return &GoPrm.Value{Query: s.number},nil
}

// ---------------------------- //
// ----------- BULLY ---------- //
// ---------------------------- //
func (s *server) Send(id int32, query int) {
    s.lamport++
    msg := GoPrm.Task{Id: s.id, Query: int32(query)}
    log.Printf("Sendt query %v to %v (lamport %v)", query, id, s.lamport)
    if _, err := s.peers[id].Elect(context.Background(), &msg, grpc.WaitForReady(true)); err != nil {
        log.Fatalf("Failed to send %v a message: %v", id, err)
        delete(s.peers,id)
    }
}

func (s *server) Broadcast(query int, higher bool) {
    for i, _ := range s.peers {
        if higher && i > s.id || !higher {
            s.Send(i,query)
        }
    }
}

func (s *server) Elect(context context.Context, ans *GoPrm.Task) (*GoPrm.Empty, error) {
    if ans.Lamport > s.lamport {
        s.lamport = ans.Lamport + 1
    } else {
        s.lamport++
    }
    
    log.Printf("Recieved answer %v by %d (lamport %v)", ans.Query, ans.Id, s.lamport)
    
    switch ans.Query {
        case ANSWER:
            s.votes++
        case COORDINATOR:
            if s.state == CANDIDATE || s.state == COORDINATOR {
                log.Printf("Status: follower")
            }
            s.coordID = ans.Id
            s.state = FOLLOWER
            s.coordAlive = true
        case ELECTION:
            s.Send(ans.Id,ANSWER)
    }
    return &GoPrm.Empty{}, nil
}

func (s *server) Bully() {
    sendElection := false
    for {
        switch s.state {
            case CANDIDATE:
                if len(s.peers) == int(s.id) || sendElection && s.votes == 0 {
                    log.Printf("Status: coordinator")
                    s.Broadcast(COORDINATOR,false)
                    s.state = COORDINATOR
                } else if sendElection && s.votes > 0 {
                    time.Sleep(time.Second*2)
                    if s.state != FOLLOWER {
                        s.votes = 0
                    }
                } else if !sendElection {
                    log.Printf("Status: candidate")
                    sendElection = true
                    s.Broadcast(ELECTION,true)
                    time.Sleep(time.Second*2)
                }
            case FOLLOWER:
                s.coordAlive = false
                time.Sleep(time.Second*8)
                if !s.coordAlive {
                    log.Printf("Coordinator not responding, starts new election...")
                    delete(s.peers,s.coordID)
                    s.state = CANDIDATE
                    s.coordID = -1
                    s.votes = 0
                }
            case COORDINATOR:
                log.Printf("Broadcasts heartbeat")
                s.Broadcast(COORDINATOR,false)
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
    
    server := &server {
        id: int32(pid),
        lamport: 0,
        state: CANDIDATE,
        votes: 0,
        coordID: -1,
        coordAlive: false,
        number: 0,
        log: make(map[int]int32),
        peers: make(map[int32]GoPrm.PrmClient),
    }
    
    // Start Server
    lis, err := net.Listen("tcp",fmt.Sprintf("localhost:%d",int(pid)+5000))
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }
    
    s := grpc.NewServer()
    GoPrm.RegisterPrmServer(s, server)
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
            
            server.peers[int32(i)] = GoPrm.NewPrmClient(conn)
            log.Printf("Client connected to peer %v", i)
            defer conn.Close()
        }
    }
    
    server.Bully()
}
