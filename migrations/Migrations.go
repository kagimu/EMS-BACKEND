package migrations

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	academicYearModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/academicYear/models"
	reportingActivityModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/activityreporting/models"
	admissionModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/models"
	admissionLettersModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admissionLetters/models"
	applicationModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/application/models"
	assessmentScaleModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/aggregateScale/models"
	assessmentGradeModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/assessmentGrade/models"
	gradesModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/grades/models"
	assessmentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/models"
	models3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/attachments/models"
	models4 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/attendance/models"
	authModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	cateringModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/models"
	classModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	commercialModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commercial/models"
	assessmentAreaModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/models"
	conceptModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/concept/models"
	curriculumModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/importCurriculum/models"
	learningAreaModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/learningarea/models"
	strandModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/strand/models"
	subjectClusterModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/subjectcluster/models"
	subTopicModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/subtopic/models"
	topicModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/topic/models"
	currDelivery "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculumdelivery/models"
	emsAddonModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/addons/models"
	emsAssignmentModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/assignments/models"
	emsAuditModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/audit/models"
	emsAuthorityModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/authority/models"
	customfieldModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/customfields/models"
	emsLicensingModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/licensing/models"
	emsNotificationsModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/notifications/models"
	emsOrganisationModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/organisation/models"
	emsReportingModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/reporting/models"
	emsStakeholderModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/stakeholders/models"
	emsWorkflowModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/workflow/models"
	enrollmentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/enrollment/models"
	eosModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/models"
	interviewModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/interview/models"
	invoiceModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/models"
	itemModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/item/models"
	productModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/item/product/models"
	serviceModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/item/service/models"
	models2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/landingPage/models"
	lessonPlanModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/lessonPlan/models"
	parentModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/parent/models"
	paymentModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/payments/models"
	permissionModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/permissions/models"
	preInterviewModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/preAdmissionInterviewTest/models"
	previousSchoolModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/previousSchool/models"
	previousSchoolResultsModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/previousSchoolResults/models"
	progressReportModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/progressReport/models"
	promotionModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/promotion/models"
	reportSettings "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/models"
	reportTemplateModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportTemplate/models"
	schemeOfWorkmodel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/schemeOfWork/models"
	schoolModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/models"
	materialModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/learningMaterials/models"
	LevelModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/level/models"
	sourceOfSchoolInfoModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/sourceOfSchoolInfo/models"
	staffModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/staff/models"
	studentModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studentDetails"
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
		&timetableModels.StudyTimetablePeriod{},
		&timetableModels.TimetableSlots{},
		&timetableModels.TimetableTask{},
		&eosModels.GeneralCalendar{},
		&eosModels.CalendarEvent{},
		&academicYearModel.AcademicYear{},
		&schemeOfWorkmodel.SchemeOfWork{},
		&schemeOfWorkmodel.SchemeOfWorkSubTopic{},
		&lessonPlanModel.LessonPlan{},
		&lessonPlanModel.LessonPlanConcept{},
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
		&admissionLettersModels.AdmissionLetter{},
		&admissionModel.Admission{},
		&admissionModel.AdmissionCondition{},
		&previousSchoolModels.PreviousSchoolClass{},
		&materialModels.CurriculumSubtopicReference{},
		&schemeOfWorkmodel.SchemeOfWorkImport{},
		&paymentModels.Payment{},
		&preInterviewModels.PreAdmissionInterviewTestReport{},
		&interviewModel.Interview{},
		&materialModels.MaterialSubType{},
		&lessonPlanModel.LessonPlanReference{},
		&interviewModel.InterviewStudent{},
		&itemModels.Item{},
		&productModels.Product{},
		&serviceModel.Service{},
		&invoiceModel.Statement{},
		&invoiceModel.Invoice{},
		&invoiceModel.InvoiceItem{},
		&serviceModel.Discount{},
		&invoiceModel.Transaction{},
		&invoiceModel.InvoiceTransaction{},
		&assessmentScaleModels.GradingStandard{},
		&assessmentAreaModel.AssessmentArea{},
		&assessmentModel.Assessment{},
		&assessmentModel.Metric{},
		&assessmentScaleModels.GradingStandard{},
		&assessmentModel.AssessmentMetric{},
		&models4.Attendance{},
		&permissionModels.Role{},
		&permissionModels.RolePermission{},
		&permissionModels.UserRole{},
		&assessmentScaleModels.Scale{},
		&assessmentScaleModels.AggregateScale{},
		&assessmentGradeModels.AssessmentGrade{},
		&staffModel.SalaryScale{},
		&staffModel.Task{},
		&staffModel.TaskSupervision{},
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
		&promotionModels.PromotionSetting{},
		&promotionModels.PromotionAssessmentAreaValue{},
		&promotionModels.PromotionSettingValue{},
		&promotionModels.Promotion{},
		&previousSchoolResultsModels.PreviousSchoolResults{},
		&staffModel.Staff{},
		&cateringModels.FeedingGroup{},
		&cateringModels.FeedingGroupPerDay{},
		&cateringModels.FeedingSessionPerDay{},
		&cateringModels.FeedingSession{},
		&cateringModels.UnitOfMeasure{},
		&cateringModels.MenuItem{},
		&cateringModels.ItemAssembly{},
		&cateringModels.Stocking{},
		&cateringModels.Feeding{},
		&cateringModels.FeedingSessionItem{},
		&cateringModels.FeedingGroupUser{},
		&cateringModels.FeedingGroupSession{},
		&cateringModels.FeedingAllGroups{},
		&cateringModels.Meal{},
		&permissionModels.Scope{},
		&reportTemplateModel.ReportTemplate{},
		&reportTemplateModel.ClassReportTemplate{},
		&cateringModels.AssembledPart{},
		&cateringModels.MealCourseCount{},
		&cateringModels.MealCourseCount{},
		&cateringModels.MealCourse{},
		&enrollmentModel.EnrollmentItem{},
		&cateringModels.MealCourseFeedingSession{},
		&cateringModels.CampusMealCourse{},
		&cateringModels.FeedingGroupMealCourse{},
		&cateringModels.MealStock{},
		&eosModels.EntityUser{},
		&reportSettings.HeadTeacherComment{},
		&reportingActivityModels.Indicator{},
		&eosModels.ActiveUserEntity{},
		&gradesModels.GradeGroup{},
		&gradesModels.Grade{},
		&cateringModels.CalendarEventFeedingGroup{},
		&assessmentGradeModels.GradingStandardClassScale{},
		&assessmentGradeModels.AssessmentGradeScale{},
		&models2.LandingPage{},
		&studyPeriodModel.StudyPeriodName{},
		&studyPeriodModel.ClassStudyPeriod{},
		&schoolModels.ModeOfStudy{},
		&invoiceModel.Billable{},
		&invoiceModel.BillableItem{},
		&invoiceModel.BillableInvoice{},
		&invoiceModel.EnrollmentInvoice{},
		&invoiceModel.PrepaidBalance{},
		&invoiceModel.PrepaidDeposit{},
		&serviceModel.CampusItem{},
		&studentDetails.StudentPhoto{},
		&models3.Attachment{},
		&staffModel.StreamTeacher{},
		&assessmentModel.ScriptMarkingGuide{},
		&currDelivery.CurriculumDelivery{},
		&currDelivery.StudentCurriculumDelivery{},
		&currDelivery.CurriculumDeliveryStep{},
		&assessmentModel.Narrative{},
		&assessmentModel.ScheduledNarrative{},
		&assessmentModel.NarrativeSchedule{},
		&reportSettings.LearningAreaReportComment{},
		&studentModels.MedicalRecord{},
		&studentModels.FamilyMember{},
		&studentModels.FamilyMemberContact{},
		&studentModels.StudentLanguage{},
		&studentModels.EmergencyInformation{},
		&admissionModel.AdmissionLetterTemplate{},
		&reportSettings.ClassTeacherComment{},
		&reportSettings.SubjectTeacherComment{},
		&reportSettings.SubjectGroup{},
		&reportSettings.GroupAssessmentArea{},
		&assessmentModel.ScheduledNarrativeScriptMarkingGuide{},
		&admissionModel.AdmissionRequirement{},
		&studentModels.StudentBioData{},
		&admissionModel.AdmissionInterview{},
		// EMS Core Foundation A: addon/module registry (additive; ems_modules only).
		&emsAddonModels.Module{},
		// EMS Core Foundation B: organisation units (additive; ems_organisation_units only).
		&emsOrganisationModels.OrganisationUnit{},
		// EMS Core Foundation C1: generic stakeholder registry (additive; ems_stakeholders only).
		&emsStakeholderModels.Stakeholder{},
		// EMS Core Foundation C3: generic user assignments (additive; ems_user_assignments only).
		&emsAssignmentModels.UserAssignment{},
		// EMS Core Stage 4: explicit reporting lines between assignments
		// (additive; ems_reporting_lines only). No existing table is altered.
		&emsReportingModels.ReportingLine{},
		// EMS Core Foundation D: decision authority (additive; ems_decision_authorities only).
		&emsAuthorityModels.DecisionAuthority{},
		// EMS Core Foundation E1: workflow and approvals (additive; ems_approval_requests,
		// ems_approval_steps and ems_approval_assignments only).
		&emsWorkflowModels.ApprovalRequest{},
		&emsWorkflowModels.ApprovalStep{},
		&emsWorkflowModels.ApprovalAssignment{},
		// EMS Core Foundation D (unit narrowing): normalised unit-aware approval
		// anchors (ratified D6). Additive — the legacy ems_approval_steps.anchors
		// column is retained for pre-amendment rows and is never migrated.
		&emsWorkflowModels.ApprovalStepAnchor{},
		// Foundation E4 (E4-2): organisation-scoped advisor default deadline
		// configuration (ems_workflow_deadline_configs). Additive only; the
		// resolved deadline is persisted per request on ems_approval_requests.
		&emsWorkflowModels.WorkflowDeadlineConfig{},
		// EMS Core Foundation F (Increment 1): licence registry.
		&emsLicensingModels.Licence{},
		// EMS Core Foundation F (Increment 2): installation and activation.
		&emsLicensingModels.Installation{},
		// EMS Core Foundation F (Increment 5): append-only licensing facts for
		// Foundation G to reference later (§17/§18). Facts only — no engine,
		// no delivery and no derived state.
		&emsLicensingModels.LicenceEvent{},
		// EMS Core Foundation G (Increment 1): general cross-cutting audit
		// facts (ems_audit_logs). Append-only; Foundation F's ems_licence_events
		// remains the authoritative licensing fact source and is not duplicated.
		&emsAuditModels.AuditLog{},
		// EMS Core Foundation G (Increment 3): recipient-scoped in-app
		// notifications (ems_notifications). In-app only — no external delivery,
		// scheduler, worker or queue (§14/§21).
		&emsNotificationsModels.Notification{},
		// EMS Commercial Onboarding (test-payment workflow, additive): the
		// installation-package catalogue, installation setup orders and their
		// quotation snapshots, software-subscription orders, one-shot activation
		// tokens (setup/subscription; hash only) and simulated payments. No
		// existing licensing table is altered.
		&commercialModels.InstallationPackage{},
		&commercialModels.SetupOrder{},
		&commercialModels.SubscriptionOrder{},
		&commercialModels.ActivationToken{},
		&commercialModels.SimulatedPayment{},
		&customfieldModels.FieldDefinition{},
		&customfieldModels.FieldValue{},
	)

	if err != nil {
		return err
	}

	return nil
}
