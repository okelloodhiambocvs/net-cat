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

	// Ask for name
	conn.Write([]byte("Welcome to TCP-Chat!\n"))
	conn.Write([]byte("         _nnnn_\n"))
	conn.Write([]byte("        dGGGGMMb\n"))
	conn.Write([]byte("[ENTER YOUR NAME]: "))

	reader := bufio.NewReader(conn)
	nameRaw, _ := reader.ReadString('\n')

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

	// Add to hub
	s.hub.AddClient(client)

	// SYSTEM JOIN MESSAGE
	s.hub.SystemMessage(fmt.Sprintf("%s has joined our chat.", name))

	// Keep connection alive (message loop placeholder)
	// Message loop
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			s.hub.RemoveClient(conn)
			s.hub.SystemMessage(fmt.Sprintf("%s has left our chat.", name))
			return
		}

		msg := strings.TrimSpace(message)

		// ignore empty messages (REQUIREMENT)
		if msg == "" {
			continue
		}

		// BROADCAST to others
		s.hub.Broadcast(conn, msg)
	}
}