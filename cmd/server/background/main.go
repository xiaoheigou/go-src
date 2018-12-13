package main

import (
	"fmt"

	"yuudidi.com/pkg/cmd/server/background"
	"yuudidi.com/pkg/utils"
)

func testFunc(queue string, args ...interface{}) error {
	fmt.Printf("From %s, %v\n", queue, args)
	return nil
}

func main() {
	//register function, for test only! remove this in production code
	utils.RegisterWorkerFunc("test", testFunc)
	if err := background.LaunchBackgroundEngine(); err != nil {
		utils.Log.Fatalf("Background engine initialization failed: %v", err)
		panic("Initialization failure: Background engine does not launch!")
	}

}
