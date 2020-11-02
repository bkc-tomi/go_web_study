package main

import (
	"container/list"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

var conns list.List

func echo(ws *websocket.Conn) {
	var err error
	var conn *list.Element
	contain := false
	for c := conns.Front(); c != nil; c = c.Next() {
		if c.Value == ws {
			contain = true
		}
	}
	if !contain {
		conn = conns.PushBack(ws)
	}
	for {
		var reply string
		clientHost := ws.Config().Origin

		if err = websocket.Message.Receive(ws, &reply); err != nil {
			conns.Remove(conn)
			ws.Close()
			fmt.Println("cant recieve", clientHost)
			break
		}

		fmt.Println("Recieve back from client:", reply)
		msg := "Recieved:" + reply + "from:" + clientHost.Host + "at:" + string(time.Now().Format("2006/1/2 15:04:05"))
		fmt.Println("Sending to client:", msg)

		for c := conns.Front(); c != nil; c = c.Next() {
			if err = websocket.Message.Send(c.Value.(*websocket.Conn), msg); err != nil {
				fmt.Println("cant send.")
				break
			}
			reply = ""
		}

	}
}

func main() {
	http.Handle("/", websocket.Handler(echo))

	if err := http.ListenAndServe(":1234", nil); err != nil {
		log.Fatal("Listen and Serve:", err)
	}
}
