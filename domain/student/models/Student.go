package models

import (
	"fmt"
	"sync"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	classModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	schoolModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/models"
	sourceOfSchoolInfoModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/sourceOfSchoolInfo/models"
	"gorm.io/gorm"
)

type Student struct {
	gorm.Model
	Number                   string                                      `json:"number" gorm:"size:36;unique;not null"`
	APN                      string                                      `gorm:"unique;" json:"apn"`
	APNAlternative           string                                      `gorm:"unique;" json:"apn_alternative"`
	Surname                  string                                      `json:"surname"`
	SecondName               string                                      `json:"second_name"`
	OtherName                string                                      `json:"other_name"`
	Age                      uint                                        `json:"age"`
	Gender                   string                                      `json:"gender"  binding:"required" gorm:"size:8;not null"`
	Reason                   string                                      `json:"reason"`
	SchoolNumber             string                                      `json:"school_number" gorm:"size:36;not null"`
	School                   schoolModels.School                         `gorm:"foreignKey:SchoolNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	ClassNumber              string                                      `json:"class_number" gorm:"size:36;not null"`
	Class                    classModels.Class                           `gorm:"foreignKey:ClassNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	SourceOfSchoolInfoNumber string                                      `json:"source_of_school_info_number" gorm:"size:36;not null"`
	SourceOfSchoolInfo       sourceOfSchoolInfoModels.SourceOfSchoolInfo `gorm:"foreignKey:SourceOfSchoolInfoNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	UserNumber               string                                      `json:"user_number" gorm:"size:36;index;not null"`
	User                     models.User                                 `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}

type APNGenerator struct {
	YearCounter      map[int]int
	YearCounterMutex sync.Mutex
}

type LearnerApplication struct {
	Year               int
	Code               int
	Serial             int
	APN                string
	APNFormatOPtionOne string
	APNFormatOptionTwo string
}

func NewAPNGenerator() *APNGenerator {
	return &APNGenerator{
		YearCounter: make(map[int]int),
	}
}

func (g *APNGenerator) GenerateAPN() (*LearnerApplication, error) {
	currentYear := time.Now().Year()

	g.YearCounterMutex.Lock()
	defer g.YearCounterMutex.Unlock()

	g.YearCounter[currentYear]++
	serial := g.YearCounter[currentYear]

	apn := &LearnerApplication{
		Year:               currentYear,
		Code:               1,
		Serial:             serial,
		APNFormatOPtionOne: "[%04d][%02d][%05d]",
	}

	apn.APN = fmt.Sprintf(apn.APNFormatOPtionOne, apn.Year, apn.Code, apn.Serial)

	return apn, nil
}
