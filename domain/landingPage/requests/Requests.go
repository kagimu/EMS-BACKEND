package requests

import "errors"

type LandingPageRequest struct {
	Title string `json:"title"`
}

func (l *LandingPageRequest) Validate() error {
	if l.Title == "" {
		return errors.New("title is required")
	}

	return nil
}
