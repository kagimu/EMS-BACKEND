package controllers

import (
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	paymentModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/payments/models"
	paymentServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/payments/services"
	"github.com/gin-gonic/gin"
)

const manualUploadLimit = 12 << 20

func manualPaymentError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, paymentServices.ErrManualNotFound):
		c.JSON(http.StatusNotFound, commons.Feedback(err.Error()))
	case errors.Is(err, paymentServices.ErrManualNotPayable),
		errors.Is(err, paymentServices.ErrManualDuplicate),
		errors.Is(err, paymentServices.ErrManualNotApproved),
		errors.Is(err, paymentServices.ErrManualRejected),
		errors.Is(err, paymentServices.ErrManualAlready):
		c.JSON(http.StatusConflict, commons.Feedback(err.Error()))
	case errors.Is(err, paymentServices.ErrManualProof),
		errors.Is(err, paymentServices.ErrManualInvalid),
		errors.Is(err, paymentServices.ErrManualRejectReason):
		c.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
	default:
		c.JSON(http.StatusInternalServerError, commons.Feedback("The manual payment request could not be completed"))
	}
}

func ManualOptions(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": paymentServices.ManualOptions()})
}

func manualUpload(c *gin.Context, purpose, organisation string, allowOnboarding bool) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, manualUploadLimit)
	if err := c.Request.ParseMultipartForm(manualUploadLimit); err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("The proof upload is too large or malformed"))
		return
	}
	file, _ := c.FormFile("proof")
	submission, err := paymentServices.UploadManualPayment(paymentServices.ManualUploadInput{
		OrderNumber:        c.Param("number"),
		Purpose:            purpose,
		OrganisationNumber: organisation,
		AllowOnboarding:    allowOnboarding,
		Method:             c.PostForm("method"),
		PayerName:          c.PostForm("payer_name"),
		PayerPhone:         c.PostForm("payer_phone"),
		PayerReference:     c.PostForm("payer_reference"),
		CustomerNote:       c.PostForm("note"),
		CustomerName:       c.PostForm("customer_name"),
		CustomerEmail:      c.PostForm("customer_email"),
		CustomerPhone:      c.PostForm("customer_phone"),
		RequesterName:      c.PostForm("requester_name"),
		File:               file,
	})
	if err != nil {
		manualPaymentError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": paymentServices.ManualView(submission)})
}

func UploadSetupManualPayment(c *gin.Context) {
	manualUpload(c, paymentModels.PurposeSetup, "", false)
}

func UploadOnboardingSubscriptionManualPayment(c *gin.Context) {
	manualUpload(c, paymentModels.PurposeSubscription, "", true)
}

func UploadSubscriptionManualPayment(c *gin.Context) {
	org, ok := currentOrganisation(c)
	if !ok {
		return
	}
	manualUpload(c, paymentModels.PurposeSubscription, org, false)
}

func PublicManualPaymentStatus(c *gin.Context) {
	view, err := paymentServices.GetManualPayment(c.Param("number"), "")
	if err != nil {
		manualPaymentError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": view})
}

func ClientManualPaymentStatus(c *gin.Context) {
	org, ok := currentOrganisation(c)
	if !ok {
		return
	}
	view, err := paymentServices.GetManualPayment(c.Param("number"), org)
	if err != nil {
		manualPaymentError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": view})
}

func PublicManualPaymentClaim(c *gin.Context) {
	raw, view, err := paymentServices.ClaimManualPayment(c.Param("number"), "", "")
	if err != nil {
		manualPaymentError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"token": raw, "submission": view}})
}

func ClientManualPaymentClaim(c *gin.Context) {
	org, ok := currentOrganisation(c)
	if !ok {
		return
	}
	user, _ := commons.GetUserFromHeader(c)
	actor := ""
	if user != nil {
		actor = user.Number
	}
	raw, view, err := paymentServices.ClaimManualPayment(c.Param("number"), org, actor)
	if err != nil {
		manualPaymentError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"token": raw, "submission": view}})
}

func OwnerManualPaymentList(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	perPage, _ := strconv.Atoi(c.Query("per_page"))
	items, total, err := paymentServices.ListManualPayments(paymentServices.ManualListFilter{
		Status: c.Query("status"), Search: c.Query("search"), Page: page, PerPage: perPage,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, commons.Feedback("Manual payments could not be loaded"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"items": items, "total": total, "page": page, "per_page": perPage}})
}

func OwnerManualPaymentDetail(c *gin.Context) {
	detail, err := paymentServices.GetManualPaymentDetail(c.Param("number"))
	if err != nil {
		manualPaymentError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": detail})
}

func OwnerManualPaymentProof(c *gin.Context) {
	detail, err := paymentServices.GetManualPaymentDetail(c.Param("number"))
	if err != nil {
		manualPaymentError(c, err)
		return
	}
	if detail.ProofPath == "" {
		c.JSON(http.StatusNotFound, commons.Feedback("No proof file is attached"))
		return
	}
	if _, err := os.Stat(detail.ProofPath); err != nil {
		c.JSON(http.StatusNotFound, commons.Feedback("The proof file is no longer available"))
		return
	}
	c.FileAttachment(detail.ProofPath, detail.ProofName)
}

func OwnerManualPaymentApprove(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil || user == nil {
		c.JSON(http.StatusUnauthorized, commons.Feedback("Unauthorized"))
		return
	}
	var body struct {
		Note         string `json:"note"`
		EmailSubject string `json:"email_subject"`
		EmailBody    string `json:"email_body"`
		EmailTo      string `json:"email_to"`
	}
	_ = c.ShouldBindJSON(&body)
	view, err := paymentServices.ApproveManualPayment(c.Param("number"), user.Number, body.Note, body.EmailSubject, body.EmailBody, body.EmailTo)
	if err != nil {
		manualPaymentError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Manual payment approved", "data": view})
}

func OwnerManualPaymentReject(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil || user == nil {
		c.JSON(http.StatusUnauthorized, commons.Feedback("Unauthorized"))
		return
	}
	var body struct {
		Note string `json:"note"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || strings.TrimSpace(body.Note) == "" {
		c.JSON(http.StatusBadRequest, commons.Feedback("A rejection reason is required"))
		return
	}
	view, err := paymentServices.RejectManualPayment(c.Param("number"), user.Number, body.Note)
	if err != nil {
		manualPaymentError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Manual payment rejected", "data": view})
}
