package route

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"yuudidi.com/pkg/api"
)

func AppServer(r *gin.Engine) {
	r.POST("/merchant/login", api.AppLogin)
	r.POST("/merchant/register", api.Register)
	r.GET("/merchant/randomcode", api.GetRandomCode)
	r.POST("/merchant/resetpassword", api.ResetPw)
	r.GET("/merchant/auditstatus", api.GetAuditStatus)

	g := r.Group("/")
	g.Use()
	{

		merchants := g.Group("/merchant")
		{
			merchants.POST("logout", api.AppLogout)
			merchants.GET("profile", api.GetProfile)
			merchants.GET("order", api.GetOrder)
			merchants.PUT("settings/nickname", api.SetNickName)
			merchants.GET("settings/workmode", api.GetWorkMode)
			merchants.PUT("settings/workmode", api.SetWorkMode)
			merchants.GET("settings/identify", api.GetIdentify)
			merchants.PUT("settings/identify", api.SetIdentify)
			merchants.GET("settings/payments", api.GetPayments)
			merchants.POST("settings/payments", api.AddPayment)
			merchants.PUT("settings/payments", api.SetPayment)
			merchants.DELETE("settings/payments", api.DeletePayment)

		}
		g.GET("/order", api.GetOrders)
	}
}

func WebServer(r *gin.Engine) {
	r.Any("/login", api.WebLogin)
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("session", store))

	g := r.Group("/")
	g.Use()
	{
		merchants := g.Group("merchants")
		{
			merchants.GET("", api.GetMerchants)
			merchants.PUT(":uid/assets", api.Recharge)
			merchants.GET(":uid/assets/history", api.GetMerchantAssetHistory)
			merchants.PUT(":uid/approve", api.ApproveMerchant)
			merchants.PUT(":uid/freeze", api.FreezeMerchant)
		}
		distributors := g.Group("distributors")
		{
			distributors.GET("", api.GetDistributors)
			distributors.POST("", api.CreateDistributors)
			distributors.PUT(":uid", api.UpdateDistributors)
		}
		orders := g.Group("orders")
		{
			orders.GET("", api.GetOrders)
		}
		complaints := g.Group("complaints")
		{
			complaints.GET("", api.GetComplaints)
			complaints.PUT(":id", api.HandleComplaints)
		}
	}
}
