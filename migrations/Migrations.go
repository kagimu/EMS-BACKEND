package migrations

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	academicYearModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/academicYear/models"
	admissionModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/models"
	applicationModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/application/models"
	assessmentScaleModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/aggregateScale/models"
	assessmentGradeModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/assessmentGrade/models"
	assessmentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/models"
	authModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	cateringModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/models"
	classModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	assessmentAreaModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/models"
	conceptModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/concept/models"
	curriculumModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/importCurriculum/models"
	learningAreaModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/learningarea/models"
	strandModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/strand/models"
	subjectClusterModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/subjectcluster/models"
	subTopicModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/subtopic/models"
	topicModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/topic/models"
	enrollmentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/enrollment/models"
	eosModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/models"
	gradeModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/grades/models"
	interviewModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/interview/models"
	invoiceModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/models"
	invoicePaymentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoicePayment/models"
	itemModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/item/models"
	productModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/item/product/models"
	serviceModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/item/service/models"
	lessonPlanModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/lessonPlan/models"
	parentModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/parent/models"
	paymentModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/payments/models"
	preInterviewModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/preAdmissionInterviewTest/models"
	previousSchoolModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/previousSchool/models"
	previousSchoolResultsModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/previousSchoolResults/models"
	progressReportModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/progressReport/models"
	promotionModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/promotion/models"
	reportSettings "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/models"
	schemeOfWorkmodel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/schemeOfWork/models"
	schoolModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/models"
	materialModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/learningMaterials/models"
	LevelModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/level/models"
	sourceOfSchoolInfoModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/sourceOfSchoolInfo/models"
	//staffModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/staff/models"
	studentModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	studentDetailsModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studentDetails/models"
	studyPeriodModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/models"
	subjectModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/subject/models"
	timetableModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/models"
)

func Init() error {
	err := config.DB.AutoMigrate(
		&authModels.User{},
		&eosModels.Entity{},
		&authModels.Contact{},
		&authModels.PasswordReset{},
		&authModels.UserPassword{},
		&schoolModels.School{},
		&schoolModels.SchoolCurriculum{},
		&curriculumModel.Curriculum{},
		&models.Campus{},
		&models.CampusCurriculum{},
		&LevelModels.CampusLevel{},
		&LevelModels.Level{},
		&LevelModels.CurriculumLevel{},
		&classModel.Class{},
		&learningAreaModels.LearningArea{},
		&classModel.ClassLevelSubject{},
		&subjectClusterModels.SubjectCluster{},
		&subjectModels.Subject{},
		&strandModels.Strand{},
		&topicModels.Topic{},
		&subTopicModels.SubTopic{},
		&conceptModels.Concept{},
		&classModel.CampusClass{},
		&curriculumModel.ImportCurriculum{},
		&classModel.Stream{},
		&academicYearModel.AcademicYear{},
		&studyPeriodModel.StudyPeriod{},
		&curriculumModel.ImportCurriculum{},
		&classModel.Stream{},
		&timetableModels.StudyTimetable{},
		&timetableModels.StudyTimetablePeriods{},
		&timetableModels.TimetableSlots{},
		&eosModels.GeneralCalendar{},
		&eosModels.CalendarEvent{},
		&academicYearModel.AcademicYear{},
		&schemeOfWorkmodel.SchemeOfWork{},
		&schemeOfWorkmodel.SchemeOfWorkSubTopic{},
		&lessonPlanModel.LessonPlan{},
		&lessonPlanModel.LessonPlanConcepts{},
		&lessonPlanModel.LessonPlanStep{},
		&studentDetailsModels.StudentDetail{},
		&classModel.CurriculumClass{},
		&sourceOfSchoolInfoModels.SourceOfSchoolInfo{},
		&applicationModel.Application{},
		&parentModels.Parent{},
		&studentModels.Student{},
		&previousSchoolModels.PreviousSchool{},
		&previousSchoolModels.PreviousSchoolClass{},
		&studentModels.APN{},
		&studentModels.ADN{},
		&gradeModel.Grade{},
		&admissionModel.Admission{},
		&admissionModel.AdmissionOfferCondition{},
		&admissionModel.AdmissionOffer{},
		&previousSchoolModels.PreviousSchoolClass{},
		&materialModels.CurriculumSubtopicReference{},
		&schemeOfWorkmodel.SchemeOfWorkImport{},
		&paymentModels.Payment{},
		&preInterviewModels.PreAdmissionInterviewTestReport{},
		&interviewModel.Interview{},
		&materialModels.MaterialSubType{},
		&lessonPlanModel.LessonPlanReference{},
		&interviewModel.InterviewQuestion{},
		&interviewModel.Answer{},
		&interviewModel.InterviewResult{},
		&interviewModel.InterviewStudent{},
		&itemModels.Item{},
		&productModels.Product{},
		&serviceModel.Service{},
		&invoiceModel.Invoice{},
		&invoiceModel.InvoiceItem{},
		&serviceModel.SerialCode{},
		&serviceModel.CampusServices{},
		&serviceModel.Discount{},
		&invoicePaymentModel.InvoicePayment{},
		&invoiceModel.InvoiceNumber{},
		&assessmentScaleModels.GradingStandard{},
		&assessmentAreaModel.AssessmentArea{},
		&assessmentModel.Assessment{},
		&assessmentModel.Metric{},
		&assessmentScaleModels.GradingStandard{},
		&assessmentModel.AssessmentMetric{},
		&assessmentModel.AssessmentSubMetric{},
		&assessmentModel.Attendance{},
		&authModels.Role{},
		&authModels.RolePermission{},
		&authModels.UserRole{},
		&assessmentScaleModels.AggregateScale{},
		&assessmentScaleModels.AlternativeAggregateScale{},
		&assessmentGradeModels.AssessmentGrade{},
		//&staffModel.JobTitle{},
		//&staffModel.SalaryScale{},
		//&staffModel.Task{},
		//&staffModel.TaskSupervision{},
		//&staffModel.Staff{},
		&interviewModel.InterviewAssessmentMetric{},
		&progressReportModels.ProgressReport{},
		&reportSettings.ReportSetting{},
		&reportSettings.ReportColumn{},
		&reportSettings.ReportColumnAssessment{},
		&reportSettings.ReportSummary{},
		&reportSettings.ReportScore{},
		&reportSettings.ReportScoreSummary{},
		&enrollmentModel.Enrollment{},
		&assessmentAreaModel.AssessmentAreaGradingStandard{},
		&promotionModels.Promotion{},
		&promotionModels.PromotionSetting{},
		&promotionModels.PromotionSettingsValue{},
		&previousSchoolResultsModels.PreviousSchoolResults{},
		&cateringModels.FeedingGroup{},
		&cateringModels.FeedingGroupPerDay{},
		&cateringModels.FeedingSessionPerDay{},
		&cateringModels.FeedingSession{},
		&cateringModels.UnitOfMeasure{},
		&cateringModels.MenuItem{},
		&cateringModels.ItemAssembly{},
		&cateringModels.Feeding{},
		&cateringModels.FeedingSessionItem{},
		&cateringModels.Stocking{},
	)

	if err != nil {
		return err
	}

	return nil
}
