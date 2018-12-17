package route

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"yuudidi.com/pkg/controller"
	"yuudidi.com/pkg/protocol/web/middleware"
)

func AppServer(t *gin.Engine) {
	r := t.Group("m")
	r.POST("/merchant/login", controller.AppLogin)
	r.POST("/merchant/register", controller.Register)
	r.POST("/merchant/random-code", controller.SendRandomCode)
	r.POST("/merchant/verify-identity", controller.VerifyRandomCode)
	r.POST("/merchant/reset-password", controller.ResetPw)
	r.GET("/merchants/:uid/audit-status", controller.GetAuditStatus)
	r.POST("/merchant/start-geetest", controller.RegisterGeetest)
	r.POST("/merchant/verify-geetest", controller.VerifyGeetest)

	g := r.Group("/")
	g.Use()
	{

		r.POST("merchant/logout", controller.AppLogout)
		r.POST("orders/:order-id/complaint", controller.OrderComplaint)
		merchants := g.Group("/merchants")
		{
			merchants.POST(":uid/change-password", controller.ChangePw)
			merchants.GET(":uid/profile", controller.GetProfile)
			merchants.GET(":uid/orders", controller.GetOrdersByMerchant)
			merchants.GET(":uid/orders/:order-id", controller.GetOrderDetail)
			merchants.PUT(":uid/settings/nickname", controller.SetNickname)
			merchants.GET(":uid/settings/work-mode", controller.GetWorkMode)
			merchants.PUT(":uid/settings/work-mode", controller.SetWorkMode)
			merchants.POST(":uid/settings/identities", controller.SetIdentities)
			merchants.PUT(":uid/settings/identities", controller.UpdateIdentities)
			merchants.POST(":uid/settings/identity/upload", controller.UploadIdentityFile)
			merchants.GET(":uid/settings/payments", controller.GetPayments)
			merchants.POST(":uid/settings/payments", controller.AddPayment)
			merchants.PUT(":uid/settings/payments/:id", controller.SetPayment)
			merchants.DELETE(":uid/settings/payments/:id", controller.DeletePayment)

		}
	}
}

func WebServer(t *gin.Engine) {
	r := t.Group("w")
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("session", store))

	r.Any("/login", controller.WebLogin)
	createOrder := t.Group("c")
	createOrder.Use()
	{
		createOrder.POST("create-order", controller.CreateOrder)
		createOrder.POST("reprocess-order", controller.ReprocessOrder)
		createOrder.GET("order/list", controller.GetOrderList)
	}

	g := r.Group("/")
	g.Use(middleware.Authenticated())
	{
		merchants := g.Group("merchants")
		{
			merchants.GET("", controller.GetMerchants)
			merchants.GET(":uid", controller.GetMerchant)
			merchants.POST(":uid/assets", controller.Recharge)
			merchants.GET(":uid/assets/history", controller.GetMerchantAssetHistory)
			merchants.PUT(":uid/assets/apply/:applyId", controller.RechargeConfirm)
			merchants.PUT(":uid/approve", controller.ApproveMerchant)
			merchants.PUT(":uid/freeze", controller.FreezeMerchant)
		}
		distributors := g.Group("distributors")
		{
			distributors.GET("", controller.GetDistributors)
			distributors.GET(":uid", controller.GetDistributor)
			distributors.GET(":uid/assets/history", controller.GetDistributorAssetHistory)
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
		users := g.Group("users")
		{
			users.POST("", controller.CreateUser)
			users.GET(":uid", controller.GetUser)
			users.GET("", controller.GetUsers)
		}
		recharges := g.Group("recharge")
		{
			recharges.GET("applies", controller.GetRechargeApplies)
		}
	}
}
