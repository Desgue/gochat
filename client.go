package main

import (
	"log"

	"github.com/gorilla/websocket"
)

type ClientList map[*Client]bool

type Client struct {
	conn    *websocket.Conn
	manager *Manager
}

func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		conn:    conn,
		manager: manager,
	}
}

func (c *Client) readMessages() {
	defer func() {
		//cleanup connection
		c.manager.removeClient(c)
	}()

	for {
		msgType, payload, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("error reading message: ", err)
				break
			}
			log.Println(err)
			break
		}
		log.Println(msgType)
		log.Println(string(payload))
	}
}
