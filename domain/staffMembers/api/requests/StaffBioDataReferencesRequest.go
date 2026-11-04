package requests

type GuarantorRequest struct {
	Name         string `json:"name"`
	Address      string `json:"address"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Relationship string `json:"relationship"`
}

type NextOfKinRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
}

type RefereeRequest struct {
	Name         string `json:"name"`
	Address      string `json:"address"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Relationship string `json:"relationship"`
	YearKnown    string `json:"yearKnown"`
}

type StaffBioDataReferencesRequest struct {
	Guarantor1 GuarantorRequest `json:"guarantor1"`
	Guarantor2 GuarantorRequest `json:"guarantor2"`
	NextOfKin  NextOfKinRequest `json:"nextOfKin"`
	Referee1   RefereeRequest   `json:"referee1"`
	Referee2   RefereeRequest   `json:"referee2"`
}

func (r *StaffBioDataReferencesRequest) Validate() error {
	return nil
}
