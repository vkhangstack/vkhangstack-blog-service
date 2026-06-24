package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/webhook"
	"github.com/vkhangstack/hexagonal-architecture/internal/config"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
	"github.com/vkhangstack/hexagonal-architecture/internal/utils"
)

type WebhookRequest struct {
	Event  string `json:"event"`
	UserId string `json:"user_id"`
}

func (h *UserHandler) UpdateMembershipStatus(ctx *gin.Context) {
	apiCfg := config.LoadConfig()

	// get api key from config
	apiKey := apiCfg.App.APIKey

	// check if api key is valid
	authHeader := ctx.Request.Header.Get("Authorization")
	if authHeader == "" {
		HandleError(ctx, domain.ErrorCodeUnAuthorization, nil, errors.New("no api key provided").Error())
		return
	}
	apiString := strings.TrimPrefix(authHeader, "ApiKey ")

	if apiString != apiKey {
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, errors.New("invalid api key").Error())
		return
	}

	var req WebhookRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}

	if req.Event != "membership_status_updated" {
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, errors.New("invalid event type").Error())
		return
	}

	userId, _ := utils.ParseUint64(req.UserId)

	err := h.svc.UpdateMembershipStatus(userId, true)
	if err != nil {
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}

	HandleSuccess(ctx, nil, "User's membership status updated successfully")
}

// stripe webhook
func handleWebhook(c *gin.Context) {
	const MaxBodyBytes = int64(65536)
	payload, err := c.GetRawData()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading request body: %v\n", err)
		HandleError(c, domain.ErrorCodePayloadBadRequest, nil, "Error reading request body")
		return
	}

	event := stripe.Event{}

	if err := json.Unmarshal(payload, &event); err != nil {
		fmt.Fprintf(os.Stderr, "⚠️  Webhook error while parsing basic request. %v\n", err.Error())
		HandleError(c, domain.ErrorCodePayloadBadRequest, nil, "Webhook error while parsing basic request")
		return
	}

	// Replace this endpoint secret with your endpoint's unique secret
	// If you are testing with the CLI, find the secret by running 'stripe listen'
	// If you are using an endpoint defined with the API or dashboard, look in your webhook settings
	// at https://dashboard.stripe.com/webhooks
	endpointSecret := "whsec_"
	signatureHeader := c.GetHeader("Stripe-Signature")
	event, err = webhook.ConstructEvent(payload, signatureHeader, endpointSecret)
	if err != nil {
		fmt.Fprintf(os.Stderr, "⚠️  Webhook signature verification failed. %v\n", err)
		HandleError(c, domain.ErrorCodePayloadBadRequest, nil, "Webhook signature verification failed")
		return
	}
	// Unmarshal the event data into an appropriate struct depending on its Type
	switch event.Type {
	case "payment_intent.succeeded":
		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
			HandleError(c, domain.ErrorCodePayloadBadRequest, nil, "Error parsing webhook JSON")
			return
		}
		log.Printf("Successful payment for %d.", paymentIntent.Amount)
		// Then define and call a func to handle the successful payment intent.
		// handlePaymentIntentSucceeded(paymentIntent)
	case "payment_method.attached":
		var paymentMethod stripe.PaymentMethod
		err := json.Unmarshal(event.Data.Raw, &paymentMethod)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
			HandleError(c, domain.ErrorCodePayloadBadRequest, nil, "Error parsing webhook JSON")
			return
		}
		// Then define and call a func to handle the successful attachment of a PaymentMethod.
		// handlePaymentMethodAttached(paymentMethod)
	default:
		fmt.Fprintf(os.Stderr, "Unhandled event type: %s\n", event.Type)
	}

	HandleSuccess(c, nil, "Webhook updated successfully")
}
