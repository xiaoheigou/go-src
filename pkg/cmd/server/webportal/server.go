package webportal

import (
	"yuudidi.com/pkg/protocol/webportal"
	"yuudidi.com/pkg/utils"
)

func RunServer() error {
	// get configuration
	port := utils.Config.GetString("gin.port")
	return webportal.RunServer(port)
}
