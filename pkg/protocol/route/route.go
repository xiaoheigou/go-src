package route

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"yuudidi.com/pkg/controller"
)

func AppServer(t *gin.Engine) {
	r := t.Group("m")
	r.POST("/merchant/login", controller.AppLogin)
	r.POST("/merchant/register", controller.Register)
	r.GET("/merchant/random-code", controller.GetRandomCode)
	r.POST("/merchant/verify-identity", controller.VerifyRandomCode)
	r.POST("/merchant/reset-password", controller.ResetPw)
	r.GET("/merchants/:uid/audit-status", controller.GetAuditStatus)

	g := r.Group("/")
	g.Use()
	{

		r.POST("merchant/logout", controller.AppLogout)
		r.POST("orders/:order-id/complaint", controller.OrderComplaint)
		merchants := g.Group("/merchants")
		{
			r.POST("/merchant/change-password", controller.ChangePw)
			merchants.GET(":uid/profile", controller.GetProfile)
			merchants.GET(":uid/orders", controller.GetOrdersByMerchant)
			merchants.GET(":uid/orders/:order-id", controller.GetOrderDetail)
			merchants.PUT(":uid/settings/nickname", controller.SetNickName)
			merchants.GET(":uid/settings/work-mode", controller.GetWorkMode)
			merchants.PUT(":uid/settings/work-mode", controller.SetWorkMode)
			merchants.POST(":uid/settings/identities", controller.SetIdentities)
			merchants.PUT(":uid/settings/identities", controller.UpdateIdentities)
			merchants.POST(":uid/settings/identity/upload", controller.UploadIdentityFile)
			merchants.GET(":uid/settings/payments", controller.GetPayments)
			merchants.POST(":uid/settings/payments", controller.AddPayment)
			merchants.PUT(":uid/settings/payments", controller.SetPayment)
			merchants.DELETE(":uid/settings/payments", controller.DeletePayment)

		}
	}
}

func WebServer(t *gin.Engine) {
	r := t.Group("w")
	r.Any("/login", controller.WebLogin)
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("session", store))


	createOrder := t.Group("c")
	createOrder.Use()
	{
		createOrder.POST("create-order",controller.CreateOrder)
	}

	g := r.Group("/")
	g.Use()
	{
		merchants := g.Group("merchants")
		{
			merchants.GET("", controller.GetMerchants)
			merchants.GET(":uid", controller.GetMerchant)
			merchants.PUT(":uid/assets", controller.Recharge)
			merchants.GET(":uid/assets/history", controller.GetMerchantAssetHistory)
			merchants.PUT(":uid/approve", controller.ApproveMerchant)
			merchants.PUT(":uid/freeze", controller.FreezeMerchant)
		}
		distributors := g.Group("distributors")
		{
			distributors.GET("", controller.GetDistributors)
			distributors.GET(":uid", controller.GetDistributor)
			distributors.POST("", controller.CreateDistributors)
			distributors.PUT(":uid", controller.UpdateDistributors)
		}
		orders := g.Group("orders")
		{
			orders.GET("", controller.GetOrders)
			orders.GET(":id", controller.GetOrderByOrderNumber)
		}
		complaints := g.Group("complaints")
		{
			complaints.GET("", controller.GetComplaints)
			complaints.PUT(":id", controller.HandleComplaints)
		}

	}
}
