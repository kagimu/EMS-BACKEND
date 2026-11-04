package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/services"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateFeedingGroupStaff(c *gin.Context) {
	if _, err := commons.GetUserFromHeader(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.FeedingGroupStaffRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input- " + err.Error()})
		return
	}
	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, staff := range payload.StaffNumber {
		err := services.CreateFeedingGroupStaff(staff, payload.CampusNumber, payload.FeedingGroupNumber, payload.Status, tx)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save staff feeding group assignments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Staff successfully assigned to feeding group", "data": "Staff successfully assigned to feeding group"})
}

func FetchFeedingGroupStaffs(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	schoolNumber := c.Query("school_number")
	campusNumber := c.Query("campus_number")
	groupNumber := c.Query("group_number")

	response, err := services.MealStaffs(schoolNumber, campusNumber, groupNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

func UpdateFeedingGroupStaff(c *gin.Context) {
	if _, err := commons.GetUserFromHeader(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var payload struct {
		Number             string      `json:"number"`
		StaffNumber        interface{} `json:"staff_number"`
		FeedingGroupNumber string      `json:"feeding_group_number"`
		CampusNumber       string      `json:"campus_number"`
		Status             string      `json:"status"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input- " + err.Error()})
		return
	}

	var staffList []string
	if payload.StaffNumber != nil {
		switch v := payload.StaffNumber.(type) {
		case string:
			if v != "" {
				staffList = append(staffList, v)
			}
		case []interface{}:
			for _, item := range v {
				if str, ok := item.(string); ok && str != "" {
					staffList = append(staffList, str)
				}
			}
		}
	}

	if payload.Number == "" && len(staffList) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Assignment number or staff number is required"})
		return
	}

	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if payload.Number != "" {
		err := services.UpdateFeedingGroupStaffByNumber(payload.Number, payload.CampusNumber, payload.FeedingGroupNumber, payload.Status, tx)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
	} else {
		for _, staff := range staffList {
			err := services.UpdateFeedingGroupStaffByStaffNumber(staff, payload.CampusNumber, payload.FeedingGroupNumber, payload.Status, tx)
			if err != nil {
				tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update staff feeding group assignment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Staff feeding group assignment updated successfully", "data": "Staff feeding group assignment updated successfully"})
}

func DeleteFeedingGroupStaff(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	number := c.Param("number")
	if number == "" {
		number = c.Query("number")
	}
	if number == "" {
		number = c.Query("feeding_group_staff_number")
	}
	staffNumber := c.Query("staff_number")

	if number == "" && staffNumber == "" {
		var req struct {
			Number                 string `json:"number"`
			FeedingGroupStaffNum   string `json:"feeding_group_staff_number"`
			StaffNumber            string `json:"staff_number"`
		}
		_ = c.ShouldBindJSON(&req)
		number = req.Number
		if number == "" {
			number = req.FeedingGroupStaffNum
		}
		staffNumber = req.StaffNumber
	}

	if number == "" && staffNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Assignment number or staff number is required for deletion"})
		return
	}

	err = services.DeleteFeedingGroupStaff(number, staffNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to remove staff from feeding group: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Staff removed from feeding group successfully", "data": "Staff removed from feeding group successfully"})
}

