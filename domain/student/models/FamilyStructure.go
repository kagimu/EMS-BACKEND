package models

type FamilyMember struct {
	Number                 string `json:"number" gorm:"size:36;unique;not null"`
	StudentNumber          string `json:"student_number" gorm:"size:36;not null"`
	Surname                string `json:"surname" gorm:"size:50;not null"`
	OtherName              string `json:"other_name" gorm:"size:50;not null"`
	PlaceOfEmployment      string `json:"place_of_employment" gorm:"size:50;not null"`
	PhysicalContactAddress string `json:"physical_contact_address" gorm:"size:100;not null"`
	Relationship           string `json:"relationship" gorm:"size:100;not null"`
}

type FamilyMemberContact struct {
	FamilyMemberNumber string `json:"family_member_number" gorm:"size:36;not null"`
	Contact            string `json:"contact" gorm:"size:50;not null"`
	Code               uint   `json:"code" gorm:"size:50;not null"`
}

type StudentLanguage struct {
	StudentNumber string `json:"student_number" gorm:"size:36;not null"`
	Code          int64  `json:"code"`
	Name          string `json:"name" gorm:"size:50;not null"`
}
