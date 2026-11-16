package controllers

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	licensingServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/licensing/services"
	paymentServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/payments/services"
	"github.com/gin-gonic/gin"
)

type paymentRequest struct {
	Provider string `json:"provider"`
	Phone    string `json:"phone"`
}

func paymentError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, paymentServices.ErrPaymentNotFound):
		c.JSON(http.StatusNotFound, commons.Feedback(err.Error()))
	case errors.Is(err, paymentServices.ErrOrderNotPayable),
		errors.Is(err, paymentServices.ErrPaymentInProgress),
		errors.Is(err, paymentServices.ErrTokenAlreadyClaimed):
		c.JSON(http.StatusConflict, commons.Feedback(err.Error()))
	case errors.Is(err, paymentServices.ErrPaymentNotPaid):
		c.JSON(http.StatusPaymentRequired, commons.Feedback(err.Error()))
	case errors.Is(err, paymentServices.ErrInvalidPhone),
		errors.Is(err, paymentServices.ErrInvalidAmount),
		errors.Is(err, paymentServices.ErrProviderNotEnabled):
		c.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
	case errors.Is(err, paymentServices.ErrPaymentNotAllowed):
		c.JSON(http.StatusForbidden, commons.Feedback(err.Error()))
	case errors.Is(err, paymentServices.ErrProviderUnavailable):
		c.JSON(http.StatusBadGateway, commons.Feedback("The mobile-money provider could not be reached. Please try again."))
	default:
		c.JSON(http.StatusInternalServerError, commons.Feedback("The payment request could not be completed"))
	}
}

func currentOrganisation(c *gin.Context) (string, bool) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, commons.Feedback("Unauthorized"))
		return "", false
	}
	org, err := licensingServices.OrganisationNumberForUser(user.Number)
	if err != nil {
		c.JSON(http.StatusConflict, commons.Feedback("no active organisation for this account"))
		return "", false
	}
	return org, true
}

func paymentResponse(c *gin.Context, status int, number, organisation string) {
	view, err := paymentServices.Get(number, organisation)
	if err != nil {
		paymentError(c, err)
		return
	}
	c.JSON(status, gin.H{"data": view})
}

// Options lists the enabled mobile-money providers without revealing any
// credentials. Public because onboarding payments happen before sign-in.
func Options(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": paymentServices.Options()})
}

// InitiateSetupPayment starts an MTN MoMo or Airtel Money collection for a
// public installation order. The amount is taken from the order, not the body.
func InitiateSetupPayment(c *gin.Context) {
	var req paymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("invalid payment payload"))
		return
	}
	payment, err := paymentServices.Initiate(paymentServices.InitiateInput{
		Purpose:   "SETUP",
		Reference: c.Param("number"),
		Provider:  req.Provider,
		Phone:     req.Phone,
	})
	if err != nil && payment == nil {
		paymentError(c, err)
		return
	}
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"message": "The payment request was created but the provider could not be reached.",
			"data":    gin.H{"payment": payment},
		})
		return
	}
	paymentResponse(c, http.StatusCreated, payment.Number, "")
}

// InitiateOnboardingSubscriptionPayment starts a collection for a subscription
// order created before sign-in (bound to a paid setup token).
func InitiateOnboardingSubscriptionPayment(c *gin.Context) {
	var req paymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("invalid payment payload"))
		return
	}
	payment, err := paymentServices.Initiate(paymentServices.InitiateInput{
		Purpose:              "SUBSCRIPTION",
		Reference:            c.Param("number"),
		Provider:             req.Provider,
		Phone:                req.Phone,
		AllowOnboardingOrder: true,
	})
	if err != nil && payment == nil {
		paymentError(c, err)
		return
	}
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"message": "The payment request was created but the provider could not be reached.",
			"data":    gin.H{"payment": payment},
		})
		return
	}
	paymentResponse(c, http.StatusCreated, payment.Number, "")
}

// InitiateSubscriptionPayment starts a collection for an authenticated
// organisation's subscription order.
func InitiateSubscriptionPayment(c *gin.Context) {
	org, ok := currentOrganisation(c)
	if !ok {
		return
	}
	user, _ := commons.GetUserFromHeader(c)
	var req paymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("invalid payment payload"))
		return
	}
	payment, err := paymentServices.Initiate(paymentServices.InitiateInput{
		Purpose:            "SUBSCRIPTION",
		Reference:          c.Param("number"),
		Provider:           req.Provider,
		Phone:              req.Phone,
		ActorUserNumber:    user.Number,
		OrganisationNumber: org,
	})
	if err != nil && payment == nil {
		paymentError(c, err)
		return
	}
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"message": "The payment request was created but the provider could not be reached.",
			"data":    gin.H{"payment": payment},
		})
		return
	}
	paymentResponse(c, http.StatusCreated, payment.Number, org)
}

