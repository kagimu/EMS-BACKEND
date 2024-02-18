package requests

import "errors"

type ServiceDiscountRequest struct {
	Reason         string `json:"reason"`
	DiscountAmount uint64 `json:"discount_amount"`
	AuthorizedBy   string `json:"authorized_by"`
	StudentNumber  string `json:"student_number"`
	ServiceNumber  string `json:"service_number"`
}

func (d *ServiceDiscountRequest) Validate() error {
	if d.DiscountAmount == 0 {
		return errors.New("discount amount is required")
	}
	if d.Reason == "" {
		return errors.New("reason is required")
	}
	if d.StudentNumber == "" {
		return errors.New("student number is required")
	}
	if d.ServiceNumber == "" {
		return errors.New("service number is required")
	}
	return nil
}
