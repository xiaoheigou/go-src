package websocket

import (
	"yuudidi.com/pkg/protocol/websocket"
	"yuudidi.com/pkg/utils"
)

func RunServer() error {
	// get configuration
	url := utils.Config.GetString("websocket.port")
	return websocket.RunServer(url)
}
