package models

import (
	"fmt"
	"html"
	"strings"
	"sync"
	"time"

	classModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	schoolModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/models"
	sourceOfSchoolInfoModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/sourceOfSchoolInfo/models"
	"gorm.io/gorm"
)

type Student struct {
	gorm.Model
	ID                 uint64                                        `gorm:"primaryKey"`
	APN                string                                        `gorm:"unique;" json:"apn"`
	Surname            string                                        `json:"surname"`
	SecondName         string                                        `json:"second_name"`
	OtherName          string                                        `json:"other_name"`
	Age                uint                                          `json:"age"`
	Gender             string                                        `json:"gender"  binding:"required" gorm:"size:8;not null"`
	School             []schoolModels.School                         `gorm:"many2many:student_school;"`
	Class              []classModels.Class                           `gorm:"many2many:student_class;"`
	Reason             string                                        `json:"reason"`
	SourceOfSchoolInfo []sourceOfSchoolInfoModels.SourceOfSchoolInfo `gorm:"many2many:student_source_of_school_info;" json:"source_of_school_info"`
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
		APNFormatOptionTwo: "%d%05d",
	}

	apn.APN = fmt.Sprintf(apn.APNFormatOPtionOne, apn.APNFormatOptionTwo, apn.Year, apn.Code, apn.Serial)

	return apn, nil
}

func (s *Student) Validate() error {
	if s.Surname == "" {
		return fmt.Errorf("surname is required")
	}

	if s.SecondName == "" {
		return fmt.Errorf("second name is required")
	}

	if s.OtherName == "" {
		return fmt.Errorf("other name is required")
	}

	if s.Age == 0 {
		return fmt.Errorf("age is required")
	}

	if s.Gender == "" {
		return fmt.Errorf("gender is required")
	}

	if s.Reason == "" {
		return fmt.Errorf("reason is required")
	}

	return nil
}

func (s *Student) Prepare() {
	s.Surname = html.EscapeString(strings.TrimSpace(strings.ToLower(s.Surname)))
	s.SecondName = html.EscapeString(strings.TrimSpace(strings.ToLower(s.SecondName)))
	s.Reason = html.EscapeString(strings.TrimSpace(strings.ToLower(s.Reason)))
	s.Gender = html.EscapeString(strings.TrimSpace(strings.ToLower(s.Gender)))
}
