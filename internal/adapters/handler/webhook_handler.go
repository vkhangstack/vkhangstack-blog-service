package handler

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/webhook"
	"github.com/vkhangstack/hexagonal-architecture/internal/config"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
	"github.com/vkhangstack/hexagonal-architecture/internal/core/services"
	"github.com/vkhangstack/hexagonal-architecture/internal/logger"
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

// ZaloBotWebhook handles POST /v1/webhooks/zalo — called by Zalo's servers (not our
// authenticated clients) whenever a user sends a message to our Zalo bot. If the message
// carries a valid verification code, it links the zalo_bot channel to the requesting user.
func (h *NotificationHandler) ZaloBotWebhook(ctx *gin.Context) {
	payload, err := ctx.GetRawData()
	if err != nil {
		HandleError(ctx, domain.ErrorCodePayloadBadRequest, nil, err.Error())
		return
	}
	signature := ctx.GetHeader(h.zaloBot.GetFieldSecretToken())

	senderID, text, err := h.zaloBot.ProcessWebhook(payload, signature)
	if errors.Is(err, domain.ErrZaloBotMessage) {
		HandleSuccess(ctx, nil, "Success")
		return
	}
	if err != nil {
		logger.Log.WithError(err).Error("ZaloBotWebhook: Failed to process webhook")
		HandleSuccess(ctx, nil, "Success")
		return
	}

	// Only messages carrying the NEXION-HUB-###### verification-code prefix are relevant here;
	// skip everything else to avoid noisy failed-verification logs for ordinary chat messages.
	if !services.VerificationCodePattern.MatchString(text) {
		HandleSuccess(ctx, nil, "Success")
		return
	}

	if err := h.settingSvc.VerifyChannelNotificationLinked(ctx, domain.NotificationChannelZaloBot, text, senderID); err != nil {
		logger.Log.WithError(err).Error("ZaloBotWebhook: Failed to verify channel")
		HandleSuccess(ctx, nil, "Success")
		return
	}
	HandleSuccess(ctx, nil, "Success")
}

// stripe webhook
func handleWebhook(c *gin.Context) {
	const MaxBodyBytes = int64(65536)
	payload, err := c.GetRawData()
	if err != nil {
		logger.Log.WithError(err).Error("error reading stripe webhook body")
		HandleError(c, domain.ErrorCodePayloadBadRequest, nil, "Error reading request body")
		return
	}

	event := stripe.Event{}

	if err := json.Unmarshal(payload, &event); err != nil {
		logger.Log.WithError(err).Error("stripe webhook: error parsing basic request")
		HandleError(c, domain.ErrorCodePayloadBadRequest, nil, "Webhook error while parsing basic request")
		return
	}

	endpointSecret := "whsec_"
	signatureHeader := c.GetHeader("Stripe-Signature")
	event, err = webhook.ConstructEvent(payload, signatureHeader, endpointSecret)
	if err != nil {
		logger.Log.WithError(err).Error("stripe webhook: signature verification failed")
		HandleError(c, domain.ErrorCodePayloadBadRequest, nil, "Webhook signature verification failed")
		return
	}
	switch event.Type {
	case "payment_intent.succeeded":
		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if err != nil {
			logger.Log.WithError(err).Error("stripe webhook: error parsing payment_intent JSON")
			HandleError(c, domain.ErrorCodePayloadBadRequest, nil, "Error parsing webhook JSON")
			return
		}
		logger.Log.WithField("amount", paymentIntent.Amount).Info("stripe: payment_intent.succeeded")
	case "payment_method.attached":
		var paymentMethod stripe.PaymentMethod
		err := json.Unmarshal(event.Data.Raw, &paymentMethod)
		if err != nil {
			logger.Log.WithError(err).Error("stripe webhook: error parsing payment_method JSON")
			HandleError(c, domain.ErrorCodePayloadBadRequest, nil, "Error parsing webhook JSON")
			return
		}
	default:
		logger.Log.WithField("event_type", event.Type).Warn("stripe webhook: unhandled event type")
	}

	HandleSuccess(c, nil, "Webhook updated successfully")
}
