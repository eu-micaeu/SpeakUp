package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"

	"speakup/pkg/middlewares"
	"speakup/pkg/repositories"
	"speakup/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v79"
	billingportal "github.com/stripe/stripe-go/v79/billingportal/session"
	checkout "github.com/stripe/stripe-go/v79/checkout/session"
	"github.com/stripe/stripe-go/v79/customer"
	"github.com/stripe/stripe-go/v79/subscription"
	"github.com/stripe/stripe-go/v79/webhook"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type BillingHandler struct {
	Repo repositories.UserRepository
}

type CheckoutRequest struct {
	Plan string `json:"plan"`
}

type BillingStatusResponse struct {
	StripeCustomerID       string `json:"stripe_customer_id"`
	StripeSubscriptionID   string `json:"stripe_subscription_id"`
	StripePriceID          string `json:"stripe_price_id"`
	StripeStatus           string `json:"stripe_status"`
	StripeCurrentPeriodEnd int64  `json:"stripe_current_period_end"`
}

func NewBillingHandler(repo repositories.UserRepository) *BillingHandler {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	return &BillingHandler{Repo: repo}
}

func (h *BillingHandler) CreateCheckoutSession(c *gin.Context) {
	if err := validateStripeSecretKey(); err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	userID := middlewares.GetUserIDFromContext(c)
	if userID == "" {
		utils.RespondWithError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req CheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	priceID, err := resolvePriceID(req.Plan)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.Repo.FindByID(c, userID)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to load user")
		return
	}

	customerID := user.StripeCustomerID
	if customerID == "" {
		params := &stripe.CustomerParams{
			Email: stripe.String(user.Email),
			Name:  stripe.String(user.Name),
		}
		params.AddMetadata("user_id", user.ID)

		cust, err := customer.New(params)
		if err != nil {
			utils.RespondWithError(c, http.StatusInternalServerError, "Failed to create Stripe customer")
			return
		}

		customerID = cust.ID
		_ = h.Repo.Update(c, user.ID, mapStripeFields(customerID, "", "", "", 0))
	}

	successURL := os.Getenv("STRIPE_SUCCESS_URL")
	cancelURL := os.Getenv("STRIPE_CANCEL_URL")
	if successURL == "" || cancelURL == "" {
		utils.RespondWithError(c, http.StatusInternalServerError, "Stripe URLs not configured")
		return
	}

	params := &stripe.CheckoutSessionParams{
		Mode:                stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		Customer:            stripe.String(customerID),
		SuccessURL:          stripe.String(successURL),
		CancelURL:           stripe.String(cancelURL),
		AllowPromotionCodes: stripe.Bool(true),
		ClientReferenceID:   stripe.String(user.ID),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(priceID),
				Quantity: stripe.Int64(1),
			},
		},
		SubscriptionData: &stripe.CheckoutSessionSubscriptionDataParams{
			Metadata: map[string]string{
				"user_id": user.ID,
			},
		},
	}

	sess, err := checkout.New(params)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to create checkout session")
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, gin.H{"url": sess.URL})
}

func (h *BillingHandler) CreatePortalSession(c *gin.Context) {
	if err := validateStripeSecretKey(); err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	userID := middlewares.GetUserIDFromContext(c)
	if userID == "" {
		utils.RespondWithError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	user, err := h.Repo.FindByID(c, userID)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to load user")
		return
	}

	if user.StripeCustomerID == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Stripe customer not found")
		return
	}

	returnURL := os.Getenv("STRIPE_PORTAL_RETURN_URL")
	if returnURL == "" {
		utils.RespondWithError(c, http.StatusInternalServerError, "Stripe portal return URL not configured")
		return
	}

	params := &stripe.BillingPortalSessionParams{
		Customer:  stripe.String(user.StripeCustomerID),
		ReturnURL: stripe.String(returnURL),
	}

	portalSession, err := billingportal.New(params)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to create portal session")
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, gin.H{"url": portalSession.URL})
}

func (h *BillingHandler) GetBillingStatus(c *gin.Context) {
	userID := middlewares.GetUserIDFromContext(c)
	if userID == "" {
		utils.RespondWithError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	user, err := h.Repo.FindByID(c, userID)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to load user")
		return
	}

	// Sempre sincroniza com o Stripe se houver um customer ID para garantir dados atualizados
	if user.StripeCustomerID != "" {
		subscriptionData, priceID, err := fetchLatestSubscription(user.StripeCustomerID)
		if err == nil && subscriptionData != nil {
			user.StripeSubscriptionID = subscriptionData.ID
			user.StripePriceID = priceID
			user.StripeStatus = string(subscriptionData.Status)
			user.StripeCurrentPeriodEnd = subscriptionData.CurrentPeriodEnd
			_ = h.Repo.Update(c, user.ID, mapStripeFields(user.StripeCustomerID, subscriptionData.ID, priceID, user.StripeStatus, user.StripeCurrentPeriodEnd))
		}
	}

	utils.RespondWithJSON(c, http.StatusOK, BillingStatusResponse{
		StripeCustomerID:       user.StripeCustomerID,
		StripeSubscriptionID:   user.StripeSubscriptionID,
		StripePriceID:          user.StripePriceID,
		StripeStatus:           user.StripeStatus,
		StripeCurrentPeriodEnd: user.StripeCurrentPeriodEnd,
	})
}

