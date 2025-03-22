package filters

import (
	userModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strings"
	"time"
)

func ByCurriculumLevelNumber(level string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if level != "" {
			return db.Where("curriculum_level_number = ?", level)
		}
		return db
	}
}

func ByPrimaryKey(level string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if level != "" {
			return db.Where("number = ?", level)
		}
		return db
	}
}

func ByLevelNumber(level string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if level != "" {
			return db.Where("level_number = ?", level)
		}
		return db
	}
}

func ByAssessmentAreaNumber(assessmentArea string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if assessmentArea != "" {
			return db.Where("assessment_area_number = ?", assessmentArea)
		}
		return db
	}
}

func ByStudyPeriodNumber(studyPeriod string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if studyPeriod != "" {
			return db.Where("study_period_number = ?", studyPeriod)
		}
		return db
	}
}

func ByCampusClassNumber(campusClass string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if campusClass != "" {
			return db.Where("campus_class_number = ?", campusClass)
		}
		return db
	}
}

func BySchoolNumber(school string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if school != "" {
			return db.Where("school_number = ?", school)
		}
		return db
	}
}

func ByStreamNumber(stream string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if stream != "" {
			return db.Where("stream_number = ?", stream)
		}
		return db
	}
}

func ByAcademicYearNumber(academicYear string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if academicYear != "" {
			return db.Where("academic_year_number = ?", academicYear)
		}
		return db
	}
}

func ByCampusNumber(campus string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if campus != "" {
			return db.Where("campus_number = ?", campus)
		}
		return db
	}
}

func ByInterviewNumber(interview string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if interview != "" {
			return db.Where("interview_number = ?", interview)
		}
		return db
	}
}

func ByStudentNumber(student string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if student != "" {
			return db.Where("student_number = ?", student)
		}
		return db
	}
}

func ByUserNumber(user string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if user != "" {
			return db.Where("user_number = ?", user)
		}
		return db
	}
}

func ByAnyField(fieldName string, fieldValue string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if fieldName != "" && fieldValue != "" {
			return db.Where(fieldName+" = ?", fieldValue)
		}
		return db
	}
}

func ByParent(fieldName string, user *userModel.User) func(db *gorm.DB) *gorm.DB {

	if user == nil {
		return func(db *gorm.DB) *gorm.DB {
			return db
		}
	}

	if user.Type == constants.TypeAdmin || user.Type == constants.TypeStaff {
		return func(db *gorm.DB) *gorm.DB {
			return db
		}
	}

	return func(db *gorm.DB) *gorm.DB {
		if fieldName != "" && user.Number != "" {
			return db.Where(fieldName+" = ?", user.Number)
		}
		return db
	}
}

func ByStatus(status string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if status != "" {
			return db.Where("status = ?", status)
		}
		return db
	}
}

func ByGradingStandardNumber(gradingStandardNumber string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if gradingStandardNumber != "" {
			return db.Where("grading_standard_number = ?", gradingStandardNumber)
		}
		return db
	}
}

func ByName(search string) func(*gorm.DB) *gorm.DB {

	return func(db *gorm.DB) *gorm.DB {

		if search != "" {
			return db.Where("name like ?", "%"+search+"%")
		}

		return db
	}
}

func ByReportTemplateNumber(templateNumber string) func(*gorm.DB) *gorm.DB {

	return func(db *gorm.DB) *gorm.DB {

		if templateNumber != "" {
			return db.Where("report_template_number = ?", templateNumber)
		}

		return db
	}
}

func ByNumber(templateNumber string) func(*gorm.DB) *gorm.DB {

	return func(db *gorm.DB) *gorm.DB {

		if templateNumber != "" {
			return db.Where("number = ?", templateNumber)
		}

		return db
	}
}

func ByGrade(value string) func(*gorm.DB) *gorm.DB {

	return func(db *gorm.DB) *gorm.DB {

		if value != "" {
			return db.Where("grade = ?", value)
		}

		return db
	}
}

func ByRefNumber(refNumber string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if refNumber != "" {
			return db.Where("ref_number = ?", refNumber)
		}
		return db
	}
}

func ByMetricNumber(metricNumber string) func(*gorm.DB) *gorm.DB {

	return func(db *gorm.DB) *gorm.DB {

		if metricNumber != "" {
			return db.Where("metric_number = ?", metricNumber)
		}

		return db
	}
}

func ByAssessmentMetricNumber(assessmentMetricNumber string) func(*gorm.DB) *gorm.DB {

	return func(db *gorm.DB) *gorm.DB {

		if assessmentMetricNumber != "" {
			return db.Where("assessment_metric_number = ?", assessmentMetricNumber)
		}

		return db
	}

}

func ByType(itemType string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if itemType != "" {
			return db.Where("type = ?", itemType)
		}
		return db
	}
}

func ByValueIn(fieldName string, fieldValues []string) func(db *gorm.DB) *gorm.DB {

	return func(db *gorm.DB) *gorm.DB {

		if fieldValues != nil {
			for i, fieldValue := range fieldValues {
				fieldValues[i] = strings.TrimSpace(fieldValue)
			}
		}

		if fieldName != "" {
			fieldName = strings.TrimSpace(fieldName)
		}

		if fieldName != "" && fieldValues != nil && len(fieldValues) > 0 {
			return db.Where(fieldName+" IN ?", fieldValues)
		}

		return db
	}
}

// ByDateBetween Will determine if all dates are set, between will be used else, greater or less than will be used depending on which date is set
func ByDateBetween(dateField string, startDate time.Time, endDate time.Time) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {

		if dateField == "" {
			return db
		}

		if !(startDate.IsZero() && endDate.IsZero()) {

			if startDate.After(endDate) {
				_startDate := startDate
				_endDate := endDate

				startDate = _endDate
				endDate = _startDate
			}

			return db.Where("DATE( "+dateField+") BETWEEN ? AND ?", startDate, endDate)
		}

		if !endDate.IsZero() {
			return db.Where("DATE( "+dateField+") <= ?", endDate)
		}

		if !startDate.IsZero() {
			return db.Where("DATE( "+dateField+") >= ?", startDate)
		}

		return db
	}

}

func ByDateAfter(dateField string, afterDate time.Time) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if dateField != "" && !afterDate.IsZero() {
			return db.Where("DATE( "+dateField+") >= ?", afterDate)
		}
		return db
	}
}

func ByDateBefore(dateField string, beforeDate time.Time) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if dateField != "" && !beforeDate.IsZero() {
			return db.Where("DATE( "+dateField+") < ?", beforeDate)
		}
		return db
	}
}

func ByEntityNumber(entityNumber string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if entityNumber != "" {
			return db.Where("entity_number = ?", entityNumber)
		}
		return db
	}
}

func ByScopedItems(category string, ctx *gin.Context, fieldName ...string) func(db *gorm.DB) *gorm.DB {
	searchField := "number"
	if len(fieldName) > 0 {
		searchField = fieldName[0]
	}
	return func(db *gorm.DB) *gorm.DB {
		items := GetUserScopedItems(ctx, category)
		if len(items) > 0 {
			return db.Scopes(ByValueIn(searchField, items))
		}
		return db
	}
}
