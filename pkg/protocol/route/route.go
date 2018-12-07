package route

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"yuudidi.com/pkg/controller"
)

func AppServer(r *gin.Engine) {
	r.POST("/merchant/login", controller.AppLogin)
	r.POST("/merchant/register", controller.Register)
	r.GET("/merchant/randomcode", controller.GetRandomCode)
	r.POST("/merchant/resetpassword", controller.ResetPw)
	r.GET("/merchant/auditstatus", controller.GetAuditStatus)

	g := r.Group("/")
	g.Use()
	{

		merchants := g.Group("/merchant")
		{
			merchants.POST("logout", controller.AppLogout)
			merchants.GET("profile", controller.GetProfile)
			merchants.GET("order", controller.GetOrder)
			merchants.PUT("settings/nickname", controller.SetNickName)
			merchants.GET("settings/workmode", controller.GetWorkMode)
			merchants.PUT("settings/workmode", controller.SetWorkMode)
			merchants.GET("settings/identify", controller.GetIdentify)
			merchants.PUT("settings/identify", controller.SetIdentify)
			merchants.GET("settings/payments", controller.GetPayments)
			merchants.POST("settings/payments", controller.AddPayment)
			merchants.PUT("settings/payments", controller.SetPayment)
			merchants.DELETE("settings/payments", controller.DeletePayment)

		}
		g.GET("/order", controller.GetOrders)
	}
}

func WebServer(r *gin.Engine) {
	r.Any("/login", controller.WebLogin)
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("session", store))

	g := r.Group("/")
	g.Use()
	{
		merchants := g.Group("merchants")
		{
			merchants.GET("", controller.GetMerchants)
			merchants.PUT(":uid/assets", controller.Recharge)
			merchants.GET(":uid/assets/history", controller.GetMerchantAssetHistory)
			merchants.PUT(":uid/approve", controller.ApproveMerchant)
			merchants.PUT(":uid/freeze", controller.FreezeMerchant)
		}
		distributors := g.Group("distributors")
		{
			distributors.GET("", controller.GetDistributors)
			distributors.POST("", controller.CreateDistributors)
			distributors.PUT(":uid", controller.UpdateDistributors)
		}
		orders := g.Group("orders")
		{
			orders.GET("", controller.GetOrders)
		}
		complaints := g.Group("complaints")
		{
			complaints.GET("", controller.GetComplaints)
			complaints.PUT(":id", controller.HandleComplaints)
		}
	}
}
