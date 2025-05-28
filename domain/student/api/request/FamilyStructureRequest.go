package request

import (
	"errors"
)

type FamilyStructureRequest struct {
	StudentNumber string         `json:"student_number" gorm:"size:36;not null"`
	FamilyMembers []FamilyMember `json:"family_members" `
	Languages     []struct {
		Name string
		Code int64
	} `json:"languages"`
}

func (req *FamilyStructureRequest) Validate() error {
	if req.StudentNumber == "" {
		return errors.New("student number is required")
	}

	if len(req.FamilyMembers) == 0 {
		return errors.New("family members(father/mother/guardian) is required")
	}
	return nil
}

type FamilyMember struct {
	Number                 string   `json:"number"`
	StudentNumber          string   `json:"student_number"`
	Surname                string   `json:"surname"`
	OtherName              string   `json:"other_name"`
	PlaceOfEmployment      string   `json:"place_of_employment"`
	PhysicalContactAddress string   `json:"physical_contact_address"`
	Contacts               []string `json:"contacts"`
	Relationship           string   `json:"relationship"`
}

type FamilyMemberContacts struct {
	FamilyMemberNumber string `json:"family_member_number" gorm:"size:36;not null"`
	Contact            string `json:"contact" gorm:"size:50;not null"`
}