func PublicPaymentStatus(c *gin.Context) {
	// A status poll is also the client-driven verification trigger. Provider
	// outages return the last known state instead of hiding it behind a 502.
	if _, err := paymentServices.Get(c.Param("number"), ""); err != nil {
		paymentError(c, err)
		return
	}
	_, _ = paymentServices.Refresh(c.Param("number"))
	paymentResponse(c, http.StatusOK, c.Param("number"), "")
}

func ClientList(c *gin.Context) {
	org, ok := currentOrganisation(c)
	if !ok {
		return
	}
	page, _ := strconv.Atoi(c.Query("page"))
	perPage, _ := strconv.Atoi(c.Query("per_page"))
	result, err := paymentServices.List(paymentServices.ListFilter{
		OrganisationNumber: org,
		Provider:           c.Query("provider"),
		Status:             c.Query("status"),
		Search:             c.Query("search"),
		Page:               page,
		PerPage:            perPage,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, commons.Feedback("Payment history could not be loaded"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

func ClientPaymentStatus(c *gin.Context) {
	org, ok := currentOrganisation(c)
	if !ok {
		return
	}
	if _, err := paymentServices.Get(c.Param("number"), org); err != nil {
		paymentError(c, err)
		return
	}
	_, _ = paymentServices.Refresh(c.Param("number"))
	paymentResponse(c, http.StatusOK, c.Param("number"), org)
}

func PublicClaimPayment(c *gin.Context) {
	raw, payment, err := paymentServices.Claim(c.Param("number"), "")
	if err != nil {
		if payment != nil && errors.Is(err, paymentServices.ErrPaymentNotPaid) {
			c.JSON(http.StatusPaymentRequired, gin.H{"message": err.Error(), "data": gin.H{"payment": payment}})
			return
		}
		paymentError(c, err)
		return
	}
	view, _ := paymentServices.Get(payment.Number, "")
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"token": raw, "payment": view}})
}

func ClientClaimPayment(c *gin.Context) {
	org, ok := currentOrganisation(c)
	if !ok {
		return
	}
	raw, payment, err := paymentServices.Claim(c.Param("number"), org)
	if err != nil {
		if payment != nil && errors.Is(err, paymentServices.ErrPaymentNotPaid) {
			c.JSON(http.StatusPaymentRequired, gin.H{"message": err.Error(), "data": gin.H{"payment": payment}})
			return
		}
		paymentError(c, err)
		return
	}
	view, _ := paymentServices.Get(payment.Number, org)
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"token": raw, "payment": view}})
}

func MTNCallback(c *gin.Context) {
	handleCallback(c, "MTN_MOMO")
}

func AirtelCallback(c *gin.Context) {
	handleCallback(c, "AIRTEL_MONEY")
}

func handleCallback(c *gin.Context, provider string) {
	body, err := io.ReadAll(io.LimitReader(c.Request.Body, 1<<20))
	if err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("invalid callback payload"))
		return
	}
	if err := paymentServices.HandleCallback(provider, body); err != nil {
		// A provider outage is retryable; never acknowledge it as processed.
		c.JSON(http.StatusBadGateway, commons.Feedback("callback could not be verified yet"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "callback accepted"})
}

func OwnerList(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	perPage, _ := strconv.Atoi(c.Query("per_page"))
	filter := paymentServices.ListFilter{
		Search:             c.Query("search"),
		Provider:           c.Query("provider"),
		Status:             c.Query("status"),
		OrganisationNumber: c.Query("organisation"),
		Page:               page,
		PerPage:            perPage,
	}
	if value := strings.TrimSpace(c.Query("from")); value != "" {
		if parsed, err := time.Parse(time.RFC3339, value); err == nil {
			filter.From = &parsed
		}
	}
	if value := strings.TrimSpace(c.Query("to")); value != "" {
		if parsed, err := time.Parse(time.RFC3339, value); err == nil {
			filter.To = &parsed
		}
	}
	result, err := paymentServices.List(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, commons.Feedback("Payments could not be loaded"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

func OwnerDetail(c *gin.Context) {
	detail, err := paymentServices.Detail(c.Param("number"))
	if err != nil {
		paymentError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": detail})
}

func OwnerVerify(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil || user == nil {
		c.JSON(http.StatusUnauthorized, commons.Feedback("Unauthorized"))
		return
	}
	detail, err := paymentServices.Verify(c.Param("number"), user.Number)
	if err != nil {
		paymentError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Payment verification checked", "data": detail})
}
