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

		go s.handleClient(conn) // GOROUTINE (IMPORTANT)
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
		conn.Write([]byte("Name cannot be empty. Disconnecting...\n"))
		return
	}

	client := &chat.Client{
		Conn: conn,
		Name: name,
	}

	// Add to hub
	s.hub.AddClient(client)

	fmt.Println(name, "joined the chat")

	// Keep connection alive (message loop placeholder)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			s.hub.RemoveClient(conn)
			fmt.Println(name, "left the chat")
			return
		}

		msg := strings.TrimSpace(message)

		// ignore empty messages (REQUIREMENT)
		if msg == "" {
			continue
		}

		// For now we just print server-side
		fmt.Printf("[%s]: %s\n", name, msg)
	}
}