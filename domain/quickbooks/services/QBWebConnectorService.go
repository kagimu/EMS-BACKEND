package services

import (
	"fmt"
	"strings"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/quickbooks/models"
	uuid "github.com/satori/go.uuid"
)

type AuthResponse struct {
	Ticket string
	Status string
}

// Authenticate checks the provided username and password against the QBSettings table.
func Authenticate(username, password string) AuthResponse {
	var settings models.QBSettings
	if err := config.DB.Where("username = ? AND is_active = ?", username, true).First(&settings).Error; err != nil {
		return AuthResponse{Ticket: "", Status: "nvu"} // nvu = invalid user
	}

	// For simplicity in this example, we compare passwords.
	// In production, compare hashes.
	if settings.PasswordHash != password {
		return AuthResponse{Ticket: "", Status: "nvu"}
	}

	// Generate a unique session ticket
	ticket := uuid.NewV4().String()
	
	// Check if there's any work in the queue
	var count int64
	config.DB.Model(&models.QBSyncQueue{}).Where("status = ?", "QUEUED").Count(&count)
	
	if count == 0 {
		return AuthResponse{Ticket: ticket, Status: "none"} // "none" tells QBWC there is no work
	}

	return AuthResponse{Ticket: ticket, Status: ""}
}

// SendRequestXML retrieves the oldest queued item and returns it as a qbXML request
func SendRequestXML(ticket string) string {
	var queueItem models.QBSyncQueue
	
	err := config.DB.Where("status = ?", "QUEUED").Order("created_at asc").First(&queueItem).Error
	if err != nil {
		return ""
	}

	// Mark it as processing
	queueItem.Status = "PROCESSING"
	config.DB.Save(&queueItem)

	// In a complete implementation, this qbXML string would be constructed dynamically 
	// using the ItemReceiptNumber and Amount from the transactions table.
	// For this prototype, we return the raw payload stored in the queue.
	return `<?xml version="1.0" encoding="utf-8"?><?qbxml version="13.0"?><QBXML><QBXMLMsgsRq onError="stopOnError">` + queueItem.QbXmlPayload + `</QBXMLMsgsRq></QBXML>`
}

// ReceiveResponseXML handles the result from QuickBooks
func ReceiveResponseXML(ticket, response, hresult, message string) string {
	var queueItem models.QBSyncQueue
	
	// Find the item we just processed
	err := config.DB.Where("status = ?", "PROCESSING").Order("created_at asc").First(&queueItem).Error
	if err == nil {
		if strings.Contains(response, `statusCode="0"`) {
			queueItem.Status = "SUCCESS"
		} else {
			queueItem.Status = "FAILED"
			queueItem.ErrorMessage = message
		}
		config.DB.Save(&queueItem)
	}

	// Check if there is more work to do
	var count int64
	config.DB.Model(&models.QBSyncQueue{}).Where("status = ?", "QUEUED").Count(&count)
	
	if count == 0 {
		return "100" // 100% complete
	}
	
	// Rough estimate of completion, e.g., "50"
	return "50" 
}

func CloseConnection(ticket string) {
	fmt.Println("Closing QBWC session for ticket:", ticket)
}
