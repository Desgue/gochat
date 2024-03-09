package main

import (
	"fmt"
	"time"

	_ "github.com/gorilla/websocket"
)

func main() {
	ch := make(chan int)
	num := 1000000000000000
	go wait1(ch, num)

	wait2(num)
}

func wait1(ch chan<- int, num int) {
	start := time.Now()
	for i := range num {
		if i%100000000 == 0 {
			ch <- i
		}
	}
	fmt.Printf("goroutine total time elapsesed %s", time.Since(start))
	defer close(ch)
}

func wait2(num int) {
	start := time.Now()
	for i := range num {
		if i%100000000 == 0 {
			continue
		}
	}
	fmt.Printf("normal total time elapsesed %s", time.Since(start))
}

/* func Run() error {
	return http.ListenAndServe(":8000", nil)
}

func setuAPI() {
	manager := NewManager()

	http.HandleFunc("/ws", manager.serveWs)
}
*/
