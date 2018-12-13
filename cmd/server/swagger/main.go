package main

import (
	"fmt"
	"os"

	"yuudidi.com/pkg/cmd/server/swagger"
)

func main() {
	if err := swagger.RunServer(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
