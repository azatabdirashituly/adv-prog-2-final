package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type Client struct {
	conn net.Conn
	name string
}

type Message struct {
	Name    string
	Content string
	Time    string
}

var (
	clients     []*Client
	clientsLock sync.Mutex
	messageLog    []Message
	messageLock sync.Mutex
)

func main() {
	fmt.Println("Launching server...")
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer ln.Close()

	go storeMessageLog()
	
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		client := &Client{conn: conn}

		clientsLock.Lock()
		clients = append(clients, client)
		clientsLock.Unlock()

		go handleConnection(client)
	}
}

func handleConnection(client *Client) {
	defer func() {
		client.conn.Close()

		clientsLock.Lock()
		defer clientsLock.Unlock()

		for i, c := range clients {
			if c == client {
				clients = append(clients[:i], clients[i+1:]...)
				break
			}
		}
	}()

	reader := bufio.NewReader(client.conn)
	fmt.Fprint(client.conn, "Enter your name: ")
	name, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading name:", err)
		return
	}
	client.name = strings.TrimSpace(name)
	broadcast(fmt.Sprintf("%s has joined the chat. ", client.name))

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading:", err)
			return
		}
		message = strings.TrimSpace(message)

		currentTime := time.Now().Format("15:04:05")

		modifiedMessage := fmt.Sprintf("%s: %s %s\n", client.name, message, currentTime)
        broadcast(modifiedMessage)

		historyMessage := Message{Name: client.name, Content: message, Time: currentTime}
		addToMessageLog(historyMessage)
		
		messageLock.Lock()
		messageLog = append(messageLog, historyMessage)
		messageLock.Unlock()
	}
}


func broadcast(message string) {
	clientsLock.Lock()
	defer clientsLock.Unlock()

	message = "\n" + message

	for _, client := range clients {
		_, err := client.conn.Write([]byte(message + "\n"))
		if err != nil {
			fmt.Println("Error writing to client:", err)
			continue
		}
	}
}

func addToMessageLog(message Message) {
	messageLock.Lock()
	defer messageLock.Unlock()

	messageLog = append(messageLog, message)
}

func storeMessageLog() {
    for {
        time.Sleep(5 * time.Second)
        messageLock.Lock()

        file, err := os.OpenFile("message_log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
        if err != nil {
            fmt.Println("Error opening file:", err)
            messageLock.Unlock()
            continue
        }

        writer := bufio.NewWriter(file)
        for _, msg := range messageLog {
            _, err := fmt.Fprintf(writer, "[%s] %s: %s\n", msg.Time, msg.Name, msg.Content)
            if err != nil {
                fmt.Println("Error writing to file:", err)
                break
            }
        }
        writer.Flush()
        file.Close()

        messageLock.Unlock()
    }
}


