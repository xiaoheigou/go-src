package main

import (
	"yuudidi.com/pkg/cmd/server/background"
	"yuudidi.com/pkg/utils"
)

func main() {
	if err := background.LaunchBackgroundEngine(); err != nil {
		utils.Log.Fatalf("Background engine initialization failed: %v", err)
		panic("Initialization failure: Background engine does not launch!")
	}
}