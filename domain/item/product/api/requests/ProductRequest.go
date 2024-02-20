package requests

import "errors"

type ProductRequest struct {
	Name               string `json:"name"`
	Description        string `json:"description"`
	SerialCode         string `json:"serial_code"`
	DiscountPercentage uint   `json:"discount_percentage"`
	DiscountAmount     uint   `json:"discount_amount"`
	Cost               int64  `json:"cost"`
	Prefix             string `json:"prefix"`
	RepeatCycle        string `json:"repeat_cycle"`
	OneTimePurchase    bool   `json:"one_time_purchase"`
}

func (p *ProductRequest) Validate() error {
	if p.Name == " " {
		return errors.New("name is required")
	}
	if p.Description == " " {
		return errors.New("description is required")
	}
	if p.SerialCode == " " {
		return errors.New("serial code is required")
	}
	if p.Cost == 0 {
		return errors.New("cost is required")
	}
	if p.Prefix == " " {
		return errors.New("prefix is required")
	}

	return nil
}

type ProductUpdateRequest struct {
	Name               string `json:"name"`
	Description        string `json:"description"`
	DiscountPercentage uint   `json:"discount_percentage"`
	DiscountAmount     uint   `json:"discount_amount"`
	Cost               uint   `json:"cost"`
	Prefix             string `json:"prefix"`
}
