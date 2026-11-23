package reporsitory

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/models"
	"gorm.io/gorm"
)

type TeacherCommentsRepository interface {
	CreatClasTeacherComment(teacherComment models.ClassTeacherComment) error
	GetClassTeacherCommentByNumber(number string) (models.ClassTeacherComment, error)
	DeleteClassTeacherComment(number string) error
	CreateSubjectTeacherComment(teacherComment models.SubjectTeacherComment) error
	GetSubjectTeacherCommentByNumber(number string) (models.SubjectTeacherComment, error)
	DeleteSubjectTeacherComment(number string) error
	FindClassTeacherComments() ([]models.ClassTeacherComment, error)
	FindSubjectTeacherComments(campusClassNumber, academicYearNumber, studyPeriodNumber, assessmentAreaNumber string) ([]models.SubjectTeacherComment, error)
	UpdateClassTeacherComment(comment models.ClassTeacherComment) error
	UpdateSubjectTeacherComment(comment models.SubjectTeacherComment) error
}

type classTeacherCommentsRepository struct {
	DB *gorm.DB
}

func (c classTeacherCommentsRepository) UpdateSubjectTeacherComment(comment models.SubjectTeacherComment) error {
	if err := c.DB.Save(&comment).Error; err != nil {
		return err
	}
	return nil
}

func (c classTeacherCommentsRepository) UpdateClassTeacherComment(comment models.ClassTeacherComment) error {
	if err := c.DB.Save(&comment).Error; err != nil {
		return err
	}
	return nil
}

func (c classTeacherCommentsRepository) FindClassTeacherComments() ([]models.ClassTeacherComment, error) {
	var classTeacherComments = make([]models.ClassTeacherComment, 0)
	if err := c.DB.Find(&classTeacherComments).Error; err != nil {
		return classTeacherComments, err
	}
	return classTeacherComments, nil
}

func (c classTeacherCommentsRepository) FindSubjectTeacherComments(campusClassNumber, academicYearNumber, studyPeriodNumber, assessmentAreaNumber string) ([]models.SubjectTeacherComment, error) {
	var subjectTeacherComments = make([]models.SubjectTeacherComment, 0)
	if err := c.DB.
		Model(&models.SubjectTeacherComment{}).
		Joins("JOIN assessment_areas ON assessment_areas.number = subject_teacher_comments.assessment_area_number").
		Joins("JOIN curriculum_classes ON curriculum_classes.number = assessment_areas.curriculum_class_number").
		Joins("JOIN campus_classes ON campus_classes.curriculum_class_number = curriculum_classes.number").
		Joins("JOIN study_periods ON study_periods.number = subject_teacher_comments.study_period_number").
		Scopes(
			filters.ByAnyField("campus_classes.number", campusClassNumber),
			filters.ByAnyField("study_periods.academic_year_number", academicYearNumber),
			filters.ByAnyField("subject_teacher_comments.study_period_number", studyPeriodNumber),
			filters.ByAnyField("subject_teacher_comments.assessment_area_number", assessmentAreaNumber),
		).
		Distinct("subject_teacher_comments.*").
		Find(&subjectTeacherComments).Error; err != nil {
		return subjectTeacherComments, err
	}
	return subjectTeacherComments, nil
}

func (c classTeacherCommentsRepository) CreatClasTeacherComment(teacherComment models.ClassTeacherComment) error {
	if err := c.DB.Create(&teacherComment).Error; err != nil {
		return err
	}
	return nil
}

func (c classTeacherCommentsRepository) CreateSubjectTeacherComment(teacherComment models.SubjectTeacherComment) error {
	if err := c.DB.Create(&teacherComment).Error; err != nil {
		return err
	}
	return nil
}

func (c classTeacherCommentsRepository) GetSubjectTeacherCommentByNumber(number string) (models.SubjectTeacherComment, error) {
	var teacherComment models.SubjectTeacherComment
	if err := c.DB.Where("number = ?", number).First(&teacherComment).Error; err != nil {
		return teacherComment, err
	}
	return teacherComment, nil
}

func (c classTeacherCommentsRepository) DeleteSubjectTeacherComment(number string) error {
	if err := c.DB.Where("number = ?", number).Delete(&models.SubjectTeacherComment{}).Error; err != nil {
		return err
	}
	return nil
}

func (c classTeacherCommentsRepository) GetClassTeacherCommentByNumber(number string) (models.ClassTeacherComment, error) {
	var teacherComment models.ClassTeacherComment
	if err := c.DB.Where("number = ?", number).First(&teacherComment).Error; err != nil {
		return teacherComment, err
	}
	return teacherComment, nil
}

func (c classTeacherCommentsRepository) DeleteClassTeacherComment(number string) error {
	if err := c.DB.Where("id = ?", number).Delete(&models.ClassTeacherComment{}).Error; err != nil {
		return err
	}
	return nil
}

func NewClassTeacherCommentsRepository(db *gorm.DB) TeacherCommentsRepository {
	return &classTeacherCommentsRepository{
		DB: db,
	}
}