func fetchLatestSubscription(customerID string) (*stripe.Subscription, string, error) {
	params := &stripe.SubscriptionListParams{
		Customer: stripe.String(customerID),
		Status:   stripe.String("all"),
	}
	params.Limit = stripe.Int64(1)
	params.AddExpand("data.items.data.price")

	iter := subscription.List(params)
	for iter.Next() {
		sub := iter.Subscription()
		priceID := ""
		if sub.Items != nil && len(sub.Items.Data) > 0 && sub.Items.Data[0].Price != nil {
			priceID = sub.Items.Data[0].Price.ID
		}
		return sub, priceID, nil
	}

	if err := iter.Err(); err != nil {
		return nil, "", err
	}

	return nil, "", nil
}

func (h *BillingHandler) HandleWebhook(c *gin.Context) {
	webhookSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	if webhookSecret == "" {
		utils.RespondWithError(c, http.StatusInternalServerError, "Stripe webhook secret not configured")
		return
	}

	payload, err := c.GetRawData()
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Failed to read request body")
		return
	}

	sigHeader := c.GetHeader("Stripe-Signature")
	event, err := webhook.ConstructEvent(payload, sigHeader, webhookSecret)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid webhook signature")
		return
	}

	switch event.Type {
	case "checkout.session.completed":
		var session stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &session); err == nil {
			h.handleCheckoutCompleted(c, &session)
		}
	case "customer.subscription.updated", "customer.subscription.deleted", "customer.subscription.created":
		var subscription stripe.Subscription
		if err := json.Unmarshal(event.Data.Raw, &subscription); err == nil {
			h.handleSubscriptionUpdate(c, &subscription)
		}
	}

	c.Status(http.StatusOK)
}

func (h *BillingHandler) handleCheckoutCompleted(c *gin.Context, session *stripe.CheckoutSession) {
	userID := session.ClientReferenceID
	if userID == "" {
		return
	}

	customerID := ""
	if session.Customer != nil {
		customerID = session.Customer.ID
	}
	subscriptionID := ""
	if session.Subscription != nil {
		subscriptionID = session.Subscription.ID
	}

	priceID := ""
	status := ""
	currentPeriodEnd := int64(0)
	if subscriptionID != "" {
		sub, err := subscription.Get(subscriptionID, nil)
		if err == nil && sub != nil {
			status = string(sub.Status)
			currentPeriodEnd = sub.CurrentPeriodEnd
			if sub.Items != nil && len(sub.Items.Data) > 0 && sub.Items.Data[0].Price != nil {
				priceID = sub.Items.Data[0].Price.ID
			}
		}
	}

	_ = h.Repo.Update(c, userID, mapStripeFields(customerID, subscriptionID, priceID, status, currentPeriodEnd))
}

func (h *BillingHandler) handleSubscriptionUpdate(c *gin.Context, subscription *stripe.Subscription) {
	customerID := ""
	if subscription.Customer != nil {
		customerID = subscription.Customer.ID
	}
	if customerID == "" {
		return
	}

	user, err := h.Repo.FindByStripeCustomerID(c, customerID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return
		}
		return
	}

	priceID := ""
	if subscription.Items != nil && len(subscription.Items.Data) > 0 && subscription.Items.Data[0].Price != nil {
		priceID = subscription.Items.Data[0].Price.ID
	}

	status := string(subscription.Status)
	currentPeriodEnd := subscription.CurrentPeriodEnd

	_ = h.Repo.Update(c, user.ID, mapStripeFields(customerID, subscription.ID, priceID, status, currentPeriodEnd))
}

func mapStripeFields(customerID, subscriptionID, priceID, status string, currentPeriodEnd int64) bson.M {
	update := bson.M{}
	if customerID != "" {
		update["stripe_customer_id"] = customerID
	}
	if subscriptionID != "" {
		update["stripe_subscription_id"] = subscriptionID
	}
	if priceID != "" {
		update["stripe_price_id"] = priceID
	}
	if status != "" {
		update["stripe_status"] = status
	}
	if currentPeriodEnd > 0 {
		update["stripe_current_period_end"] = currentPeriodEnd
	}
	return update
}

func resolvePriceID(plan string) (string, error) {
	switch plan {
	case "monthly":
		priceID := os.Getenv("STRIPE_PRICE_MONTHLY")
		if priceID == "" {
			return "", errEnv("STRIPE_PRICE_MONTHLY")
		}
		return priceID, nil
	case "annual":
		priceID := os.Getenv("STRIPE_PRICE_ANNUAL")
		if priceID == "" {
			return "", errEnv("STRIPE_PRICE_ANNUAL")
		}
		return priceID, nil
	default:
		return "", errBadPlan()
	}
}

func errBadPlan() error {
	return errors.New("Plano inválido. Use 'monthly' ou 'annual'.")
}

func errEnv(name string) error {
	return errors.New(name + " não configurado")
}

func validateStripeSecretKey() error {
	key := strings.TrimSpace(stripe.Key)
	if key == "" {
		return errors.New("Stripe key not configured")
	}
	if strings.HasPrefix(key, "pk_") {
		return errors.New("Stripe secret key must start with 'sk_'")
	}
	if !strings.HasPrefix(key, "sk_") {
		return errors.New("Invalid Stripe secret key format")
	}
	return nil
}
