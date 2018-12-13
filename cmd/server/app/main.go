package main

import (
	"fmt"
	"os"

	"yuudidi.com/pkg/cmd/server/app"
)

func main() {
	if err := app.RunServer(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
