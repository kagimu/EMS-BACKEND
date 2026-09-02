package controllers

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	authModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	invoiceModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/models"
	qbModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/quickbooks/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/quickbooks/services"
	studentModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

// Simplified SOAP Envelope
type Envelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    Body     `xml:"Body"`
}

type Body struct {
	Authenticate       *Authenticate       `xml:"authenticate"`
	SendRequestXML     *SendRequestXML     `xml:"sendRequestXML"`
	ReceiveResponseXML *ReceiveResponseXML `xml:"receiveResponseXML"`
	CloseConnection    *CloseConnection    `xml:"closeConnection"`
	ServerVersion      *ServerVersion      `xml:"serverVersion"`
	ClientVersion      *ClientVersion      `xml:"clientVersion"`
	ConnectionError    *ConnectionError    `xml:"connectionError"`
	GetLastError       *GetLastError       `xml:"getLastError"`
}

type Authenticate struct {
	UserName string `xml:"strUserName"`
	Password string `xml:"strPassword"`
}

type SendRequestXML struct {
	Ticket string `xml:"ticket"`
}

type ReceiveResponseXML struct {
	Ticket   string `xml:"ticket"`
	Response string `xml:"response"`
	HResult  string `xml:"hresult"`
	Message  string `xml:"message"`
}

type CloseConnection struct {
	Ticket string `xml:"ticket"`
}

type ServerVersion struct{}
type ClientVersion struct{}
type ConnectionError struct{}
type GetLastError struct {
	Ticket string `xml:"ticket"`
}

func HandleQBWC(c *gin.Context) {
	if c.Request.Method == "GET" {
		c.String(http.StatusOK, "QuickBooks Web Connector SOAP Endpoint is Active.")
		return
	}

	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid request")
		return
	}

	var env Envelope
	if err := xml.Unmarshal(bodyBytes, &env); err != nil {
		// Log the error if necessary
		c.String(http.StatusBadRequest, "Invalid XML")
		return
	}

	var responseXML string

	if env.Body.ServerVersion != nil {
		responseXML = `<serverVersionResponse xmlns="http://developer.intuit.com/"><serverVersionResult>1.0</serverVersionResult></serverVersionResponse>`
	} else if env.Body.ClientVersion != nil {
		responseXML = `<clientVersionResponse xmlns="http://developer.intuit.com/"><clientVersionResult></clientVersionResult></clientVersionResponse>`
	} else if env.Body.Authenticate != nil {
		res := services.Authenticate(env.Body.Authenticate.UserName, env.Body.Authenticate.Password)
		responseXML = `<authenticateResponse xmlns="http://developer.intuit.com/"><authenticateResult><string>` + res.Ticket + `</string><string>` + res.Status + `</string></authenticateResult></authenticateResponse>`
	} else if env.Body.SendRequestXML != nil {
		reqXML := services.SendRequestXML(env.Body.SendRequestXML.Ticket)
		responseXML = `<sendRequestXMLResponse xmlns="http://developer.intuit.com/"><sendRequestXMLResult>` + escapeXML(reqXML) + `</sendRequestXMLResult></sendRequestXMLResponse>`
	} else if env.Body.ReceiveResponseXML != nil {
		progress := services.ReceiveResponseXML(env.Body.ReceiveResponseXML.Ticket, env.Body.ReceiveResponseXML.Response, env.Body.ReceiveResponseXML.HResult, env.Body.ReceiveResponseXML.Message)
		responseXML = `<receiveResponseXMLResponse xmlns="http://developer.intuit.com/"><receiveResponseXMLResult>` + progress + `</receiveResponseXMLResult></receiveResponseXMLResponse>`
	} else if env.Body.ConnectionError != nil {
		responseXML = `<connectionErrorResponse xmlns="http://developer.intuit.com/"><connectionErrorResult>done</connectionErrorResult></connectionErrorResponse>`
	} else if env.Body.GetLastError != nil {
		responseXML = `<getLastErrorResponse xmlns="http://developer.intuit.com/"><getLastErrorResult>Error occurred</getLastErrorResult></getLastErrorResponse>`
	} else if env.Body.CloseConnection != nil {
		services.CloseConnection(env.Body.CloseConnection.Ticket)
		responseXML = `<closeConnectionResponse xmlns="http://developer.intuit.com/"><closeConnectionResult>OK</closeConnectionResult></closeConnectionResponse>`
	} else {
		c.String(http.StatusBadRequest, "Unknown SOAP Action")
		return
	}

	// Wrap in SOAP Envelope
	finalResponse := `<?xml version="1.0" encoding="utf-8"?><soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema"><soap:Body>` + responseXML + `</soap:Body></soap:Envelope>`
	
	c.Data(http.StatusOK, "text/xml; charset=utf-8", []byte(finalResponse))
}

