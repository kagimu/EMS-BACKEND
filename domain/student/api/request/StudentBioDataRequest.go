package request

import (
	"fmt"

	"mime/multipart"
)

type StudentBioDataRequest struct {
	StudentNumber string                `form:"student_number" binding:"required"`
	File          *multipart.FileHeader `form:"file" binding:"required"`
}

func (req *StudentBioDataRequest) Validate() error {
	if req.StudentNumber == "" {
		return fmt.Errorf("student_number is required")
	}
	if req.File == nil {
		return fmt.Errorf("file is required")
	}
	return nil
}
