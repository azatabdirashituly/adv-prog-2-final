package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

type Client struct {
	conn net.Conn
	name string
}

var (
	clients     []*Client
	clientsLock sync.Mutex
)

func main() {
	fmt.Println("Launching server...")
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer ln.Close()

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
		fmt.Printf("[%s]: %s\n", client.name, message)

		currentTime := time.Now().Format("15:04:05")
		fmt.Println("Received at:", currentTime)

		broadcast(fmt.Sprintf("%s: %s %s\n", client.name, message, currentTime))
	}
}

func broadcast(message string) {
	clientsLock.Lock()
	defer clientsLock.Unlock()

	for _, client := range clients {
		_, err := client.conn.Write([]byte(message + "\n"))
		if err != nil {
			fmt.Println("Error writing to client:", err)
			continue
		}
	}
}

