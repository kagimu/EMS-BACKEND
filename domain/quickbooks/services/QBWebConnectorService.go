package services

import (
	"encoding/xml"
	"fmt"
	"strings"
	"sync"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/quickbooks/models"
	uuid "github.com/satori/go.uuid"
)

// activeSessions maps session tickets to QB profile numbers
var activeSessions sync.Map

// Helper to check if a profile is the default one
func isDefaultProfile(profileNumber string) bool {
	if profileNumber == "" {
		return true
	}
	var p models.QBProfile
	if err := config.DB.Where("number = ?", profileNumber).First(&p).Error; err == nil {
		return p.IsDefault
	}
	return false
}

type QBCustomerRet struct {
	ListID       string `xml:"ListID"`
	EditSequence string `xml:"EditSequence"`
	FullName     string `xml:"FullName"`
	Name         string `xml:"Name"`
	FirstName    string `xml:"FirstName"`
	LastName     string `xml:"LastName"`
}

type QBCustomerQueryRs struct {
	XMLName     xml.Name        `xml:"CustomerQueryRs"`
	StatusCode  string          `xml:"statusCode,attr"`
	CustomerRet []QBCustomerRet `xml:"CustomerRet"`
}

type QBCustomerAddRs struct {
	XMLName       xml.Name      `xml:"CustomerAddRs"`
	StatusCode    string        `xml:"statusCode,attr"`
	StatusMessage string        `xml:"statusMessage,attr"`
	CustomerRet   QBCustomerRet `xml:"CustomerRet"`
}

type QBSalesReceiptRet struct {
	TxnID string `xml:"TxnID"`
}

type QBSalesReceiptAddRs struct {
	XMLName         xml.Name          `xml:"SalesReceiptAddRs"`
	StatusCode      string            `xml:"statusCode,attr"`
	StatusMessage   string            `xml:"statusMessage,attr"`
	SalesReceiptRet QBSalesReceiptRet `xml:"SalesReceiptRet"`
}

type QBItemServiceRet struct {
	ListID       string `xml:"ListID"`
	EditSequence string `xml:"EditSequence"`
	Name         string `xml:"Name"`
	FullName     string `xml:"FullName"`
}

type QBItemServiceAddRs struct {
	XMLName        xml.Name         `xml:"ItemServiceAddRs"`
	StatusCode     string           `xml:"statusCode,attr"`
	StatusMessage  string           `xml:"statusMessage,attr"`
	ItemServiceRet QBItemServiceRet `xml:"ItemServiceRet"`
}

// Inventory Item response structs
type QBItemInventoryRet struct {
	ListID       string `xml:"ListID"`
	EditSequence string `xml:"EditSequence"`
	Name         string `xml:"Name"`
	FullName     string `xml:"FullName"`
}

type QBItemInventoryAddRs struct {
	XMLName           xml.Name           `xml:"ItemInventoryAddRs"`
	StatusCode        string             `xml:"statusCode,attr"`
	StatusMessage     string             `xml:"statusMessage,attr"`
	ItemInventoryRet  QBItemInventoryRet `xml:"ItemInventoryRet"`
}

// Non-Inventory Item response structs
type QBItemNonInventoryRet struct {
	ListID       string `xml:"ListID"`
	EditSequence string `xml:"EditSequence"`
	Name         string `xml:"Name"`
	FullName     string `xml:"FullName"`
}

type QBItemNonInventoryAddRs struct {
	XMLName              xml.Name              `xml:"ItemNonInventoryAddRs"`
	StatusCode           string                `xml:"statusCode,attr"`
	StatusMessage        string                `xml:"statusMessage,attr"`
	ItemNonInventoryRet  QBItemNonInventoryRet `xml:"ItemNonInventoryRet"`
}

