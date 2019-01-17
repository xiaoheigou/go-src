package main

import (
	"fmt"
	"os"
	"yuudidi.com/pkg/cmd/server/h5backend"
)

func main() {
	if err := h5backend.RunServer(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
