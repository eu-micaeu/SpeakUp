package handlers

import (
	"net/http"

	"speakup/pkg/middlewares"
	"speakup/pkg/repositories"
	"speakup/pkg/utils"

	"github.com/gin-gonic/gin"
)

type BillingHandler struct {
	Repo repositories.UserRepository
}

type BillingStatusResponse struct {
	StripeCustomerID       string `json:"stripe_customer_id"`
	StripeSubscriptionID   string `json:"stripe_subscription_id"`
	StripePriceID          string `json:"stripe_price_id"`
	StripeStatus           string `json:"stripe_status"`
	StripeCurrentPeriodEnd int64  `json:"stripe_current_period_end"`
}

func NewBillingHandler(repo repositories.UserRepository) *BillingHandler {
	return &BillingHandler{Repo: repo}
}

func (h *BillingHandler) CreateCheckoutSession(c *gin.Context) {
	utils.RespondWithError(c, http.StatusNotImplemented, "O SpeakUp é 100% gratuito!")
}

func (h *BillingHandler) CreatePortalSession(c *gin.Context) {
	utils.RespondWithError(c, http.StatusNotImplemented, "O SpeakUp é 100% gratuito!")
}

func (h *BillingHandler) GetBillingStatus(c *gin.Context) {
	userID := middlewares.GetUserIDFromContext(c)
	if userID == "" {
		utils.RespondWithError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, BillingStatusResponse{
		StripeCustomerID:       "free_user",
		StripeSubscriptionID:   "free_subscription",
		StripePriceID:          "free_price",
		StripeStatus:           "active",
		StripeCurrentPeriodEnd: 9999999999,
	})
}

