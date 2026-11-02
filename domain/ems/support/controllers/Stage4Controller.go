package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	supportRequests "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/support/api/requests"
	supportResponses "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/support/api/responses"
	supportMiddleware "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/support/middleware"
	supportModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/support/models"
	supportServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/support/services"
	"github.com/gin-gonic/gin"
)

func staffError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, supportServices.ErrSupportStaffInvalid):
		c.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
	case errors.Is(err, supportServices.ErrAttachmentInvalid), errors.Is(err, supportServices.ErrAttachmentTooLarge):
		c.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
	case errors.Is(err, supportServices.ErrAttachmentNotFound):
		c.JSON(http.StatusNotFound, commons.Feedback("Attachment not found"))
	case errors.Is(err, supportServices.ErrTicketNotFound):
		c.JSON(http.StatusNotFound, commons.Feedback("Support ticket not found"))
	case errors.Is(err, supportServices.ErrTicketAuthorisation):
		c.JSON(http.StatusForbidden, commons.Feedback("You cannot access this support ticket"))
	case errors.Is(err, supportServices.ErrTicketTransition), errors.Is(err, supportServices.ErrEscalationReason):
		c.JSON(http.StatusConflict, commons.Feedback(err.Error()))
	default:
		c.JSON(http.StatusBadRequest, commons.Feedback("The support operation could not be completed"))
	}
}

func supportPage(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.Query("page"))
	perPage, _ := strconv.Atoi(c.Query("per_page"))
	return page, perPage
}

func OwnerStaff(c *gin.Context) {
	staff, err := supportServices.ListSupportStaff()
	if err != nil {
		staffError(c, err)
		return
	}
	eligible, _ := supportServices.EligibleSupportStaff()
	c.JSON(http.StatusOK, gin.H{"message": "Support staff fetched", "data": gin.H{"items": staff, "eligible": eligible}})
}

