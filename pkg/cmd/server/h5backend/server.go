package h5backend

import (
	"yuudidi.com/pkg/protocol/h5backend"
	"yuudidi.com/pkg/utils"
)

func RunServer() error {
	// get configuration
	port := utils.Config.GetString("gin.port")
	return h5backend.RunServer(port)
}