type QBXMLResponse struct {
	XMLName     xml.Name `xml:"QBXML"`
	QBXMLMsgsRs struct {
		CustomerQueryRs        *QBCustomerQueryRs        `xml:"CustomerQueryRs"`
		CustomerAddRs          *QBCustomerAddRs          `xml:"CustomerAddRs"`
		SalesReceiptAddRs      *QBSalesReceiptAddRs      `xml:"SalesReceiptAddRs"`
		ItemServiceAddRs       *QBItemServiceAddRs       `xml:"ItemServiceAddRs"`
		ItemInventoryAddRs     *QBItemInventoryAddRs     `xml:"ItemInventoryAddRs"`
		ItemNonInventoryAddRs  *QBItemNonInventoryAddRs  `xml:"ItemNonInventoryAddRs"`
	} `xml:"QBXMLMsgsRs"`
}

type AuthResponse struct {
	Ticket string
	Status string
}

// Authenticate checks the provided username and password against QBProfile (and fallback to QBSettings).
func Authenticate(username, password string) AuthResponse {
	var profile models.QBProfile
	var profileNumber string

	if err := config.DB.Where("username = ? AND is_active = ?", username, true).First(&profile).Error; err == nil {
		if profile.PasswordHash != password {
			return AuthResponse{Ticket: "", Status: "nvu"} // nvu = invalid user
		}
		profileNumber = profile.Number
	} else {
		// Fallback to legacy QBSettings
		var settings models.QBSettings
		if err := config.DB.Where("username = ? AND is_active = ?", username, true).First(&settings).Error; err != nil {
			return AuthResponse{Ticket: "", Status: "nvu"} // nvu = invalid user
		}
		if settings.PasswordHash != password {
			return AuthResponse{Ticket: "", Status: "nvu"}
		}
	}

	// Generate a unique session ticket
	ticket := uuid.NewV4().String()
	activeSessions.Store(ticket, profileNumber)

	// Check if there's any work in the queue for this profile
	var count int64
	query := config.DB.Model(&models.QBSyncQueue{}).Where("status = ?", "QUEUED")
	if isDefaultProfile(profileNumber) {
		query = query.Where("qb_profile_number = ? OR qb_profile_number = '' OR qb_profile_number IS NULL", profileNumber)
	} else {
		query = query.Where("qb_profile_number = ?", profileNumber)
	}
	query.Count(&count)

	if count == 0 {
		return AuthResponse{Ticket: ticket, Status: "none"} // "none" tells QBWC there is no work
	}

	return AuthResponse{Ticket: ticket, Status: ""}
}

// SendRequestXML retrieves the oldest queued item and returns it as a qbXML request
func SendRequestXML(ticket string) string {
	var profileNumber string
	if val, ok := activeSessions.Load(ticket); ok {
		profileNumber, _ = val.(string)
	}

	var queueItem models.QBSyncQueue
	query := config.DB.Where("status = ?", "QUEUED")
	if isDefaultProfile(profileNumber) {
		query = query.Where("qb_profile_number = ? OR qb_profile_number = '' OR qb_profile_number IS NULL", profileNumber)
	} else {
		query = query.Where("qb_profile_number = ?", profileNumber)
	}
	err := query.Order("priority asc, created_at asc").First(&queueItem).Error
	if err != nil {
		return ""
	}

	// Mark it as processing
	queueItem.Status = "PROCESSING"
	config.DB.Save(&queueItem)

	return `<?xml version="1.0" encoding="utf-8"?><?qbxml version="13.0"?><QBXML><QBXMLMsgsRq onError="stopOnError">` + queueItem.QbXmlPayload + `</QBXMLMsgsRq></QBXML>`
}

