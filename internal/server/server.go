package server

import (
	"fmt"
	"net"
)

type Server struct {
	port string
}

func NewServer(port string) *Server {
	return &Server{port: port}
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
			fmt.Println("failed to accept connection:", err)
			continue
		}

		fmt.Println("New client connected:", conn.RemoteAddr())

		// For now we just close immediately (STEP 2 only)
		conn.Close()
	}
}