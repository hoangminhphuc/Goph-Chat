package websocket

import "github.com/gorilla/websocket"


type Client struct {
	ID         string
	Connection *websocket.Conn
	Send chan []byte
}

func NewClient(id string, conn *websocket.Conn) *Client{
	return &Client {
		ID: id,
		Connection: conn,
		Send: make(chan []byte),
	}
}