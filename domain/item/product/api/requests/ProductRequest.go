package requests

import "errors"

type ProductRequest struct {
	Description     string   `json:"description"`
	Cost            float64  `json:"cost"`
	Name            string   `json:"name"`
	SerialCode      string   `json:"serial_code"`
	Mandatory       bool     `json:"mandatory"`
	CampusNumbers   []string `json:"campus_numbers"`
	OneTimePurchase bool     `json:"one_time_purchase"`
	RepeatCycle     string   `json:"repeat_cycle"`
	PaymentMode     string   `json:"payment_mode"`
	SchoolNumber    string   `json:"school_number"`
}

func (p *ProductRequest) Validate() error {
	if len(p.Name) == 0 {
		return errors.New("name is required")
	}
	if len(p.Description) == 0 {
		return errors.New("description is required")
	}
	if len(p.SerialCode) == 0 {
		return errors.New("serial code is required")
	}
	if p.Cost == 0 {
		return errors.New("cost is required")
	}
	return nil
}
