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
func (c *client) Frontend(serverCount *int64) {
    if c.primary == nil {
        log.Println("Searching for new primary server...")
        SearchId := int(*serverCount-1)
        
        for {
            conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", SearchId+5000), grpc.WithTimeout(time.Second*5), grpc.WithBlock(), grpc.WithInsecure())
            if err != nil {
                log.Printf("Failed to connect to %v: %v", SearchId, err)
                SearchId--
                continue
            }
            
            log.Printf("Client connected to server %v", SearchId)
            server := GoPrm.NewPrmClient(conn)
            
            if c.primary != nil {
                for {
                    if primary, _ := server.GetPrimary(context.Background(), &GoPrm.Empty{}); int(primary.Id) != -1 {
                        SearchId = int(primary.Id)
                        conn.Close()
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
        c.Frontend(&sc)
        time.Sleep(time.Second*5)
        number, err := c.primary.GetVal(context.Background(),&GoPrm.Empty{})
        if err != nil {
            sc--
            log.Printf("Server connection error: %v", err)
            c.primary = nil
            continue
        }
        c.primary.SetVal(context.Background(),&GoPrm.Task{Id: c.id, Query: number.Query+1, Lamport: -1})
        log.Printf("Increasing number by %v to %v", 1, number.Query)
    }
}
