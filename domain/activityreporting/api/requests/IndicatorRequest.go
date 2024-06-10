package requests

import "errors"

type IndicatorRequest struct {
	Name         string `json:"name"`
	EntityNumber string `json:"entity_number"`
	Description  string `json:"description"`
	Code         int    `json:"code"`
	InputType    string `json:"input_type"`
}

func isAllowedTypes(inputType string) bool {

	var allowed = []string{"VIDEO", "IMAGE", "TEXT", "NUMBER"}

	for _, t := range allowed {
		if t == inputType {
			return true
		}
	}

	return false
}

func (indicatorRequest *IndicatorRequest) Validate() error {

	if indicatorRequest.Name == "" {
		return errors.New("provide the indicator name field")
	}

	if indicatorRequest.EntityNumber == "" {
		return errors.New("provide the indicator entity")
	}

	if indicatorRequest.InputType == "" {
		return errors.New("provide the indicator input type, I.E IMAGE, VIDEO, NUMBER or TEXT")
	}

	if !isAllowedTypes(indicatorRequest.InputType) {
		return errors.New("the provided indicator input type is not allowed")
	}

	if indicatorRequest.Code <= 0 {
		return errors.New("provide the indicator code field")
	}

	return nil
}