// ReceiveResponseXML handles the result from QuickBooks
func ReceiveResponseXML(ticket, response, hresult, message string) string {
	var profileNumber string
	if val, ok := activeSessions.Load(ticket); ok {
		profileNumber, _ = val.(string)
	}

	var queueItem models.QBSyncQueue

	// Find the item we just processed for this profile
	query := config.DB.Where("status = ?", "PROCESSING")
	if isDefaultProfile(profileNumber) {
		query = query.Where("qb_profile_number = ? OR qb_profile_number = '' OR qb_profile_number IS NULL", profileNumber)
	} else {
		query = query.Where("qb_profile_number = ?", profileNumber)
	}
	err := query.Order("created_at asc").First(&queueItem).Error
	if err == nil {
		if hresult != "" && hresult != "0x00000000" {
			// QB itself returned a COM-level error (e.g. bad XML)
			queueItem.ErrorMessage = message
			queueItem.RetryCount++
			if queueItem.RetryCount < queueItem.MaxRetries {
				queueItem.Status = "QUEUED"
			} else {
				queueItem.Status = "FAILED"
			}
		} else if strings.Contains(response, `statusCode="0"`) {
			queueItem.Status = "SUCCESS"

			if queueItem.ActionType == "CUSTOMER_ADD" || queueItem.ActionType == "CUSTOMER_QUERY" || queueItem.ActionType == "SALES_RECEIPT_ADD" || queueItem.ActionType == "ITEM_SERVICE_ADD" || queueItem.ActionType == "ITEM_INVENTORY_ADD" || queueItem.ActionType == "ITEM_NON_INVENTORY_ADD" {
				var qbxml QBXMLResponse
				xml.Unmarshal([]byte(response), &qbxml)

				targetProfileNumber := queueItem.QBProfileNumber
				if targetProfileNumber == "" {
					targetProfileNumber = profileNumber
				}

				if queueItem.ActionType == "CUSTOMER_ADD" && qbxml.QBXMLMsgsRs.CustomerAddRs != nil {
					if qbxml.QBXMLMsgsRs.CustomerAddRs.StatusCode == "0" {
						ret := qbxml.QBXMLMsgsRs.CustomerAddRs.CustomerRet
						mapping := models.QBCustomerMapping{
							UserNumber:      queueItem.UserNumber,
							QBProfileNumber: targetProfileNumber,
							QbListID:        ret.ListID,
							QbEditSequence:  ret.EditSequence,
							QbFullName:      ret.FullName,
						}
						config.DB.Create(&mapping)
					}
				} else if queueItem.ActionType == "CUSTOMER_QUERY" && qbxml.QBXMLMsgsRs.CustomerQueryRs != nil {
					for _, ret := range qbxml.QBXMLMsgsRs.CustomerQueryRs.CustomerRet {
						var mapping models.QBCustomerMapping
						q := config.DB.Where("qb_list_id = ?", ret.ListID)
						if targetProfileNumber != "" {
							q = q.Where("qb_profile_number = ?", targetProfileNumber)
						}
						
						if err := q.First(&mapping).Error; err == nil {
							// Only update if it already exists for this profile
							mapping.QbEditSequence = ret.EditSequence
							mapping.QbFullName = ret.FullName
							config.DB.Save(&mapping)
						} else if targetProfileNumber != "" {
							// If it doesn't exist for this profile, check if we have it ANYWHERE else (matching by FullName)
							var existingMapping models.QBCustomerMapping
							if err := config.DB.Where("qb_full_name = ?", ret.FullName).Where("qb_full_name != ''").First(&existingMapping).Error; err == nil {
								// Duplicate the mapping for the new profile, securely linking it to the same UserNumber
								newMapping := models.QBCustomerMapping{
									UserNumber:      existingMapping.UserNumber,
									QBProfileNumber: targetProfileNumber,
									QbListID:        ret.ListID,
									QbEditSequence:  ret.EditSequence,
									QbFullName:      ret.FullName,
								}
								config.DB.Create(&newMapping)
							}
						}
					}
				} else if queueItem.ActionType == "SALES_RECEIPT_ADD" && qbxml.QBXMLMsgsRs.SalesReceiptAddRs != nil {
					// Fix 4: Store QB TxnID for tracking
					queueItem.QbTxnID = qbxml.QBXMLMsgsRs.SalesReceiptAddRs.SalesReceiptRet.TxnID
				} else if queueItem.ActionType == "ITEM_SERVICE_ADD" && qbxml.QBXMLMsgsRs.ItemServiceAddRs != nil {
					if qbxml.QBXMLMsgsRs.ItemServiceAddRs.StatusCode == "0" {
						ret := qbxml.QBXMLMsgsRs.ItemServiceAddRs.ItemServiceRet
						mapping := models.QBItemMapping{
							ItemNumber:      queueItem.ItemNumber,
							QBProfileNumber: targetProfileNumber,
							QbListID:        ret.ListID,
							QbEditSequence:  ret.EditSequence,
							QbFullName:      ret.FullName,
						}
						config.DB.Create(&mapping)
					}
				} else if queueItem.ActionType == "ITEM_INVENTORY_ADD" && qbxml.QBXMLMsgsRs.ItemInventoryAddRs != nil {
					if qbxml.QBXMLMsgsRs.ItemInventoryAddRs.StatusCode == "0" {
						ret := qbxml.QBXMLMsgsRs.ItemInventoryAddRs.ItemInventoryRet
						mapping := models.QBItemMapping{
							ItemNumber:      queueItem.ItemNumber,
							QBProfileNumber: targetProfileNumber,
							QbListID:        ret.ListID,
							QbEditSequence:  ret.EditSequence,
							QbFullName:      ret.FullName,
						}
						config.DB.Create(&mapping)
					}
				} else if queueItem.ActionType == "ITEM_NON_INVENTORY_ADD" && qbxml.QBXMLMsgsRs.ItemNonInventoryAddRs != nil {
					if qbxml.QBXMLMsgsRs.ItemNonInventoryAddRs.StatusCode == "0" {
						ret := qbxml.QBXMLMsgsRs.ItemNonInventoryAddRs.ItemNonInventoryRet
						mapping := models.QBItemMapping{
							ItemNumber:      queueItem.ItemNumber,
							QBProfileNumber: targetProfileNumber,
							QbListID:        ret.ListID,
							QbEditSequence:  ret.EditSequence,
							QbFullName:      ret.FullName,
						}
						config.DB.Create(&mapping)
					}
				}
			}
		} else {
			// QB returned a non-zero statusCode — extract the actual error from the XML
			var qbxml QBXMLResponse
			xml.Unmarshal([]byte(response), &qbxml)

			errMsg := message // fallback to SOAP message
			if qbxml.QBXMLMsgsRs.SalesReceiptAddRs != nil && qbxml.QBXMLMsgsRs.SalesReceiptAddRs.StatusMessage != "" {
				errMsg = fmt.Sprintf("[%s] %s", qbxml.QBXMLMsgsRs.SalesReceiptAddRs.StatusCode, qbxml.QBXMLMsgsRs.SalesReceiptAddRs.StatusMessage)
			} else if qbxml.QBXMLMsgsRs.ItemServiceAddRs != nil && qbxml.QBXMLMsgsRs.ItemServiceAddRs.StatusMessage != "" {
				errMsg = fmt.Sprintf("[%s] %s", qbxml.QBXMLMsgsRs.ItemServiceAddRs.StatusCode, qbxml.QBXMLMsgsRs.ItemServiceAddRs.StatusMessage)
			} else if qbxml.QBXMLMsgsRs.CustomerAddRs != nil && qbxml.QBXMLMsgsRs.CustomerAddRs.StatusMessage != "" {
				errMsg = fmt.Sprintf("[%s] %s", qbxml.QBXMLMsgsRs.CustomerAddRs.StatusCode, qbxml.QBXMLMsgsRs.CustomerAddRs.StatusMessage)
			}
			queueItem.ErrorMessage = errMsg
			queueItem.RetryCount++
			if queueItem.RetryCount < queueItem.MaxRetries {
				queueItem.Status = "QUEUED"
			} else {
				queueItem.Status = "FAILED"
			}
		}
		config.DB.Save(&queueItem)
	}

	// Check if there is more work to do for this profile
	var count int64
	remainingQuery := config.DB.Model(&models.QBSyncQueue{}).Where("status = ?", "QUEUED")
	if isDefaultProfile(profileNumber) {
		remainingQuery = remainingQuery.Where("qb_profile_number = ? OR qb_profile_number = '' OR qb_profile_number IS NULL", profileNumber)
	} else {
		remainingQuery = remainingQuery.Where("qb_profile_number = ?", profileNumber)
	}
	remainingQuery.Count(&count)

	if count == 0 {
		return "100" // 100% complete
	}

	// Rough estimate of completion, e.g., "50"
	return "50"
}

func CloseConnection(ticket string) {
	fmt.Println("Closing QBWC session for ticket:", ticket)
	activeSessions.Delete(ticket)
}
