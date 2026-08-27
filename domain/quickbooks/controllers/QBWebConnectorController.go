package controllers

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	invoiceModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/models"
	qbModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/quickbooks/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/quickbooks/services"
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
		responseXML = `<serverVersionResult>1.0</serverVersionResult>`
	} else if env.Body.ClientVersion != nil {
		responseXML = `<clientVersionResult></clientVersionResult>`
	} else if env.Body.Authenticate != nil {
		res := services.Authenticate(env.Body.Authenticate.UserName, env.Body.Authenticate.Password)
		responseXML = `<authenticateResult><string>` + res.Ticket + `</string><string>` + res.Status + `</string></authenticateResult>`
	} else if env.Body.SendRequestXML != nil {
		reqXML := services.SendRequestXML(env.Body.SendRequestXML.Ticket)
		responseXML = `<sendRequestXMLResult>` + escapeXML(reqXML) + `</sendRequestXMLResult>`
	} else if env.Body.ReceiveResponseXML != nil {
		progress := services.ReceiveResponseXML(env.Body.ReceiveResponseXML.Ticket, env.Body.ReceiveResponseXML.Response, env.Body.ReceiveResponseXML.HResult, env.Body.ReceiveResponseXML.Message)
		responseXML = `<receiveResponseXMLResult>` + progress + `</receiveResponseXMLResult>`
	} else if env.Body.ConnectionError != nil {
		responseXML = `<connectionErrorResult>done</connectionErrorResult>`
	} else if env.Body.GetLastError != nil {
		responseXML = `<getLastErrorResult>Error occurred</getLastErrorResult>`
	} else if env.Body.CloseConnection != nil {
		services.CloseConnection(env.Body.CloseConnection.Ticket)
		responseXML = `<closeConnectionResult>OK</closeConnectionResult>`
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
