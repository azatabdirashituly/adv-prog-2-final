package core

import (
    "sync"
)

type ChatRoom struct {
    Name    string
    Clients []*Client
    Creator *Client
    Lock    sync.Mutex
}

func NewChatRoom(name string, creator *Client) *ChatRoom {
    return &ChatRoom{Name: name, Clients: []*Client{}, Creator: creator}
}
