package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/danny-vo/fibonacci-backend/pkg/server"
)

func main() {
	log.Println("Starting server...")
	for {
		if err := run(); nil != err {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			// os.Exit(1)
		}
	}
}

func run() error {
	var err error = nil
	s := server.InitializeServer()
	err = http.ListenAndServe("0.0.0.0:8080", s.GetRouter())

	return err
}
