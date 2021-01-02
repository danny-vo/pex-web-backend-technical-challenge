package main

import (
	"log"
	"net/http"
	"os"

	"github.com/danny-vo/fibonacci-backend/pkg/server"
)

func main() {
	log.Println("Starting server...")
	for {
		if err := run(); nil != err {
			log.Printf("Error occurred while serving: %v\n", err)
		}
	}
}

func run() error {
	var err error = nil
	s := server.InitializeServer()

	hostPort := os.Getenv("SERVING_HOST_PORT")
	if 0 == len(hostPort) {
		hostPort = "0.0.0.0:8080"
	}

	log.Println("Server has been initialized, now serving...")
	err = http.ListenAndServe(hostPort, s.GetRouter())

	return err
}
