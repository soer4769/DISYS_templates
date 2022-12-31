// --------------------------- //
// ---------- IMPORT --------- //
// --------------------------- //
package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/GoArm/proto"
	"google.golang.org/grpc"
)

// --------------------------- //
// --------- METHODS --------- //
// --------------------------- //
type server struct {
	GoArm.UnimplementedArmServer
	number int32
}

func (s *server) SetNum(context context.Context, in *GoArm.In) (*GoArm.Empty, error) {
    if in.Number > s.number {
        log.Printf("Client %v sat number to %v", in.Id, in.Number)
        s.number = in.Number
    }
    return &GoArm.Empty{}, nil
}

func (s *server) GetNum(context context.Context, empty *GoArm.Empty) (*GoArm.Out, error) {
    log.Printf("Sending number %v...", s.number)
    return &GoArm.Out{Number: s.number}, nil
}

// --------------------------- //
// ---------- SETUP ---------- //
// --------------------------- //
func main() {
	pid, _ := strconv.ParseInt(os.Args[1], 10, 32) // server id

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", int(pid)+5000))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	GoArm.RegisterArmServer(s, &server{})
	log.Printf("Server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
