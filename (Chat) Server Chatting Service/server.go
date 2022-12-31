// ---------------------------- //
// ---------- IMPORT ---------- //
// ---------------------------- //
package main

import (
    "context"
    "log"
    "net"
    "io"
    "fmt"

	"github.com/GoChat/proto"
	"google.golang.org/grpc"
)

// ----------------------------- //
// ---------- GLOBALS ---------- //
// ----------------------------- //
type server struct {
    GoChat.UnimplementedChatServer
    users    map[int32]GoChat.Chat_ConnectServer
    usrCount int32
    usrDis   int32
}

// ---------------------------- //
// ---------- METHODS --------- //
// ---------------------------- //
func (s *server) Broadcast(msg *GoChat.Post) {
    for i, u := range s.users {
        if int32(i) != msg.Id {
            u.Send(msg)
        }
    }
}

func (s *server) Connect(in *GoChat.Post, srv GoChat.Chat_ConnectServer) error {
    s.usrCount++
    usrId := s.usrCount
    userMsg := fmt.Sprintf("Participant %v joined Chitty-Chat at Lamport time #L#", usrId)
    s.users[usrId] = srv

    s.Broadcast(&GoChat.Post{Id: usrId, Lamport: in.Lamport, Message: userMsg})
    srv.Send(&GoChat.Post{Id: usrId, Lamport: in.Lamport, Message: userMsg})
    srv.Send(&GoChat.Post{Message: "You joined the chat!\n--exit            Leave Server"})
    log.Printf("Client %v connected to the server...", usrId)
    for {
        if s.usrDis == usrId {
            return nil
        }
    }
}

func (s *server) Disconnect(context context.Context, in *GoChat.Post) (*GoChat.Empty, error) {
    log.Printf("Client %v disconnected from the server...", in.Id)
    s.Broadcast(&GoChat.Post{Id: in.Id, Lamport: in.Lamport, Message: fmt.Sprintf("Participant %v left Chitty-Chat at Lamport time #L#", in.Id)})
    delete(s.users,in.Id)
    s.usrDis = in.Id
    return &GoChat.Empty{}, nil
}

func (s *server) Messages(srv GoChat.Chat_MessagesServer) error {
    for {
        resp, err := srv.Recv()
        if err == io.EOF || err != nil || resp.Message == "--exit" {
            return nil
        }

        log.Printf("Client %v wrote the message: %v", resp.Id, resp.Message)
        s.Broadcast(&GoChat.Post{Id: resp.Id, Lamport: resp.Lamport, Message: fmt.Sprintf("Participant %v wrote at Lamport time #L#: %v", resp.Id, resp.Message)})
    }
}

// -------------------------- //
// ---------- SETUP---------- //
// -------------------------- //
func main() {
    lis, err := net.Listen("tcp", "localhost:5000")
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    s := grpc.NewServer()
	GoChat.RegisterChatServer(s, &server {
        users: make(map[int32]GoChat.Chat_ConnectServer),
        usrCount: 0,
        usrDis: -1,
    })
	log.Printf("Server listening at %v", lis.Addr())

    if err := s.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}
