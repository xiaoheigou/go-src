package main

import (
	"fmt"
	"os"
	"yuudidi.com/pkg/cmd/server/webportal"
)

func main() {
	if err := webportal.RunServer(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
