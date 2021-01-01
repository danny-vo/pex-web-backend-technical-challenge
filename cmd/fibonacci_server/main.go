package main

import (
	"log"
	"net/http"

	"github.com/danny-vo/fibonacci-backend/pkg/server"
)

func main() {
	log.Println("Starting server...")
	for {
		if err := run(); nil != err {
			log.Printf("Error occured while serving: %v\n", err)
		}
	}
}

func run() error {
	var err error = nil
	s := server.InitializeServer()

	log.Println("Server has been initialized, now serving...")
	err = http.ListenAndServe("0.0.0.0:8080", s.GetRouter())

	return err
}
