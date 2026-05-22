package main

import (
	"fmt"
	"log"
	"os"

	"net-cat/internal/server"
)

func main() {
	port := "8989"

	// If user provides port
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	// Invalid usage handling (assignment requirement)
	if len(os.Args) > 2 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}

	srv := server.NewServer(port)

	log.Printf("Listening on the port :%s\n", port)

	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}