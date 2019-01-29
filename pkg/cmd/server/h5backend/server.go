package h5backend

import (
	"yuudidi.com/pkg/protocol/h5backend"
	"yuudidi.com/pkg/service"
	"yuudidi.com/pkg/utils"
)

func RunServer() error {
	// get configuration
	service.InitSendNotifyWheel()
	port := utils.Config.GetString("gin.port")
	return h5backend.RunServer(port)
}

