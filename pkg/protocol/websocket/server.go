package websocket

import (
	"github.com/gin-gonic/gin"
	"yuudidi.com/pkg/controller"
	"yuudidi.com/pkg/utils"
)

func RunServer(port string) error {

	r := gin.Default()
	controller.InitWheel()
	r.GET("/ws", controller.HandleWs)
	r.POST("/notify", controller.Notify)
	return r.Run(":" + port)
}

func init() {
	//服务重启删掉redis里面的key
	utils.RedisClient.Del(utils.RedisKeyMerchantOnline())
}
