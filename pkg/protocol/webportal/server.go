package webportal

import (
	"github.com/gin-gonic/gin"
	"yuudidi.com/pkg/protocol/route"
	"yuudidi.com/pkg/utils"
)

func RunServer(port string) error {
	defer utils.DB.Close()
	defer utils.Log.OSFile.Close()
	r := gin.Default()

	route.WebServer(r)

	r.Run(":" + port)
	return nil
}
