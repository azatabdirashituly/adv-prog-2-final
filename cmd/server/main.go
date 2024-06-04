package main

import (
	"fmt"
	"net"
	"tcp-chat/internal/chat"
	"tcp-chat/internal/core"
	"tcp-chat/internal/handlers"
)

func main() {
	fmt.Println("Launching server...")
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer ln.Close()

	chat.InitializeChatRooms()

	handler := &handlers.ChatMessageHandler{}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		client := core.NewClient(conn)
		go handleConnection(client, handler)
	}
}

func handleConnection(client *core.Client, handler core.MessageHandler) {
	defer client.Conn.Close()

	client.Greet()

	for {
		message, err := client.ReadMessage()
		if err != nil {
			fmt.Println("Error reading:", err)
			return
		}

		handler.HandleMessage(client, message)
	}
}
