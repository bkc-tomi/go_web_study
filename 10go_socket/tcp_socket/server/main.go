package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	service := ":7777"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)
	// _, err = net.ListenTCP("tcp", tcpAddr)
	listener, err := net.Listen("tcp", tcpAddr.String())
	checkError(err)
	log.Println("concurrency socket\nlisten on port", service)
	for {
		conn, err := listener.Accept()
		listener.Close()
		if err != nil {
			continue
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	buf := make([]byte, 24)
	defer conn.Close()
	for {
		n, _ := conn.Read(buf)
		req := string(buf[:n])
		fmt.Println("request:", req)
		if n == 0 {
			break
		} else {
			time.Sleep(5 * time.Second)
			conn.Write([]byte(string(time.Now().Format("Mon Jan 2 15:04:05"))))
		}
	}
	fmt.Println("conn close.")
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