func escapeXML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}

func DownloadQWC(c *gin.Context) {
	// Determine protocol (default to https since QB Web Connector requires it in production unless using localhost)
	proto := c.GetHeader("X-Forwarded-Proto")
	if proto == "" {
		proto = "https"
		if c.Request.TLS == nil && strings.Contains(c.Request.Host, "localhost") {
			proto = "http"
		}
	}
	baseURL := fmt.Sprintf("%s://%s", proto, c.Request.Host)

	// A minimal standard qwc file template
	qwc := fmt.Sprintf(`<?xml version="1.0"?>
<QBWCXML>
   <AppName>WholesomeEd QB Sync</AppName>
   <AppID></AppID>
   <AppURL>%s/api/quickbooks/soap</AppURL>
   <AppDescription>Sync payments to QuickBooks Desktop</AppDescription>
   <AppSupport>%s/support</AppSupport>
   <UserName>admin</UserName>
   <OwnerID>{90A44FB7-33D9-4815-AC85-BC87A7E7D1EB}</OwnerID>
   <FileID>{57F3B9B1-86F1-4fcc-B1EE-566DE1813D20}</FileID>
   <QBType>QBFS</QBType>
   <Scheduler>
      <RunEveryNMinutes>5</RunEveryNMinutes>
   </Scheduler>
   <IsReadOnly>false</IsReadOnly>
</QBWCXML>`, baseURL, baseURL)

	c.Header("Content-Disposition", `attachment; filename="wholesome_ed.qwc"`)
	c.Data(http.StatusOK, "application/x-qbwc", []byte(qwc))
}

type QueueTransactionRequest struct {
	TransactionNumber string `json:"transaction_number" binding:"required"`
}

func QueueTransaction(c *gin.Context) {
	var req QueueTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var itemTx invoiceModel.InvoiceItemTransaction
	if err := config.DB.Preload("Transaction.Statement.User").
		Preload("InvoiceItem.Item").
		Where("item_receipt_number = ?", req.TransactionNumber).
		First(&itemTx).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	// Generate SalesReceiptAddRq qbXML
	// Mapping:
	// CustomerRef -> Student Name
	// TxnDate -> Transaction Date
	// RefNumber -> Item Receipt Number
	// ItemRef -> Item Code
	// Rate -> Amount Paid Off
	// DepositToAccountRef -> "Undeposited Funds" (default mapping)

	studentName := "Unknown Student"
	if itemTx.Transaction.Statement.User.Number != "" {
		studentName = itemTx.Transaction.Statement.User.FirstName + " " + itemTx.Transaction.Statement.User.LastName
	}
	
	itemCode := fmt.Sprintf("%d", itemTx.InvoiceItem.Item.Code)
	if itemTx.InvoiceItem.Item.Code == 0 {
		itemCode = itemTx.InvoiceItem.Item.Name // fallback
	}

	dateStr := itemTx.Transaction.Date.Format("2006-01-02")
	amount := fmt.Sprintf("%.2f", itemTx.AmountPaidOff)

	qbxml := `
	<SalesReceiptAdd>
		<CustomerRef>
			<FullName>` + escapeXML(studentName) + `</FullName>
		</CustomerRef>
		<TxnDate>` + dateStr + `</TxnDate>
		<RefNumber>` + escapeXML(itemTx.ItemReceiptNumber) + `</RefNumber>
		<DepositToAccountRef>
			<FullName>Undeposited Funds</FullName>
		</DepositToAccountRef>
		<SalesReceiptLineAdd>
			<ItemRef>
				<FullName>` + escapeXML(itemCode) + `</FullName>
			</ItemRef>
			<Rate>` + amount + `</Rate>
			<Amount>` + amount + `</Amount>
		</SalesReceiptLineAdd>
	</SalesReceiptAdd>`

	queue := qbModels.QBSyncQueue{
		Number:            uuid.NewV4().String(),
		TransactionNumber: req.TransactionNumber,
		Status:            "QUEUED",
		QbXmlPayload:      qbxml,
	}

	if err := config.DB.Create(&queue).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to queue transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction queued for QuickBooks sync successfully."})
}

