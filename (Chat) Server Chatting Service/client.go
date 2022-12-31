// ---------------------------- //
// ---------- IMPORT ---------- //
// ---------------------------- //
package main

import (
    "context"
    "io"
    "log"
    "os"
    "bufio"
    "strings"
    "strconv"

    "github.com/GoChat/proto"
    "google.golang.org/grpc"
)

// ----------------------------- //
// ---------- GLOBALS ---------- //
// ----------------------------- //
type client struct {
    id      int32
    lamport int64
    chat    GoChat.ChatClient
    done    chan bool
}

// ---------------------------- //
// ---------- METHODS --------- //
// ---------------------------- //
func (c *client) inPost(inStream GoChat.Chat_ConnectClient) {
    for {
        resp, err := inStream.Recv()
        if err == io.EOF {
            c.done <- true
            return
        }

        if err != nil {
            log.Fatalf("Failed to receive %v", err)
        }

        if c.id < 0 {
            c.id = resp.Id
        }

        if c.id != resp.Id {
            if resp.Lamport > c.lamport {
                c.lamport = resp.Lamport + 1
            } else {
                c.lamport++
            }
        }

        log.Println(strings.Replace(resp.Message, "#L#", strconv.FormatInt(c.lamport,10), 1))
    }
}

func (c *client) outPost(outStream GoChat.Chat_MessagesClient) {
    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        textIn := scanner.Text()
        if len(textIn) > 128 {
            log.Println("Failed to send message: more than 128 characters")
            continue
        }

        c.lamport++
        usrIn := GoChat.Post{Id: c.id, Message: textIn, Lamport: c.lamport}
        outStream.Send(&usrIn)

        if textIn == "--exit" {
            c.chat.Disconnect(context.Background(), &usrIn)
            c.done <- true
            return
        }
    }

    if err := scanner.Err(); err != nil {
       log.Printf("Failed to scan input: %v", err)
    }
}

// -------------------------- //
// ---------- SETUP---------- //
// -------------------------- //
func main() {
    // dial server
    conn, err := grpc.Dial("localhost:5000", grpc.WithInsecure())
    log.Printf("Client connecting to server...")
    if err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
    c := &client {
        id: -1,
        lamport: 0,
        chat: GoChat.NewChatClient(conn),
        done: make(chan bool),
    }
    defer log.Printf("Closing server connection...\nServer connection ended...")
    defer conn.Close()

    // setup streams
    inStream, inErr := c.chat.Connect(context.Background(), &GoChat.Post{Lamport: c.lamport})
    if inErr != nil {
        log.Fatalf("Failed to open connection stream: %v; %v", inErr)
    }
    outStream, outErr := c.chat.Messages(context.Background())
    if outErr != nil {
        log.Fatalf("Failed to open message stream: %v", outErr)
    }
    log.Println("Client connected successfully...")

    // running goroutines streams
    go c.inPost(inStream)
    go c.outPost(outStream)
    <-c.done
}
