package main

import (
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (

	// Tempop permitido para escrever uma mensagem
	writeWait = 10 * time.Second
	// Tempo permitidop para ler a proxima mensagem de pong
	pongWait = 60 * time.Second
	// Envia pings dentro desse periodo, tem que ser menor que pongWait
	pingPeriod = (pongWait * 9) / 10
	// Tamanho maximo permitido por mensagem
	maxMessageSize = 512
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

// Cliente é o middleman entre a conexao websocket e o hub
type Client struct {
	Id   string
	Hub  *Hub
	conn *websocket.Conn

	// chan []byte -> Responsavel por receber as mensagens que sao enviados pelo hub
	// e enviar como resposta escrita para o resto dos clientes
	Send chan []byte
}

// Lida com as requisiçoes websocket
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{
		Id:   uuid.New().String(),
		Hub:  hub,
		conn: conn,
		Send: make(chan []byte, 256),
	}

	// Registra o novo cliente no Hub
	client.Hub.Register <- client

	go client.writePump()
	go client.readPump()

}

func (c *Client) readPump() {
	defer func() {
		c.conn.Close()
		c.Hub.Unregister <- c
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(appData string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, text, err := c.conn.ReadMessage()
		log.Printf("Reading Message Value %v", string(text))
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		c.Hub.Broadcast <- &Message{MessageId: uuid.NewString(), ClientId: c.Id, Text: string(text)}
	}
}
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		c.conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.Send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(msg)
			log.Printf("Writing message with value %v", string(msg))
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}

	}
}
