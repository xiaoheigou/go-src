package main

import (
	"fmt"
	"os"
	"yuudidi.com/pkg/cmd/swagger-server"
)

func main() {
	if err := swagger_server.RunServer(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
