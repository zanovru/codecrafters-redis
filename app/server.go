package main

import (
	"fmt"
	"net"
	"os"
)

type Server struct {
	handler *Handler
	storage *Storage
}

func NewServer(handler *Handler, storage *Storage) *Server {
	return &Server{
		handler: handler,
		storage: storage,
	}
}
func ListenAndServe(address string) {
	l, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	storage := NewStorage()
	handler := NewHandler(storage)
	server := NewServer(handler, storage)
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go server.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	message := make([]byte, 1024)
	for {
		n, err := conn.Read(message)
		if err != nil {
			fmt.Println("Error reading from connection: ", err.Error())
			return
		}
		response, err := s.handler.Handle(message[:n])
		if err != nil {
			fmt.Println("Error handling the message:", err.Error())
		}

		conn.Write([]byte(response))
	}
}
