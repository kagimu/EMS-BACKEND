package notifier

import (
	"fmt"
	"html"
	"log"
	"os"
	"strings"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	commercialModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commercial/models"
	emailSender "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/emailSender"
	notificationServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/notifications/services"
	paymentModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/payments/models"
)

// NotifySubmitted emails the provider review inbox and confirms receipt to the
// customer. In-app notifications are not useful before an account/org exists.
func NotifySubmitted(submission *paymentModels.ManualPaymentSubmission) {
	if submission == nil {
		return
	}
	providerEmail := strings.TrimSpace(os.Getenv("MANUAL_PAYMENT_NOTIFY_EMAIL"))
	if providerEmail != "" {
		sendEmailAsync(providerEmail, "EMS manual payment proof submitted", fmt.Sprintf(
			"<p>A manual payment proof was submitted.</p><ul><li>Requester: %s</li><li>Order: %s</li><li>Method: %s</li><li>Amount: UGX %d</li><li>Reference: %s</li></ul><p>Open the Provider Console â†’ Manual Payments to review it.</p>",
			html.EscapeString(submission.RequesterName), html.EscapeString(submission.OrderNumber),
			html.EscapeString(submission.Method), submission.AmountUGX, html.EscapeString(submission.Number),
		))
	}
	if customerEmail := customerEmail(submission); customerEmail != "" {
		sendEmailAsync(customerEmail, "EMS payment proof received", fmt.Sprintf(
			"<p>Hello %s,</p><p>We received your proof of payment for order <strong>%s</strong>. It is awaiting review by the EMS team.</p><p>Reference: %s</p>",
			html.EscapeString(submission.RequesterName), html.EscapeString(submission.OrderNumber),
			html.EscapeString(submission.Number),
		))
	}
}

// ApprovalEmail builds the editable default email shown to the owner before
// approving a manual payment.
func ApprovalEmail(submission *paymentModels.ManualPaymentSubmission) (string, string) {
	if submission == nil {
		return "", ""
	}
	subject := fmt.Sprintf("EMS payment approved - order %s", submission.OrderNumber)
	body := fmt.Sprintf("Hello %s,\n\nYour manual payment for order %s has been approved.\n\nAmount: UGX %d\nReference: %s\n\nYou can now continue to claim your activation token.\n\nRegards,\nEMS", submission.RequesterName, submission.OrderNumber, submission.AmountUGX, submission.Number)
	return subject, body
}

// NotifyReviewed sends the customer email and, when an organisation exists, an
// in-app notification to its administrators. For approval, the owner-supplied
// subject/body are used; empty values fall back to the standard message.
func NotifyReviewed(submission *paymentModels.ManualPaymentSubmission, approved bool, emailSubject, emailBody, emailTo string) {
	if submission == nil {
		return
	}
	state := "rejected"
	title := "Manual payment rejected"
	if approved {
		state = "approved"
		title = "Manual payment approved"
	}
	body := fmt.Sprintf("Your manual payment for order %s was %s.", submission.OrderNumber, state)
	if !approved && strings.TrimSpace(submission.ReviewNote) != "" {
		body += " Reason: " + submission.ReviewNote
	}
	if submission.OrganisationNumber != "" {
		for _, recipient := range adminRecipients(submission.OrganisationNumber) {
			if err := notificationServices.Notify(notificationServices.NotificationInput{
				OrganisationNumber:  submission.OrganisationNumber,
				RecipientUserNumber: recipient,
				Type:                "payment.manual." + state,
				Title:               title,
				Body:                body,
				Payload: map[string]interface{}{
					"submission_number": submission.Number,
					"order_number":      submission.OrderNumber,
					"status":            submission.Status,
				},
			}); err != nil {
				log.Printf("manual payment in-app notification failed: %v", err)
			}
		}
	}
	customer := strings.TrimSpace(emailTo)
	if customer == "" {
		customer = customerEmail(submission)
	}
	if customer == "" {
		return
	}
	subject := strings.TrimSpace(emailSubject)
	text := strings.TrimSpace(emailBody)
	if approved && (subject == "" || text == "") {
		defaultSubject, defaultBody := ApprovalEmail(submission)
		if subject == "" {
			subject = defaultSubject
		}
		if text == "" {
			text = defaultBody
		}
	}
	if subject == "" {
		subject = "EMS payment " + state
	}
	if text == "" {
		text = body
	}
	subject = strings.ReplaceAll(strings.ReplaceAll(subject, "\r", " "), "\n", " ")
	sendEmailAsync(customer, subject, plainTextEmailHTML(text))
}

