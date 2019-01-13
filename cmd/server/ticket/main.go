package main

import (
	"fmt"
	"os"
	"yuudidi.com/pkg/cmd/server/ticket"
)

func main() {
	if err := ticket.RunServer(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
