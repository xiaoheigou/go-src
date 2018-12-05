package web_server

import (
	"YuuPay_core-service/pkg/protocol/web"
	"YuuPay_core-service/pkg/utils"
)

func RunServer() error {
	// get configuration
	port := utils.Config.GetString("gin.port")
	return web.RunServer(port)
}
