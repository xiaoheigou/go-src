package main

import (
	"fmt"
	"os"

	"yuudidi.com/pkg/cmd/server/web"
)

func main() {
	if err := web.RunServer(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
