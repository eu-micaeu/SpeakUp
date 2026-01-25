package models

type User struct {
    ID        string `json:"id" bson:"id"`
    Name      string `json:"name" bson:"name"`
    Email     string `json:"email" bson:"email"`
    Password  string `json:"password" bson:"password"`
    Language  string `json:"language" bson:"language"`
    Level     string `json:"level" bson:"level"`
    StripeCustomerID     string `json:"stripe_customer_id" bson:"stripe_customer_id"`
    StripeSubscriptionID string `json:"stripe_subscription_id" bson:"stripe_subscription_id"`
    StripePriceID        string `json:"stripe_price_id" bson:"stripe_price_id"`
    StripeStatus         string `json:"stripe_status" bson:"stripe_status"`
    StripeCurrentPeriodEnd int64 `json:"stripe_current_period_end" bson:"stripe_current_period_end"`
}