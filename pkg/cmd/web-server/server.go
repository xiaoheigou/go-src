package web_server

import (
	"yuudidi.com/pkg/protocol/web"
	"yuudidi.com/pkg/utils"
)

func RunServer() error {
	// get configuration
	port := utils.Config.GetString("gin.port")
	return web.RunServer(port)
}
