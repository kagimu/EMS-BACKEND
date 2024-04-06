package filters

import "gorm.io/gorm"

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
