package models

import "time"

type SubjectTeacherReportComment struct {
	Number                     string    `json:"number" gorm:"primaryKey;column:number"`
	SubjectTeacherReportNumber string    `json:"subject_teacher_report_number" gorm:"column:subject_teacher_report_number"`
	CommenterName              string    `json:"commenter_name" gorm:"column:commenter_name"`
	Comment                    string    `json:"comment" gorm:"column:comment"`
	CreatedAt                  time.Time `json:"created_at" gorm:"column:created_at"`
}
