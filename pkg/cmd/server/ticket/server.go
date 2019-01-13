package ticket

import (
	"yuudidi.com/pkg/protocol/ticket"
	"yuudidi.com/pkg/utils"
)

func RunServer() error {
	// get configuration
	port := utils.Config.GetString("gin.port")
	return ticket.RunServer(port)
}

