package https

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"os"
	"otc-project/pkg/api/v1/order"
	"otc-project/pkg/api/v1/user"
	"otc-project/pkg/utils"
	"path"
	"runtime"
)

func RunServer(port string) error {
	defer utils.DB.Close()
	defer utils.Log.OSFile.Close()
	r := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("session", store))

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
	r.RunTLS(":" + port,rootPath + "/" + "server.pem",rootPath + "/" + "server.key")
	return nil
}