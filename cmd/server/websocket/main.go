package main

import (
	"fmt"
	"os"
	"yuudidi.com/pkg/cmd/server/websocket"
)

func main() {
	//websocket config
	if err := websocket.RunServer(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}