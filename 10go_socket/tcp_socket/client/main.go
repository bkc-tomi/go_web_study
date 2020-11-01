package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s host:port ", os.Args[0])
		os.Exit(1)
	}

	service := os.Args[1]
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err, "tcpAddr")
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err, "conn")
	defer conn.Close()
	for {
		input := scanStdin("please input calc,num1,num2, or 0 to exit \ncalc: add, sub, mult, div \nexample:add,1,1 \ninput:")
		if input == "0" {
			os.Exit(0)
		}
		_, err = conn.Write([]byte(input))
		checkError(err, "conn write")
		res := make([]byte, 1024)
		len, err := conn.Read(res)
		checkError(err, "conn read")
		fmt.Println("response:", string(res[:len]))
	}
}

func scanStdin(msg string) string {
	fmt.Print(msg)
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	return input.Text()
}

func checkError(err error, msg string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s \n", err.Error())
		fmt.Fprintf(os.Stderr, "message: %s \n", msg)
		os.Exit(1)
	}
}