func plainTextEmailHTML(text string) string {
	lines := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
	var builder strings.Builder
	for _, line := range lines {
		builder.WriteString(html.EscapeString(line))
		builder.WriteString("<br/>")
	}
	return builder.String()
}
func customerEmail(submission *paymentModels.ManualPaymentSubmission) string {
	if value := strings.TrimSpace(submission.CustomerEmail); value != "" {
		return value
	}
	if submission.OrganisationNumber == "" {
		return ""
	}
	emails := adminEmails(submission.OrganisationNumber)
	if len(emails) > 0 {
		return emails[0]
	}
	return ""
}

func adminRecipients(organisationNumber string) []string {
	var rows []struct{ Number string }
	if err := config.DB.Table("users").Select("users.number").
		Joins("JOIN entity_users ON entity_users.user_number = users.number").
		Where("entity_users.entity_number = ? AND users.type = ? AND users.deleted_at IS NULL", organisationNumber, "Admin").
		Scan(&rows).Error; err != nil {
		return nil
	}
	result := make([]string, 0, len(rows))
	for _, row := range rows {
		if row.Number != "" {
			result = append(result, row.Number)
		}
	}
	return result
}

func adminEmails(organisationNumber string) []string {
	var rows []struct{ Contact string }
	if err := config.DB.Table("contacts").Select("contacts.contact").
		Joins("JOIN users ON users.number = contacts.user_number").
		Joins("JOIN entity_users ON entity_users.user_number = users.number").
		Where("entity_users.entity_number = ? AND users.type = ? AND users.deleted_at IS NULL AND LOWER(contacts.contact_type) = ?", organisationNumber, "Admin", "email").
		Scan(&rows).Error; err != nil {
		return nil
	}
	result := make([]string, 0, len(rows))
	for _, row := range rows {
		if strings.TrimSpace(row.Contact) != "" {
			result = append(result, strings.TrimSpace(row.Contact))
		}
	}
	return result
}

func reviewNoteHTML(note string) string {
	if strings.TrimSpace(note) == "" {
		return ""
	}
	return "<p>Reviewer note: " + html.EscapeString(note) + "</p>"
}

func sendEmailAsync(to, subject, body string) {
	go func() {
		if err := emailSender.SendSimpleEmail(to, subject, body); err != nil {
			log.Printf("manual payment email notification skipped: %v", err)
		}
	}()
}

// NotifyMobileMoneySuccessful sends the customer an in-app notification and,
// when an email address is available, an email after a mobile-money payment is
// independently confirmed.
func NotifyMobileMoneySuccessful(payment *paymentModels.MobileMoneyPayment) {
	if payment == nil {
		return
	}
	title := "Mobile money payment received"
	body := fmt.Sprintf("Payment of UGX %d for order %s was confirmed.", payment.AmountUGX, payment.Reference)
	if payment.OrganisationNumber != "" {
		for _, recipient := range adminRecipients(payment.OrganisationNumber) {
			if err := notificationServices.Notify(notificationServices.NotificationInput{
				OrganisationNumber:  payment.OrganisationNumber,
				RecipientUserNumber: recipient,
				Type:                "payment.mobile_money.success",
				Title:               title,
				Body:                body,
				Payload: map[string]interface{}{
					"payment_number": payment.Number, "order_number": payment.Reference,
					"provider": payment.Provider, "amount_ugx": payment.AmountUGX,
				},
			}); err != nil {
				log.Printf("mobile money in-app notification failed: %v", err)
			}
		}
	}
	if email := paymentEmail(payment); email != "" {
		sendEmailAsync(email, "EMS payment received", fmt.Sprintf(
			"<p>%s</p><p>Order: <strong>%s</strong><br/>Reference: %s<br/>Provider: %s</p>",
			html.EscapeString(body), html.EscapeString(payment.Reference),
			html.EscapeString(payment.Number), html.EscapeString(payment.Provider),
		))
	}
}

func paymentEmail(payment *paymentModels.MobileMoneyPayment) string {
	switch payment.Purpose {
	case paymentModels.PurposeSetup:
		var order commercialModels.SetupOrder
		if err := config.DB.Where("number = ?", payment.Reference).First(&order).Error; err == nil {
			return strings.TrimSpace(order.ContactEmail)
		}
	case paymentModels.PurposeSubscription:
		if payment.OrganisationNumber != "" {
			if emails := adminEmails(payment.OrganisationNumber); len(emails) > 0 {
				return emails[0]
			}
		}
		var order commercialModels.SubscriptionOrder
		if err := config.DB.Where("number = ?", payment.Reference).First(&order).Error; err == nil && order.SetupOrderNumber != "" {
			var setup commercialModels.SetupOrder
			if err := config.DB.Where("number = ?", order.SetupOrderNumber).First(&setup).Error; err == nil {
				return strings.TrimSpace(setup.ContactEmail)
			}
		}
	}
	return ""
}
