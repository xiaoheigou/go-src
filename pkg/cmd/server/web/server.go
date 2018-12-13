package web

import (
	"yuudidi.com/pkg/protocol/web"
	"yuudidi.com/pkg/utils"
)

func RunServer() error {
	// get configuration
	port := utils.Config.GetString("gin.port")
	return web.RunServer(port)
}
