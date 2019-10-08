package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/bobinette/turn-o-matic"
)

func main() {
	addr := ":8098"

	// Default
	upgrader := websocket.Upgrader{}

	srv := turnomatic.Server{
		Upgrader: upgrader,
	}

	http.HandleFunc("/display", srv.Display)
	http.HandleFunc("/display-ws", srv.HandleDisplay)
	http.HandleFunc("/desk", srv.Desk)
	http.HandleFunc("/desk-ws", srv.HandleDesk)
	log.Println("listening on:", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
