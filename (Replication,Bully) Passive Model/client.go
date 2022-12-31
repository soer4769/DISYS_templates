// --------------------------- //
// ---------- IMPORT --------- //
// --------------------------- //
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
    "strconv"

	"github.com/GoPrm/proto"
	"google.golang.org/grpc"
)

// --------------------------- //
// --------- GLOBALS --------- //
// --------------------------- //
type client struct {
	id        int32
	primary   GoPrm.PrmClient
}

// --------------------------- //
// --------- METHODS --------- //
// --------------------------- //
func (c *client) Frontend(serverCount int) {
    if c.primary == nil {
        log.Println("Searching for primary server...")
        SearchId := serverCount-1
        
        for {
            conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", SearchId+5000), grpc.WithTimeout(time.Second*2), grpc.WithBlock(), grpc.WithInsecure())
            if err != nil {
                SearchId--
                log.Printf("Failed to connect to %v: %v", SearchId, err)
                continue
            }
            
            log.Printf("Client connected to server %v", SearchId)
            server := GoPrm.NewPrmClient(conn)
            
            if c.primary != nil {
                for {
                    if primary, _ := server.GetPrimary(context.Background(), &GoPrm.Empty{}); int(primary.Id) != -1 {
                        SearchId = int(primary.Id)
                        break
                    }
                    time.Sleep(time.Second*2)
                }
                continue
            }
            
            c.primary = server
            log.Printf("Found primary server %v", SearchId)
            return
        }
    }
}

// --------------------------- //
// ---------- SETUP ---------- //
// --------------------------- //
func main() {
	aid, _ := strconv.ParseInt(os.Args[1], 10, 32) // client id
    sc, _ := strconv.ParseInt(os.Args[2], 10, 32) // server count
    
    c := &client {
        id: int32(aid),
        primary: nil,
    }
    
    for {
        c.Frontend(int(sc))
        time.Sleep(time.Second*5)
        if _, err := c.primary.AddVal(context.Background(),&GoPrm.Task{Id: c.id, Query: 1, Lamport: -1}); err == nil {
            log.Printf("Failed to connect to server %v, finding new primary server...", c.primary)
            c.primary = nil
            continue
        }
        log.Printf("Increasing number by %v", 1)
    }
}
