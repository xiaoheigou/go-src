package app

import (
	"YuuPay_core-service/pkg/api/v1"
	"YuuPay_core-service/pkg/api/v1/order"
	"YuuPay_core-service/pkg/api/v1/user"
	"YuuPay_core-service/pkg/utils"
	"github.com/gin-gonic/gin"
	"os"
	"path"
	"runtime"
)

func RunServer(port string) error {
	defer utils.DB.Close()
	defer utils.Log.OSFile.Close()
	r := gin.Default()

	//store := cookie.NewStore([]byte("secret"))
	//r.Use(sessions.Sessions("session", store))

	r.Any("/login",v1.AppLogin)
	g := r.Group("/")
	g.Use()
	{

		g.GET("/order",order.GetOrder)
		g.GET("/user",user.GetUser)
	}


	_, fileName, _, _ := runtime.Caller(0)
	rootPath := path.Join(fileName, "../../../../configs/")
	err := os.Chdir(rootPath)
	if err != nil {
		panic(err)
	}
	r.Run(":" + port)
	return nil
}