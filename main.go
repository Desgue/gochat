package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	hub := NewHub()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	go hub.Run()

	server := http.Server{
		Addr:              "localhost:8000",
		ReadHeaderTimeout: 3 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalln("err servig http: ", err)
	}
}
