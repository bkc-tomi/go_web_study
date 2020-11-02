package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

func echo(ws *websocket.Conn) {
	var err error

	for {
		var reply string

		if err = websocket.Message.Receive(ws, &reply); err != nil {
			fmt.Println("cant recieve")
			break
		}

		fmt.Println("Recieve back from client:", reply)

		msg := "Recieved:" + reply + "at:" + string(time.Now().Format("2006/1/2 15:04:05"))
		fmt.Println("Sending to client:", msg)

		if err = websocket.Message.Send(ws, msg); err != nil {
			fmt.Println("cant send.")
			break
		}

	}
}

func main() {
	http.Handle("/", websocket.Handler(echo))

	if err := http.ListenAndServe(":1234", nil); err != nil {
		log.Fatal("Listen and Serve:", err)
	}
}
