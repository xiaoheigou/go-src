package main

import (
	"YuuPay_core-service/pkg/cmd/web-server"
	"fmt"
	"os"
)

func main() {
	if err := web_server.RunServer(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
