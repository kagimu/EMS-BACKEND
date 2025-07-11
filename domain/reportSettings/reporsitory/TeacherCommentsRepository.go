package reporsitory

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/models"
	"gorm.io/gorm"
)

type TeacherCommentsRepository interface {
	CreatClasTeacherComment(teacherComment models.ClassTeacherComment) error
	GetClassTeacherCommentByNumber(number uint) (models.ClassTeacherComment, error)
	DeleteClassTeacherComment(number string) error
	CreateSubjectTeacherComment(teacherComment models.SubjectTeacherComment) error
	GetSubjectTeacherCommentByNumber(number uint) (models.SubjectTeacherComment, error)
	DeleteSubjectTeacherComment(number string) error
	FindClassTeacherComments() ([]models.ClassTeacherComment, error)
	FindSubjectTeacherComments() ([]models.SubjectTeacherComment, error)
}

type classTeacherCommentsRepository struct {
	DB *gorm.DB
}

func (c classTeacherCommentsRepository) FindClassTeacherComments() ([]models.ClassTeacherComment, error) {
	var classTeacherComments = make([]models.ClassTeacherComment, 0)
	if err := c.DB.Find(&classTeacherComments).Error; err != nil {
		return classTeacherComments, err
	}
	return classTeacherComments, nil
}

func (c classTeacherCommentsRepository) FindSubjectTeacherComments() ([]models.SubjectTeacherComment, error) {
	var subjectTeacherComments = make([]models.SubjectTeacherComment, 0)
	if err := c.DB.Find(&subjectTeacherComments).Error; err != nil {
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

func (c classTeacherCommentsRepository) GetSubjectTeacherCommentByNumber(number uint) (models.SubjectTeacherComment, error) {
	var teacherComment models.SubjectTeacherComment
	if err := c.DB.Where("id = ?", number).First(&teacherComment).Error; err != nil {
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

func (c classTeacherCommentsRepository) GetClassTeacherCommentByNumber(number uint) (models.ClassTeacherComment, error) {
	var teacherComment models.ClassTeacherComment
	if err := c.DB.Where("id = ?", number).First(&teacherComment).Error; err != nil {
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
