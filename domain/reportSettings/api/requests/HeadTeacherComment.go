package requests

import "errors"

type HeadTeacherCommentRequest struct {
	Comment      string  `json:"comment"`
	MinAggregate float64 `json:"min_aggregate"`
	MaxAggregate float64 `json:"max_aggregate"`
}

func (r *HeadTeacherCommentRequest) Validate() error {

	if r.Comment == "" {
		return errors.New("comment is required")
	}

	if r.MinAggregate == 0 {
		return errors.New("min aggregate is required")
	}

	if r.MaxAggregate == 0 {
		return errors.New("max aggregate is required")
	}

	return nil
}
