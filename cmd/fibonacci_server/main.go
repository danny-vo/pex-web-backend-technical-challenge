package main

import (
	"fmt"
	"os"
)

func main() {
	if err := run(); nil != err {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func run() error {
	return nil
}
