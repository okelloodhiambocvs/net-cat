package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"net-cat/internal/chat"
)

type Server struct {
	port string
	hub  *chat.Hub
}

func NewServer(port string) *Server {
	return &Server{
		port: port,
		hub:  chat.NewHub(),
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	defer listener.Close()

	fmt.Printf("Listening on the port :%s\n", s.port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			continue
		}

		go s.handleClient(conn)
	}
}

func (s *Server) handleClient(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	// Welcome message
	conn.Write([]byte("Welcome to TCP-Chat!\n"))
	conn.Write([]byte("         _nnnn_\n"))
	conn.Write([]byte("        dGGGGMMb\n"))
	conn.Write([]byte("[ENTER YOUR NAME]: "))

	// Read name ONCE
	nameRaw, err := reader.ReadString('\n')
	if err != nil {
		return
	}

	name := strings.TrimSpace(nameRaw)

	// Reject empty name
	if name == "" {
		conn.Write([]byte("Invalid name. Disconnecting...\n"))
		return
	}

	client := &chat.Client{
		Conn: conn,
		Name: name,
	}

	// Add client FIRST
	s.hub.AddClient(client)

	// THEN broadcast join (ONLY ONCE)
	s.hub.SystemMessage(fmt.Sprintf("%s has joined our chat.", name))

	// Message loop
	for {
		message, err := reader.ReadString('\n')

		// REAL DISCONNECT ONLY
		if err != nil {
			s.hub.RemoveClient(conn)
			s.hub.SystemMessage(fmt.Sprintf("%s has left our chat.", name))
			return
		}

		msg := strings.TrimSpace(message)

		// ignore empty messages
		if msg == "" {
			continue
		}

		s.hub.Broadcast(conn, msg)
	}
}