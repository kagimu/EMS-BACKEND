package requests

import "errors"

type RemoveTemplateFromClass struct {
	CampusClassNumber string `json:"campusClassNumber"`
}

func (removeTemplateFromClass *RemoveTemplateFromClass) Validate() error {

	if removeTemplateFromClass.CampusClassNumber == "" {
		return errors.New("provide campus class number")
	}

	return nil
}
