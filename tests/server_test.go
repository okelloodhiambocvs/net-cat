package tests

import (
	"net"
	"testing"
	"time"

	"net-cat/internal/server"
)

func TestServerStarts(t *testing.T) {
	srv := server.NewServer("9999")

	go func() {
		_ = srv.Start()
	}()

	time.Sleep(200 * time.Millisecond)

	conn, err := net.Dial("tcp", "localhost:9999")
	if err != nil {
		t.Fatalf("expected server to accept connection, got error: %v", err)
	}
	defer conn.Close()
}