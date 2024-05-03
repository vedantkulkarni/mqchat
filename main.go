package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}



func socketHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Inside socket Handler")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("A socket error occured : %s", err.Error())
		return
	}
	conn.WriteMessage(1, []byte("Hello, World"))
}



func main() {

	fmt.Println("Hello WebSockets")
	http.HandleFunc("/", socketHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println("Error occured while listening to port 3333")
	}

}
