package main

import (
	"YuuPay_core-service/pkg/cmd/app-server"
	"fmt"
	"os"
)

func main() {
	if err := app_server.RunServer(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
