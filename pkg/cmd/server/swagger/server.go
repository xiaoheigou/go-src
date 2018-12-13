package swagger

import (
	"yuudidi.com/pkg/protocol/swagger"
	"yuudidi.com/pkg/utils"
)

func RunServer() error {
	// get configuration
	port := utils.Config.GetString("gin.port")
	return swagger.RunServer(port)
}
