package web

import (
	"YuuPay_core-service/pkg/api/v1"
	"YuuPay_core-service/pkg/api/v1/asset"
	"YuuPay_core-service/pkg/api/v1/complaint"
	"YuuPay_core-service/pkg/api/v1/distributor"
	merchant2 "YuuPay_core-service/pkg/api/v1/merchant"
	"YuuPay_core-service/pkg/api/v1/order"
	"YuuPay_core-service/pkg/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"os"
	"path"
	"runtime"

	_ "YuuPay_core-service/docs"
)

func RunServer(port string) error {
	defer utils.DB.Close()
	defer utils.Log.OSFile.Close()
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Any("/login",v1.WebLogin)
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("session", store))

	g := r.Group("/")
	g.Use()
	{
		merchants := g.Group("merchants")
		{
			merchants.GET("",merchant2.GetMerchants)
			merchants.PUT(":uid/assets",merchant2.Recharge)
			merchants.GET(":uid/assets/history", asset.GetMerchantAssetHistory)
			merchants.PUT(":uid/approve",merchant2.ApproveMerchant)
			merchants.PUT(":uid/freeze",merchant2.FreezeMerchant)
		}
		distributors := g.Group("distributors")
		{
			distributors.GET("",distributor.GetDistributors)
			distributors.POST("",distributor.CreateDistributors)
			distributors.PUT(":uid",distributor.UpdateDistributors)
		}
		orders := g.Group("orders")
		{
			orders.GET("",order.GetOrders)
		}
		complaints := g.Group("complaints")
		{
			complaints.GET("",complaint.GetComplaints)
			complaints.PUT(":id",complaint.HandleComplaints)
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