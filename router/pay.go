package router

import (
	"github.com/gin-gonic/gin"
	"parking-system-go/api"
)

type PayRouter struct{}

func (b *PayRouter) InitPayRouter(Router *gin.RouterGroup) {
	payeRouter := Router.Group("pay")
	payApi := api.ApiGroupApp.PayApi
	paywechatApi := api.ApiGroupApp.PayWeChatApi
	{
		payeRouter.POST("/create_mock_order", payApi.CreateMockOrder)
		payeRouter.POST("/mock_payment", payApi.MockPayment)
		payeRouter.GET("/order_status", payApi.GetOrderStatus)
		payeRouter.POST("/unifiedorder", paywechatApi.CreatePayment)
		payeRouter.POST("/payment_notify", paywechatApi.PaymentNotify)
		payeRouter.GET("/callback", paywechatApi.MockCallback)
	}
}
