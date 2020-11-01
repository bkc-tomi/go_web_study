package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s host:port", os.Args[0])
		os.Exit(1)
	}
	services := os.Args[1]
	udpAddr, err := net.ResolveUDPAddr("udp4", services)
	checkError(err)
	conn, err := net.DialUDP("udp", nil, udpAddr)
	checkError(err)
	defer conn.Close()
	_, err = conn.Write([]byte("anything"))
	checkError(err)
	var buf []byte
	n, err := conn.Read(buf)
	checkError(err)
	fmt.Println(string(buf[:n]))
	os.Exit(0)
}

func checkError(err error) {
	if err != nil {
		fmt.Fprint(os.Stderr, "Fatal error ", err.Error())
		os.Exit(1)
	}
}
