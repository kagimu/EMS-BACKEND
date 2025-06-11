package request

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/api/responses"
)

type FamilyStructureRequest struct {
	StudentNumber string               `json:"student_number"`
	FamilyMembers []FamilyMember       `json:"family_members" `
	Languages     []responses.Language `json:"languages"`
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
	Number                 string                `json:"number"`
	StudentNumber          string                `json:"student_number"`
	Surname                string                `json:"surname"`
	OtherName              string                `json:"other_name"`
	PlaceOfEmployment      string                `json:"place_of_employment"`
	PhysicalContactAddress string                `json:"physical_contact_address"`
	Contacts               []FamilyMemberContact `json:"contacts"`
	Relationship           string                `json:"relationship"`
}

type FamilyMemberContact struct {
	PhoneNumber string `json:"phone_number" gorm:"size:50;not null"`
	Email       string `json:"email" gorm:"size:50;not null"`
}