func OwnerAddStaff(c *gin.Context) {
	actorNumber, ok := actor(c)
	if !ok {
		return
	}
	var body struct {
		UserNumber string `json:"user_number"`
		Role       string `json:"role"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("Invalid request body"))
		return
	}
	staff, err := supportServices.UpsertSupportStaff(actorNumber, body.UserNumber, body.Role)
	if err != nil {
		staffError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Support staff saved", "data": staff})
}

func OwnerStaffStatus(c *gin.Context) {
	actorNumber, ok := actor(c)
	if !ok {
		return
	}
	var body struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("Invalid request body"))
		return
	}
	staff, err := supportServices.SetSupportStaffStatus(actorNumber, c.Param("user_number"), body.Status)
	if err != nil {
		staffError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Support staff status updated", "data": staff})
}

func OwnerRemoveStaff(c *gin.Context) {
	if err := supportServices.RemoveSupportStaff(c.Param("user_number")); err != nil {
		staffError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Support staff removed"})
}

func OwnerWorkload(c *gin.Context) {
	workload, err := supportServices.SupportWorkload()
	if err != nil {
		staffError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Support workload fetched", "data": workload})
}

func OwnerSLATargets(c *gin.Context) {
	targets, err := supportServices.ListSLATargets()
	if err != nil {
		staffError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "SLA targets fetched", "data": targets})
}

func OwnerSaveSLATargets(c *gin.Context) {
	actorNumber, ok := actor(c)
	if !ok {
		return
	}
	var body []supportResponses.SLATargetResponse
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("Invalid request body"))
		return
	}
	if err := supportServices.SaveSLATargets(actorNumber, body); err != nil {
		staffError(c, err)
		return
	}
	targets, _ := supportServices.ListSLATargets()
	c.JSON(http.StatusOK, gin.H{"message": "SLA targets saved", "data": targets})
}

func OwnerSupportReport(c *gin.Context) {
	report, err := supportServices.SupportReport()
	if err != nil {
		staffError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Support report fetched", "data": report})
}

func OwnerEscalate(c *gin.Context) {
	actorNumber, ok := actor(c)
	if !ok {
		return
	}
	var body supportRequests.EscalateRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("Invalid request body"))
		return
	}
	detail, err := supportServices.EscalateTicket(actorNumber, c.Param("number"), body.Reason)
	if err != nil {
		staffError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Ticket escalated", "data": detail})
}

func StaffProfile(c *gin.Context) {
	userNumber := supportMiddleware.SupportUser(c)
	role := supportMiddleware.SupportRole(c)
	staff, _ := supportServices.ListSupportStaff()
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"user_number": userNumber, "role": role, "staff": staff}})
}

func StaffOverview(c *gin.Context) {
	overview, err := supportServices.StaffOverview(supportMiddleware.SupportRole(c), supportMiddleware.SupportUser(c))
	if err != nil {
		staffError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": overview})
}

func StaffTickets(c *gin.Context) {
	page, perPage := supportPage(c)
	result, err := supportServices.ListStaffTickets(supportMiddleware.SupportRole(c), supportMiddleware.SupportUser(c), c.Query("search"), c.Query("status"), c.Query("category"), c.Query("priority"), c.Query("organisation"), page, perPage)
	if err != nil {
		staffError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

func StaffTicket(c *gin.Context) {
	detail, err := supportServices.GetStaffTicket(supportMiddleware.SupportRole(c), supportMiddleware.SupportUser(c), c.Param("number"))
	if err != nil {
		staffError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": detail})
}

func StaffMessage(c *gin.Context) {
	actorNumber := supportMiddleware.SupportUser(c)
	role := supportMiddleware.SupportRole(c)
	if err := ensureStaffTicketAccess(c, actorNumber, role); err != nil {
		staffError(c, err)
		return
	}
	var body supportRequests.AddMessageRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("Invalid request body"))
		return
	}
	authorRole := role
	if role == "OWNER" {
		authorRole = supportModels.AuthorOwner
	}
	message, err := supportServices.AddMessage(actorNumber, c.Param("number"), "", authorRole, body.Visibility, body.Body)
	if err != nil {
		staffError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": message})
}

func StaffUpdate(c *gin.Context) {
	actorNumber := supportMiddleware.SupportUser(c)
	role := supportMiddleware.SupportRole(c)
	if err := ensureStaffTicketAccess(c, actorNumber, role); err != nil {
		staffError(c, err)
		return
	}
	var body supportRequests.UpdateTicketRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("Invalid request body"))
		return
	}
	if role == supportModels.SupportRoleAgent && body.AssignedToNumber != nil {
		c.JSON(http.StatusForbidden, commons.Feedback("Support agents cannot assign tickets"))
		return
	}
	detail, err := supportServices.UpdateTicket(actorNumber, c.Param("number"), body)
	if err != nil {
		staffError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": detail})
}

func StaffEscalate(c *gin.Context) {
	actorNumber := supportMiddleware.SupportUser(c)
	role := supportMiddleware.SupportRole(c)
	if err := ensureStaffTicketAccess(c, actorNumber, role); err != nil {
		staffError(c, err)
		return
	}
	var body supportRequests.EscalateRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("Invalid request body"))
		return
	}
	detail, err := supportServices.EscalateTicket(actorNumber, c.Param("number"), body.Reason)
	if err != nil {
		staffError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": detail})
}

func StaffWorkload(c *gin.Context) {
	workload, err := supportServices.SupportWorkload()
	if err != nil {
		staffError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": workload})
}

func StaffAttachments(c *gin.Context) {
	actorNumber := supportMiddleware.SupportUser(c)
	role := supportMiddleware.SupportRole(c)
	if err := ensureStaffTicketAccess(c, actorNumber, role); err != nil {
		staffError(c, err)
		return
	}
	var ticket supportModels.Ticket
	if err := configTicket(c.Param("number"), &ticket); err != nil {
		staffError(c, err)
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("A file is required"))
		return
	}
	visibility := c.PostForm("visibility")
	response, err := supportServices.SaveSupportAttachment(ticket, actorNumber, visibility, file)
	if err != nil {
		staffError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": response})
}

func StaffAttachmentDownload(c *gin.Context) {
	actorNumber := supportMiddleware.SupportUser(c)
	role := supportMiddleware.SupportRole(c)
	attachment, path, err := supportServices.AttachmentPath(c.Param("number"), "", true)
	if err != nil {
		staffError(c, err)
		return
	}
	if role == supportModels.SupportRoleAgent {
		var ticket supportModels.Ticket
		if err := configTicket(attachment.TicketNumber, &ticket); err != nil || ticket.AssignedToNumber == nil || *ticket.AssignedToNumber != actorNumber {
			c.JSON(http.StatusNotFound, commons.Feedback("Attachment not found"))
			return
		}
	}
	c.FileAttachment(path, attachment.OriginalName)
}

func ClientAttachments(c *gin.Context) {
	actorNumber, ok := actor(c)
	if !ok {
		return
	}
	organisationNumber, err := commons.GetUserEntityNumber(c)
	if err != nil || organisationNumber == "" {
		c.JSON(http.StatusConflict, commons.Feedback("No active organisation"))
		return
	}
	var ticket supportModels.Ticket
	if err := configTicket(c.Param("number"), &ticket); err != nil || ticket.OrganisationNumber != organisationNumber {
		c.JSON(http.StatusNotFound, commons.Feedback("Support ticket not found"))
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("A file is required"))
		return
	}
	response, err := supportServices.SaveSupportAttachment(ticket, actorNumber, supportModels.VisibilityCustomer, file)
	if err != nil {
		staffError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": response})
}

func ClientAttachmentDownload(c *gin.Context) {
	organisationNumber, err := commons.GetUserEntityNumber(c)
	if err != nil || organisationNumber == "" {
		c.JSON(http.StatusConflict, commons.Feedback("No active organisation"))
		return
	}
	attachment, path, err := supportServices.AttachmentPath(c.Param("number"), organisationNumber, false)
	if err != nil {
		staffError(c, err)
		return
	}
	c.FileAttachment(path, attachment.OriginalName)
}

func ensureStaffTicketAccess(c *gin.Context, userNumber, role string) error {
	if role != supportModels.SupportRoleAgent {
		return nil
	}
	_, err := supportServices.GetStaffTicket(role, userNumber, c.Param("number"))
	return err
}

func configTicket(number string, ticket *supportModels.Ticket) error {
	return config.DB.Where("number = ?", number).First(ticket).Error
}
