package route

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"yuudidi.com/pkg/controller"
	"yuudidi.com/pkg/protocol/web/middleware"
	"yuudidi.com/pkg/utils"
)

func AppServer(t *gin.Engine) {
	r := t.Group("m")
	r.POST("/merchant/login", controller.AppLogin)
	r.POST("/merchant/register", controller.Register)
	r.POST("/merchant/random-code", controller.SendRandomCode)
	r.POST("/merchant/verify-identity", controller.VerifyRandomCode)
	r.POST("/merchant/reset-password", controller.ResetPw)

	r.GET("/merchants/:uid/audit-status", controller.GetAuditStatus) // 这个API不用认证

	g := r.Group("/")
	if utils.Config.GetString("appauth.skipauth") == "true" {
		g.Use()
	} else {
		g.Use(middleware.Auth(utils.Config.GetString("appauth.authkey")))
	}
	{
		merchants := g.Group("/merchants")
		{
			merchants.GET(":uid/svr-config", controller.GetSvrConfig)
			merchants.POST(":uid/logout", controller.AppLogout)
			merchants.GET(":uid/refresh-token", controller.RefreshToken)
			merchants.POST(":uid/change-password", controller.ChangePw)
			merchants.GET(":uid/profile", controller.GetProfile)
			merchants.GET(":uid/orders", controller.GetOrdersByMerchant)
			merchants.PUT(":uid/orders", controller.OrderFulfill)
			merchants.GET(":uid/orders/:order_number", controller.GetOrderDetail)
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
		createOrder.POST("create-order/buy", controller.BuyOrder)
		createOrder.POST("create-order/sell", controller.SellOrder)
		createOrder.GET("order/detail", controller.ReprocessOrder)
		createOrder.GET("order/list", controller.GetOrderList)
		createOrder.PUT("order/update", controller.UpdateOrder)
		createOrder.GET("order/query/:orderNumber", controller.GetOrderByOrderNumber)
		createOrder.POST("order/add", controller.AddOrder)
		createOrder.POST("ticket",controller.CreateTicket)
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
			//merchants.PUT(":uid/approve", controller.ApproveMerchant)
			//merchants.PUT(":uid/freeze", controller.FreezeMerchant)
			merchants.PUT(":uid/status", controller.ModifyMerchantStatus)
		}
		distributors := g.Group("distributors")
		{
			distributors.GET("", controller.GetDistributors)
			distributors.GET(":uid", controller.GetDistributor)
			distributors.POST(":uid/upload", controller.UploadCaPem)
			distributors.GET(":uid/assets/history", controller.GetDistributorAssetHistory)
			distributors.POST("", controller.CreateDistributors)
			distributors.PUT(":uid", controller.UpdateDistributors)
		}
		orders := g.Group("orders")
		{
			orders.GET("", controller.GetOrders)
			orders.GET("details/:orderNumber", controller.GetOrder)
			orders.PUT("refulfill/:orderNumber", controller.RefulfillOrder)
			orders.GET("ticket/:orderNumber", controller.GetTicket)
			orders.GET("status", controller.GetOrderStatus)
			orders.PUT("release/:orderNumber",controller.ReleaseCoin)
			orders.PUT("unfreeze/:orderNumber",controller.UnFreezeCoin)

		}
		tickets := g.Group("tickets")
		{
			tickets.GET(":ticketId", controller.GetTicketUpdates)
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
			users.PUT(":uid/password/reset", controller.ResetUserPassword)
			users.PUT(":uid/password", controller.UpdateUserPassword)
			users.PUT(":uid", controller.UpdateUser)
		}
		recharges := g.Group("recharge")
		{
			recharges.GET("applies", controller.GetRechargeApplies)
		}
	}
}
