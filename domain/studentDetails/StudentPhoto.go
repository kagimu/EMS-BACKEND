package studentDetails

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	studentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	"time"
)

type StudentPhoto struct {
	ID            uint `gorm:"primarykey"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	StudentNumber string               `json:"student_number" gorm:"size:36;not null"`
	Student       studentModel.Student `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	Photo         string               `json:"photo" gorm:"size:255;not null"`
	UserNumber    string               `json:"user_number" gorm:"size:36;not null"`
	User          models.User          `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
