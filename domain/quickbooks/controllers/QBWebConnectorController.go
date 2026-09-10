package controllers

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	authModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	campusModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	invoiceModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/models"
	itemModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/item/models"
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

	// Fix 3: Duplicate check — prevent queueing the same receipt twice
	var existing qbModels.QBSyncQueue
	if err := config.DB.Where("transaction_number = ? AND status IN ?", req.TransactionNumber, []string{"QUEUED", "PROCESSING", "SUCCESS"}).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "This transaction has already been queued or synced to QuickBooks (status: " + existing.Status + ")"})
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

	// Require customer to be exported to QB first
	userNumber := itemTx.Transaction.Statement.User.Number
	studentName := itemTx.Transaction.Statement.User.FirstName + " " + itemTx.Transaction.Statement.User.LastName
	var customerRef string
	if userNumber != "" {
		var mapping qbModels.QBCustomerMapping
		if err := config.DB.Where("user_number = ?", userNumber).First(&mapping).Error; err == nil {
			customerRef = mapping.QbFullName
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Student '" + studentName + "' has not been exported to QuickBooks yet. Please export the student first from the Enrollment page."})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot identify student for this transaction."})
		return
	}

	var itemCode string
	var itemMapping qbModels.QBItemMapping
	if err := config.DB.Where("item_number = ?", itemTx.InvoiceItem.Item.Number).First(&itemMapping).Error; err == nil {
		itemCode = itemMapping.QbFullName
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Item '" + itemTx.InvoiceItem.Item.Name + "' has not been exported to QuickBooks yet. Please export the item first from the Services page."})
		return
	}

	dateStr := itemTx.Transaction.Date.Format("2006-01-02")
	amount := fmt.Sprintf("%.2f", itemTx.AmountPaidOff)

	depositAccount := "Undeposited Funds"
	
	if itemTx.Transaction.DepositAccountNumber != "" {
		var qbDepositAccount qbModels.QBDepositAccount
		if err := config.DB.Where("number = ?", itemTx.Transaction.DepositAccountNumber).First(&qbDepositAccount).Error; err == nil && qbDepositAccount.Name != "" {
			depositAccount = qbDepositAccount.Name
		}
	} else {
		var qbSettings qbModels.QBSettings
		if err := config.DB.First(&qbSettings).Error; err == nil && qbSettings.DepositToAccount != "" {
			depositAccount = qbSettings.DepositToAccount
		}
	}

	qbxml := `<SalesReceiptAddRq>
	<SalesReceiptAdd>
		<CustomerRef>
			<FullName>` + escapeXML(customerRef) + `</FullName>
		</CustomerRef>
		<TxnDate>` + dateStr + `</TxnDate>
		<RefNumber>` + escapeXML(itemTx.ItemReceiptNumber) + `</RefNumber>
		<DepositToAccountRef>
			<FullName>` + escapeXML(depositAccount) + `</FullName>
		</DepositToAccountRef>
		<SalesReceiptLineAdd>
			<ItemRef>
				<FullName>` + escapeXML(itemCode) + `</FullName>
			</ItemRef>
			<Rate>` + amount + `</Rate>
			<Amount>` + amount + `</Amount>
		</SalesReceiptLineAdd>
	</SalesReceiptAdd>
</SalesReceiptAddRq>`

	queue := qbModels.QBSyncQueue{
		Number:            uuid.NewV4().String(),
		TransactionNumber: req.TransactionNumber,
		ActionType:        "SALES_RECEIPT_ADD",
		Status:            "QUEUED",
		QbXmlPayload:      qbxml,
		UserNumber:        userNumber,
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

	// Fetch student for ADN and campus
	var student studentModels.Student
	config.DB.Where("user_number = ?", req.UserNumber).First(&student)

	fullName := user.FirstName + " " + user.LastName
	if student.ADN != "" {
		fullName = fullName + " - " + student.ADN
	}

	// Fetch campus name to use as CompanyName in QB
	companyName := ""
	if student.PreferredCampusNumber != "" {
		var campus campusModels.Campus
		if err := config.DB.Where("number = ?", student.PreferredCampusNumber).First(&campus).Error; err == nil {
			companyName = campus.Name
		}
	}

	companyNameXML := ""
	classRefXML := ""
	if companyName != "" {
		companyNameXML = `
			<CompanyName>` + escapeXML(companyName) + `</CompanyName>`
		classRefXML = `
			<ClassRef>
				<FullName>` + escapeXML(companyName) + `</FullName>
			</ClassRef>`
	}

	qbxml := `<CustomerAddRq>
		<CustomerAdd>
			<Name>` + escapeXML(fullName) + `</Name>` + classRefXML + companyNameXML + `
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
	CampusName string `json:"campus_name"`
}

func GetSyncedCustomers(c *gin.Context) {
	schoolNumber := c.Query("school_number")
	campusNumber := c.Query("campus_number")

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

		// Apply school filter
		if schoolNumber != "" && student.SchoolNumber != schoolNumber {
			continue
		}
		// Apply campus filter
		if campusNumber != "" && student.PreferredCampusNumber != campusNumber {
			continue
		}

		// Resolve campus name for display
		campusName := ""
		if student.PreferredCampusNumber != "" {
			var campus campusModels.Campus
			if err := config.DB.Where("number = ?", student.PreferredCampusNumber).First(&campus).Error; err == nil {
				campusName = campus.Name
			}
		}

		fullName := user.FirstName + " " + user.SecondName + " " + user.LastName

		responses = append(responses, SyncedCustomerResponse{
			UserNumber: mapping.UserNumber,
			FullName:   strings.TrimSpace(fullName),
			ADN:        student.ADN,
			QBListID:   mapping.QbListID,
			CampusName: campusName,
		})
	}

	c.JSON(http.StatusOK, responses)
}

type QueueItemRequest struct {
	ItemNumber string `json:"item_number" binding:"required"`
	AccountID  string `json:"account_id"`
}

func QueueItemAdd(c *gin.Context) {
	var req QueueItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if already mapped
	var mapping qbModels.QBItemMapping
	if err := config.DB.Where("item_number = ?", req.ItemNumber).First(&mapping).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "This item is already exported to QuickBooks"})
		return
	}

	// Check if already queued
	var existing qbModels.QBSyncQueue
	if err := config.DB.Where("item_number = ? AND status IN ?", req.ItemNumber, []string{"QUEUED", "PROCESSING", "SUCCESS"}).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "This item has already been queued for QuickBooks (status: " + existing.Status + ")"})
		return
	}

	var item itemModels.Item
	if err := config.DB.Where("number = ?", req.ItemNumber).First(&item).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	itemCode := item.Name
	if item.Code != 0 {
		itemCode = fmt.Sprintf("%d", item.Code)
	}

	// Always use the item Name as the description in QuickBooks
	desc := item.Name
	price := fmt.Sprintf("%.2f", item.Cost)

	// Load income account from request or fallback to settings
	accountName := "Service Income"
	if req.AccountID != "" {
		var chartOfAccount qbModels.QBChartOfAccount
		if err := config.DB.Where("number = ?", req.AccountID).First(&chartOfAccount).Error; err == nil {
			accountName = chartOfAccount.Name
		}
	} else {
		var qbSettingsForItem qbModels.QBSettings
		if err := config.DB.First(&qbSettingsForItem).Error; err == nil && qbSettingsForItem.IncomeAccount != "" {
			accountName = qbSettingsForItem.IncomeAccount
		}
	}

	qbxml := `<ItemServiceAddRq>
	<ItemServiceAdd>
		<Name>` + escapeXML(itemCode) + `</Name>
		<IsActive>true</IsActive>
		<SalesOrPurchase>
			<Desc>` + escapeXML(desc) + `</Desc>
			<Price>` + price + `</Price>
			<AccountRef>
				<FullName>` + escapeXML(accountName) + `</FullName>
			</AccountRef>
		</SalesOrPurchase>
	</ItemServiceAdd>
</ItemServiceAddRq>`

	queue := qbModels.QBSyncQueue{
		Number:       uuid.NewV4().String(),
		ActionType:   "ITEM_SERVICE_ADD",
		Status:       "QUEUED",
		QbXmlPayload: qbxml,
		ItemNumber:   req.ItemNumber,
	}

	if err := config.DB.Create(&queue).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to queue item addition"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item addition queued for QuickBooks sync successfully."})
}

// QueueProductAdd handles POST /main/quickbooks/products.
// It queues an Inventory Item or Non-Inventory Item in QuickBooks.
type QueueProductRequest struct {
	ItemNumber  string `json:"item_number" binding:"required"`
	AccountID   string `json:"account_id"`
	ProductType string `json:"product_type" binding:"required"` // "Inventory" | "NonInventory"
}

func QueueProductAdd(c *gin.Context) {
	var req QueueProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate product type
	if req.ProductType != "Inventory" && req.ProductType != "NonInventory" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product_type must be 'Inventory' or 'NonInventory'"})
		return
	}

	// Determine the QB action type string for this product type
	var actionType string
	switch req.ProductType {
	case "Inventory":
		actionType = "ITEM_INVENTORY_ADD"
	case "NonInventory":
		actionType = "ITEM_NON_INVENTORY_ADD"
	}

	// Check if already mapped in QB
	var mapping qbModels.QBItemMapping
	if err := config.DB.Where("item_number = ?", req.ItemNumber).First(&mapping).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "This product is already exported to QuickBooks"})
		return
	}

	// Check if already in the queue with a non-failed status
	var existing qbModels.QBSyncQueue
	if err := config.DB.Where(
		"item_number = ? AND status IN ?",
		req.ItemNumber, []string{"QUEUED", "PROCESSING", "SUCCESS"},
	).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "This product is already queued for QuickBooks (status: " + existing.Status + ")"})
		return
	}

	// Load the item record
	var item itemModels.Item
	if err := config.DB.Where("number = ?", req.ItemNumber).First(&item).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Derive the QB item name (use numeric code when available for uniqueness)
	itemCode := item.Name
	if item.Code != 0 {
		itemCode = fmt.Sprintf("%d", item.Code)
	}
	desc := item.Name
	price := fmt.Sprintf("%.2f", item.Cost)

	// Resolve income account name
	accountName := "Sales"
	if req.AccountID != "" {
		var chartOfAccount qbModels.QBChartOfAccount
		if err := config.DB.Where("number = ?", req.AccountID).First(&chartOfAccount).Error; err == nil {
			accountName = chartOfAccount.Name
		}
	} else {
		var qbSettings qbModels.QBSettings
		if err := config.DB.First(&qbSettings).Error; err == nil && qbSettings.IncomeAccount != "" {
			accountName = qbSettings.IncomeAccount
		}
	}

	// Build the QB XML payload based on product type
	var qbxml string
	switch req.ProductType {
	case "Inventory":
		qbxml = `<ItemInventoryAddRq>
	<ItemInventoryAdd>
		<Name>` + escapeXML(itemCode) + `</Name>
		<IsActive>true</IsActive>
		<SalesDesc>` + escapeXML(desc) + `</SalesDesc>
		<SalesPrice>` + price + `</SalesPrice>
		<IncomeAccountRef>
			<FullName>` + escapeXML(accountName) + `</FullName>
		</IncomeAccountRef>
		<AssetAccountRef>
			<FullName>Inventory Asset</FullName>
		</AssetAccountRef>
		<COGSAccountRef>
			<FullName>Cost of Goods Sold</FullName>
		</COGSAccountRef>
	</ItemInventoryAdd>
</ItemInventoryAddRq>`

	case "NonInventory":
		qbxml = `<ItemNonInventoryAddRq>
	<ItemNonInventoryAdd>
		<Name>` + escapeXML(itemCode) + `</Name>
		<IsActive>true</IsActive>
		<SalesOrPurchase>
			<Desc>` + escapeXML(desc) + `</Desc>
			<Price>` + price + `</Price>
			<AccountRef>
				<FullName>` + escapeXML(accountName) + `</FullName>
			</AccountRef>
		</SalesOrPurchase>
	</ItemNonInventoryAdd>
</ItemNonInventoryAddRq>`
	}

	queue := qbModels.QBSyncQueue{
		Number:       uuid.NewV4().String(),
		ActionType:   actionType,
		Status:       "QUEUED",
		QbXmlPayload: qbxml,
		ItemNumber:   req.ItemNumber,
	}

	if err := config.DB.Create(&queue).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to queue product addition"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product queued for QuickBooks export successfully."})
}

type SyncedItemResponse struct {
	ItemNumber string  `json:"item_number"`
	ItemName   string  `json:"item_name"`
	QBListID   string  `json:"qb_list_id"`
	QBFullName string  `json:"qb_full_name"`
	Cost       float64 `json:"cost"`
}

func GetSyncedItems(c *gin.Context) {
	var mappings []qbModels.QBItemMapping
	if err := config.DB.Find(&mappings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch synced items"})
		return
	}

	var responses []SyncedItemResponse
	for _, mapping := range mappings {
		var item itemModels.Item
		config.DB.Where("number = ?", mapping.ItemNumber).First(&item)

		responses = append(responses, SyncedItemResponse{
			ItemNumber: mapping.ItemNumber,
			ItemName:   item.Name,
			QBListID:   mapping.QbListID,
			QBFullName: mapping.QbFullName,
			Cost:       item.Cost,
		})
	}

	c.JSON(http.StatusOK, responses)
}

func GetQueueErrors(c *gin.Context) {
	var items []qbModels.QBSyncQueue
	if err := config.DB.Where("status = ?", "FAILED").Order("created_at desc").Limit(10).Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch errors"})
		return
	}
	c.JSON(http.StatusOK, items)
}

// GetSyncedReceipts returns ALL sales receipt transactions with any status
// so the UI can show QUEUED, PROCESSING, SUCCESS, and FAILED badges.
func GetSyncedReceipts(c *gin.Context) {
	var items []qbModels.QBSyncQueue
	if err := config.DB.
		Where("action_type = ?", "SALES_RECEIPT_ADD").
		Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch synced receipts"})
		return
	}
	c.JSON(http.StatusOK, items)
}

type QueueStatusEntry struct {
	QueueNumber  string `json:"queue_number"`
	Status       string `json:"status"`
	ErrorMessage string `json:"error_message"`
}

// GetCustomerQueueStatus returns a map of user_number → queue status for CUSTOMER_ADD entries.
func GetCustomerQueueStatus(c *gin.Context) {
	var items []qbModels.QBSyncQueue
	config.DB.Where("action_type = ?", "CUSTOMER_ADD").Find(&items)

	result := make(map[string]QueueStatusEntry)
	for _, item := range items {
		result[item.UserNumber] = QueueStatusEntry{
			QueueNumber:  item.Number,
			Status:       item.Status,
			ErrorMessage: item.ErrorMessage,
		}
	}
	c.JSON(http.StatusOK, result)
}

// GetItemQueueStatus returns a map of item_number → queue status.
// Covers services (ITEM_SERVICE_ADD) and both product types (ITEM_INVENTORY_ADD, ITEM_NON_INVENTORY_ADD).
func GetItemQueueStatus(c *gin.Context) {
	var items []qbModels.QBSyncQueue
	config.DB.Where(
		"action_type IN ?",
		[]string{"ITEM_SERVICE_ADD", "ITEM_INVENTORY_ADD", "ITEM_NON_INVENTORY_ADD"},
	).Find(&items)

	result := make(map[string]QueueStatusEntry)
	for _, item := range items {
		result[item.ItemNumber] = QueueStatusEntry{
			QueueNumber:  item.Number,
			Status:       item.Status,
			ErrorMessage: item.ErrorMessage,
		}
	}
	c.JSON(http.StatusOK, result)
}

// DeleteQueueEntry deletes a QUEUED or FAILED queue entry.
// SUCCESS and PROCESSING entries cannot be deleted.
func DeleteQueueEntry(c *gin.Context) {
	number := c.Param("number")

	var item qbModels.QBSyncQueue
	if err := config.DB.Where("number = ?", number).First(&item).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Queue entry not found"})
		return
	}

	if item.Status == "SUCCESS" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete a successfully synced entry — it already exists in QuickBooks."})
		return
	}
	if item.Status == "PROCESSING" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete an entry that is currently being processed by the Web Connector."})
		return
	}

	// Hard delete the queue entry so the entity can be re-exported
	if err := config.DB.Unscoped().Delete(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete queue entry"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Queue entry deleted. You can now re-export this record."})
}
