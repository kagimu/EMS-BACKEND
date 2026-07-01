package models

import (
	campusModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	staffMembersModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/staffMembers/models"
	"time"
)

type FeedingGroupStaff struct {
	ID                 uint `gorm:"primarykey"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
	Number             string                             `json:"number" gorm:"size:36;unique;not null"`
	FeedingGroupNumber string                             `json:"feeding_group_number" gorm:"size:36;not null"`
	Status             string                             `json:"status"`
	FeedingGroup       FeedingGroup                       `gorm:"foreignKey:FeedingGroupNumber;references:Number;"`
	StaffNumber        string                             `json:"staff_number" gorm:"size:36;index;not null"`
	StaffDeployment    staffMembersModels.StaffDeployment `gorm:"foreignKey:StaffNumber;references:Number;"`
	CampusNumber       string                             `json:"campus_number" gorm:"size:36;index"`
	Campus             campusModel.Campus                 `gorm:"foreignKey:CampusNumber;references:Number;"`
}
