package app

import (
	"yuudidi.com/pkg/protocol/app"
	"yuudidi.com/pkg/utils"
)

func RunServer() error {
	// get configuration
	port := utils.Config.GetString("gin.port")
	return app.RunServer(port)
}
