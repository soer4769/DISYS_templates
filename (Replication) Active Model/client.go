// --------------------------- //
// ---------- IMPORT --------- //
// --------------------------- //
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
    "math/rand"

	"github.com/GoArm/proto"
	"google.golang.org/grpc"
)

// --------------------------- //
// --------- METHODS --------- //
// --------------------------- //
type client struct {
	id          int32
	number      int32
	servers     map[int]GoArm.ArmClient
}

func (c *client) BroadcastNum(number int32) {
    log.Printf("Broadcasting number %d", c.number)
    
    for i := 0; i < len(c.servers); i++ {
        if _, err := c.servers[i].SetNum(context.Background(), &GoArm.In{Id: c.id, Number: c.number}); err != nil {
            log.Printf("Server %v unresponsive, disconnected...", i)
            delete(c.servers,i)
            i -= 1
        }
    }
}

func (c *client) GetNum() int32 {
	timeout, _ := context.WithTimeout(context.Background(), time.Second*2)
    
    for i := 0; i < len(c.servers); i++ {
        if result, err := c.servers[i].GetNum(timeout, &GoArm.Empty{}); err == nil {
            return result.Number
        }
    }
    
    return 0
}

// --------------------------- //
// ---------- SETUP ---------- //
// --------------------------- //
func main() {
    rand.Seed(time.Now().UnixNano())
	aid, _ := strconv.ParseInt(os.Args[1], 10, 32) // client id
    sc, _ := strconv.ParseInt(os.Args[2], 10, 32) // server count
    
    c := &client {
        id: int32(aid),
        number: 0,
        servers: make(map[int]GoArm.ArmClient, int(sc)),
    }
    
	for i := 0; i < int(sc); i++ {
        conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", i+5000), grpc.WithInsecure())
        if err != nil {
            log.Fatalf("Failed to connect to %v: %v", i, err)
        }

        c.servers[i] = GoArm.NewArmClient(conn)
        log.Printf("Client connected to server %v", i)
        defer conn.Close()
	}
	
    for {
        c.number = c.GetNum()+1
        c.BroadcastNum(c.number)
        time.Sleep(time.Second*time.Duration(rand.Intn(2)+1))
	}
}
