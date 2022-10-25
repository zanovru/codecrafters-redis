package main

import (
	"fmt"
	"net"
	"os"
	"strings"

	// Uncomment this block to pass the first stage
	// "net"
	// "os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleEcho(conn)
		//go handlePing(conn)
	}
}

func handleEcho(conn net.Conn) {
	for {
		defer conn.Close()
		var msg = make([]byte, 1024)
		if _, err := conn.Read(msg); err != nil {
			fmt.Println("Error reading from client")
			continue
		}
		strMsg := string(msg)
		idx := strings.Index(strMsg, "o")

		conn.Write([]byte(strMsg[idx+2:]))
	}
}
func handlePing(conn net.Conn) {
	for {
		defer conn.Close()
		if _, err := conn.Read([]byte{}); err != nil {
			fmt.Println("Error reading from client")
			continue
		}
		conn.Write([]byte("+PONG\r\n"))
	}
}
