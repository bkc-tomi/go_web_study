package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		socket()
	}
	service := os.Args[1]
	if service == "-g" {
		goroutinSocket()
	}
	fmt.Println("if you use concurrency socket, pleae set option -g.")
}

func socket() {
	service := ":7777"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	log.Println("normal socket\nlisten on port", service)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		daytime := time.Now().String()
		req := make([]byte, 1024)
		len, err := conn.Read(req)
		log.Println("riquest:", string(req[:len]))
		conn.Write([]byte(daytime))
		conn.Close()
	}
}

func goroutinSocket() {
	service := ":7777"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	log.Println("concurrency socket\nlisten on port", service)
	for {
		conn, err := listener.Accept()
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
		if n == 0 {
			break
		} else {
			c, err := parseReq(req)
			if err != nil {
				conn.Write([]byte(err.Error()))
			}
			ans, err := calc(c)
			if err != nil {
				conn.Write([]byte(err.Error()))
			}
			strAns := strconv.FormatFloat(ans, 'f', -1, 64)
			time.Sleep(1 * time.Second)
			conn.Write([]byte(strAns))
		}
	}
	fmt.Println("conn close.")
}

// Calc 演算子とオペランドをセットにした構造体
type Calc struct {
	Operator string
	x        float64
	y        float64
}

func parseReq(req string) (calc Calc, err error) {
	var temp []string
	n := 0
	for i, r := range req {
		if string(r) == "," {
			temp = append(temp, string(req[n:i]))
			n = i + 1
		}
	}
	if len(temp) < 3 {
		return
	}
	calc.Operator = temp[0]
	x, err := strconv.ParseFloat(temp[1], 64)
	if err != nil {
		return
	}
	y, err := strconv.ParseFloat(temp[2], 64)
	if err != nil {
		return
	}
	calc.x = x
	calc.y = y
	return
}

// CalcError calc実行時のエラーを返す型
type CalcError struct {
	s string
}

func (e *CalcError) Error() string {
	return "calc error." + e.s
}

func calc(c Calc) (float64, error) {
	switch c.Operator {
	case "add":
		return c.x + c.y, nil
	case "sub":
		return c.x - c.y, nil
	case "mult":
		return c.x * c.y, nil
	case "div":
		return c.x / c.y, nil
	default:
		return 0, &CalcError{"no match case"}
	}
}
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
