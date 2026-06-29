package routes

import (
	"speakup/pkg/handlers"
	"speakup/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

// BillingRoutes sets up the routes for Stripe billing
func BillingRoutes(router *gin.Engine, billingHandler *handlers.BillingHandler) {
	billingRoutes := router.Group("api/billing")
	billingRoutes.Use(middlewares.AuthMiddleware())
	{
		// billingRoutes.POST("/checkout", billingHandler.CreateCheckoutSession)
		// billingRoutes.POST("/portal", billingHandler.CreatePortalSession)
		billingRoutes.GET("/status", billingHandler.GetBillingStatus)
	}

	// router.POST("/api/billing/webhook", billingHandler.HandleWebhook)
}
