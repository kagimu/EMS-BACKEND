package requests

type StaffBioDataMedicalRequest struct {
	Height              string `json:"height"`
	Weight              string `json:"weight"`
	Complexion          string `json:"complexion"`
	BloodType           string `json:"bloodType"`
	Allergies           string `json:"allergies"`
	MedicalConditions   string `json:"medicalConditions"`
	RegularMedications  string `json:"regularMedications"`
	PhysicalImpairments string `json:"physicalImpairments"`
	DoctorName          string `json:"doctorName"`
	DoctorPhone         string `json:"doctorPhone"`
	DoctorFacility      string `json:"doctorFacility"`
}

func (r *StaffBioDataMedicalRequest) Validate() error {
	return nil
}
