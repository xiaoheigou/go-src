package webportal

import (
	"yuudidi.com/pkg/protocol/webportal"
	"yuudidi.com/pkg/service"
	"yuudidi.com/pkg/utils"
)

func RunServer() error {
	// get configuration
	service.InitSendNotifyWheel()
	port := utils.Config.GetString("gin.port")
	return webportal.RunServer(port)
}
