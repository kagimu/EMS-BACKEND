package requests

type StaffBioDataFamilyChildRequest struct {
	Surname    string `json:"surname"`
	OtherNames string `json:"otherNames"`
	Dob        string `json:"dob"`
}

type StaffBioDataFamilyRequest struct {
	FatherSurname         string `json:"fatherSurname"`
	FatherOtherNames      string `json:"fatherOtherNames"`
	FatherEmployment      string `json:"fatherEmployment"`
	FatherPhone1          string `json:"fatherPhone1"`
	FatherPhone2          string `json:"fatherPhone2"`
	FatherEmail1          string `json:"fatherEmail1"`
	FatherEmail2          string `json:"fatherEmail2"`
	FatherPhysicalAddress string `json:"fatherPhysicalAddress"`
	FatherLiving          string `json:"fatherLiving"`

	MotherSurname         string `json:"motherSurname"`
	MotherOtherNames      string `json:"motherOtherNames"`
	MotherEmployment      string `json:"motherEmployment"`
	MotherPhone1          string `json:"motherPhone1"`
	MotherPhone2          string `json:"motherPhone2"`
	MotherEmail1          string `json:"motherEmail1"`
	MotherEmail2          string `json:"motherEmail2"`
	MotherPhysicalAddress string `json:"motherPhysicalAddress"`
	MotherLiving          string `json:"motherLiving"`

	BrotherSurname        string `json:"brotherSurname"`
	BrotherOtherNames     string `json:"brotherOtherNames"`
	SisterSurname         string `json:"sisterSurname"`
	SisterOtherNames      string `json:"sisterOtherNames"`

	SpouseSurname         string `json:"spouseSurname"`
	SpouseMaidenName      string `json:"spouseMaidenName"`
	SpouseOtherNames      string `json:"spouseOtherNames"`
	SpouseEmployment      string `json:"spouseEmployment"`
	SpousePhone1          string `json:"spousePhone1"`
	SpousePhone2          string `json:"spousePhone2"`
	SpouseEmail1          string `json:"spouseEmail1"`
	SpouseEmail2          string `json:"spouseEmail2"`
	SpousePhysicalAddress string `json:"spousePhysicalAddress"`
	SpouseLiving          string `json:"spouseLiving"`

	Children []StaffBioDataFamilyChildRequest `json:"children"`
}

func (r *StaffBioDataFamilyRequest) Validate() error {
	return nil
}
