package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/danny-vo/fibonacci-backend/pkg/server"
)

func main() {
	if err := run(); nil != err {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func run() error {
	var err error = nil
	s := server.Initialize_Server()
	err = http.ListenAndServe("0.0.0.0:8080", s.Get_Router())

	return err
}
