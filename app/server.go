package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
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

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	message := make([]byte, 1024)
	for {
		n, err := conn.Read(message)
		if err != nil {
			fmt.Println("Error reading from connection: ", err.Error())
		}
		response, err := Handle(message[:n])
		if err != nil {
			fmt.Println("Error handling the message:", err.Error())
		}

		conn.Write([]byte(response))
	}
}
