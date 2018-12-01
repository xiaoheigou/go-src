package cmd

import (
	"otc-project/pkg/protocol/https"
	"otc-project/pkg/utils"
)

func RunServer() error {
	// get configuration
	port := utils.Config.GetString("gin.port")
	return https.RunServer(port)
}
