package app_server

import (
	"YuuPay_core-service/pkg/protocol/app"
	"YuuPay_core-service/pkg/utils"
)

func RunServer() error {
	// get configuration
	port := utils.Config.GetString("gin.port")
	return app.RunServer(port)
}
