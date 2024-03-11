package main

import (
	"encoding/json"
	"log"
)

type Message struct {
	MessageId string `json:"messageId"`
	ClientId  string `json:"clientId"`
	Text      string `json:"text"`
}

// Hub mantem um conjunto de clientes ativos e envia menssagens aos clientes
type Hub struct {
	// Clientes registrados
	Clients  map[*Client]bool
	Messages []*Message

	// chan *Message ->  responsavel por receber as mensagens lidas
	// e enviar para cada um dos clientes registrados
	Broadcast chan *Message

	// Registra a requisiçao dos clientes
	Register chan *Client
	// Desregistra a requisiçao dos clientes
	Unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		Broadcast:  make(chan *Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
			for _, msg := range h.Messages {
				payload, err := json.Marshal(msg)
				if err != nil {
					log.Println("error decoding message: ", err)
				}
				client.Send <- payload
			}
			log.Println("New Client ID registered: ", client.Id)
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
				log.Println("Unregistering ")
			}
		case message := <-h.Broadcast:
			h.Messages = append(h.Messages, message)
			payload, err := json.Marshal(message)
			if err != nil {
				log.Println("error decoding message: ", err)
			}
			for client := range h.Clients {
				select {
				case client.Send <- payload:
				default:
					close(client.Send)
					delete(h.Clients, client)
				}
			}
		}
	}
}
