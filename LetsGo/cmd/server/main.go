package main

import (
	"log"

	"github.com/oneeyedsunday/dist-services-play/internal/server"
)

func main() {
	srv := server.NewHTTPServer(":5000")
	log.Printf("server running on %v", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
