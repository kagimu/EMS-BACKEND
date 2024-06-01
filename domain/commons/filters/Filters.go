package filters

import (
	"gorm.io/gorm"
	"strings"
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