type QueueCustomerRequest struct {
	UserNumber string `json:"user_number" binding:"required"`
}

func QueueCustomerAdd(c *gin.Context) {
	var req QueueCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if already mapped
	var mapping qbModels.QBCustomerMapping
	if err := config.DB.Where("user_number = ?", req.UserNumber).First(&mapping).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "This user is already exported to QuickBooks"})
		return
	}

	var user authModels.User
	if err := config.DB.Where("number = ?", req.UserNumber).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Fetch student for ADN
	var student studentModels.Student
	config.DB.Where("user_number = ?", req.UserNumber).First(&student)

	fullName := user.FirstName + " " + user.LastName
	if student.ADN != "" {
		fullName = fullName + " - " + student.ADN
	}
	
	qbxml := `<CustomerAddRq>
		<CustomerAdd>
			<Name>` + escapeXML(fullName) + `</Name>
			<FirstName>` + escapeXML(user.FirstName) + `</FirstName>
			<LastName>` + escapeXML(user.LastName) + `</LastName>
		</CustomerAdd>
	</CustomerAddRq>`

	queue := qbModels.QBSyncQueue{
		Number:       uuid.NewV4().String(),
		ActionType:   "CUSTOMER_ADD",
		Status:       "QUEUED",
		QbXmlPayload: qbxml,
		UserNumber:   req.UserNumber,
	}

	if err := config.DB.Create(&queue).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to queue customer addition"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Customer addition queued for QuickBooks sync successfully."})
}

func QueueCustomerQuery(c *gin.Context) {
	qbxml := `<CustomerQueryRq></CustomerQueryRq>`

	queue := qbModels.QBSyncQueue{
		Number:       uuid.NewV4().String(),
		ActionType:   "CUSTOMER_QUERY",
		Status:       "QUEUED",
		QbXmlPayload: qbxml,
	}

	if err := config.DB.Create(&queue).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to queue customer query"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Customer query queued for QuickBooks sync successfully."})
}

type SyncedCustomerResponse struct {
	UserNumber string `json:"user_number"`
	FullName   string `json:"full_name"`
	ADN        string `json:"adn"`
	QBListID   string `json:"qb_list_id"`
}

func GetSyncedCustomers(c *gin.Context) {
	var mappings []qbModels.QBCustomerMapping
	if err := config.DB.Find(&mappings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch synced customers"})
		return
	}

	var responses []SyncedCustomerResponse
	for _, mapping := range mappings {
		var user authModels.User
		var student studentModels.Student
		
		config.DB.Where("number = ?", mapping.UserNumber).First(&user)
		config.DB.Where("user_number = ?", mapping.UserNumber).First(&student)

		fullName := user.FirstName + " " + user.SecondName + " " + user.LastName
		
		responses = append(responses, SyncedCustomerResponse{
			UserNumber: mapping.UserNumber,
			FullName:   strings.TrimSpace(fullName),
			ADN:        student.ADN,
			QBListID:   mapping.QbListID,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": responses})
}
