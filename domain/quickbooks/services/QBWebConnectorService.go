package services

import (
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/quickbooks/models"
	uuid "github.com/satori/go.uuid"
)

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
		if hresult != "" && hresult != "0x00000000" {
			// QB itself returned a COM-level error (e.g. bad XML)
			queueItem.Status = "FAILED"
			queueItem.ErrorMessage = message
		} else if strings.Contains(response, `statusCode="0"`) {
			queueItem.Status = "SUCCESS"

			if queueItem.ActionType == "CUSTOMER_ADD" || queueItem.ActionType == "CUSTOMER_QUERY" || queueItem.ActionType == "SALES_RECEIPT_ADD" || queueItem.ActionType == "ITEM_SERVICE_ADD" || queueItem.ActionType == "ITEM_INVENTORY_ADD" || queueItem.ActionType == "ITEM_NON_INVENTORY_ADD" {
				var qbxml QBXMLResponse
				xml.Unmarshal([]byte(response), &qbxml)

				if queueItem.ActionType == "CUSTOMER_ADD" && qbxml.QBXMLMsgsRs.CustomerAddRs != nil {
					ret := qbxml.QBXMLMsgsRs.CustomerAddRs.CustomerRet
					mapping := models.QBCustomerMapping{
						UserNumber:     queueItem.UserNumber,
						QbListID:       ret.ListID,
						QbEditSequence: ret.EditSequence,
						QbFullName:     ret.FullName,
					}
					config.DB.Create(&mapping)
				} else if queueItem.ActionType == "CUSTOMER_QUERY" && qbxml.QBXMLMsgsRs.CustomerQueryRs != nil {
					for _, ret := range qbxml.QBXMLMsgsRs.CustomerQueryRs.CustomerRet {
						var mapping models.QBCustomerMapping
						err := config.DB.Where("qb_list_id = ?", ret.ListID).First(&mapping).Error
						if err == nil {
							// Only update if it already exists. Do not create missing students from QB.
							mapping.QbEditSequence = ret.EditSequence
							mapping.QbFullName = ret.FullName
							config.DB.Save(&mapping)
						}
					}
				} else if queueItem.ActionType == "SALES_RECEIPT_ADD" && qbxml.QBXMLMsgsRs.SalesReceiptAddRs != nil {
					// Fix 4: Store QB TxnID for tracking
					queueItem.QbTxnID = qbxml.QBXMLMsgsRs.SalesReceiptAddRs.SalesReceiptRet.TxnID
				} else if queueItem.ActionType == "ITEM_SERVICE_ADD" && qbxml.QBXMLMsgsRs.ItemServiceAddRs != nil {
					ret := qbxml.QBXMLMsgsRs.ItemServiceAddRs.ItemServiceRet
					mapping := models.QBItemMapping{
						ItemNumber:     queueItem.ItemNumber,
						QbListID:       ret.ListID,
						QbEditSequence: ret.EditSequence,
						QbFullName:     ret.FullName,
					}
					config.DB.Create(&mapping)
				} else if queueItem.ActionType == "ITEM_INVENTORY_ADD" && qbxml.QBXMLMsgsRs.ItemInventoryAddRs != nil {
					ret := qbxml.QBXMLMsgsRs.ItemInventoryAddRs.ItemInventoryRet
					mapping := models.QBItemMapping{
						ItemNumber:     queueItem.ItemNumber,
						QbListID:       ret.ListID,
						QbEditSequence: ret.EditSequence,
						QbFullName:     ret.FullName,
					}
					config.DB.Create(&mapping)
				} else if queueItem.ActionType == "ITEM_NON_INVENTORY_ADD" && qbxml.QBXMLMsgsRs.ItemNonInventoryAddRs != nil {
					ret := qbxml.QBXMLMsgsRs.ItemNonInventoryAddRs.ItemNonInventoryRet
					mapping := models.QBItemMapping{
						ItemNumber:     queueItem.ItemNumber,
						QbListID:       ret.ListID,
						QbEditSequence: ret.EditSequence,
						QbFullName:     ret.FullName,
					}
					config.DB.Create(&mapping)
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
			queueItem.Status = "FAILED"
			queueItem.ErrorMessage = errMsg
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
