package websocket

import (
	"github.com/gin-gonic/gin"
	"os"
	"path"
	"runtime"
	"yuudidi.com/pkg/controller"
)

func RunServer(port string) error {

	r := gin.Default()

	_, fileName, _, _ := runtime.Caller(0)
	rootPath := path.Join(fileName, "../../../../configs/")
	err := os.Chdir(rootPath)
	if err != nil {
		panic(err)
	}
	r.GET("/ws", controller.HandleWs)
	r.POST("/notify", controller.Notify)
	return r.Run(":" + port)
}
