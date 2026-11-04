package models

import "time"

type ReportSummaryCorrespondence struct {
	Number              string    `json:"number" gorm:"primary_key"`
	ReportSummaryNumber string    `json:"report_summary_number" gorm:"column:report_summary_number"`
	Comment             string    `json:"comment" gorm:"column:comment;type:text"`
	CommenterName       string    `json:"commenterName" gorm:"column:commenter_name"`
	CreatedAt           time.Time `json:"createdAt" gorm:"column:created_at"`
}

func (b *ReportSummaryCorrespondence) TableName() string {
	return "report_summary_correspondence"
}
