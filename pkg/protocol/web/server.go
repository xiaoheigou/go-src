package web

import (
	"yuudidi.com/pkg/api"
	"yuudidi.com/pkg/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"os"
	"path"
	"runtime"

	_ "yuudidi.com/docs"
)

func RunServer(port string) error {
	defer utils.DB.Close()
	defer utils.Log.OSFile.Close()
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Any("/login",api.WebLogin)
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("session", store))

	g := r.Group("/")
	g.Use()
	{
		merchants := g.Group("merchants")
		{
			merchants.GET("",api.GetMerchants)
			merchants.PUT(":uid/assets",api.Recharge)
			merchants.GET(":uid/assets/history", api.GetMerchantAssetHistory)
			merchants.PUT(":uid/approve",api.ApproveMerchant)
			merchants.PUT(":uid/freeze",api.FreezeMerchant)
		}
		distributors := g.Group("distributors")
		{
			distributors.GET("",api.GetDistributors)
			distributors.POST("",api.CreateDistributors)
			distributors.PUT(":uid",api.UpdateDistributors)
		}
		orders := g.Group("orders")
		{
			orders.GET("",api.GetOrders)
		}
		complaints := g.Group("complaints")
		{
			complaints.GET("",api.GetComplaints)
			complaints.PUT(":id",api.HandleComplaints)
		}
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