package routes

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/attachments/controller"
	controllers4 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/attendance/controllers"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	controllers3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculumdelivery/controllers"
	controllers2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/landingPage/controllers"
	"time"

	academicYearController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/academicYear/controllers"
	activityReportingControllers "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/activityreporting/api/controllers"
	admissionController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/controllers"
	applicationController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/application/controllers"
	aggregateScaleController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/aggregateScale/controllers"
	assessmentGradeController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/assessmentGrade/controllers"
	assessmentController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/controllers"
	gradesController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/grades/controllers"
	authController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/controllers"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/middleware"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/controllers"
	cateringController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/controllers"
	classController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/controllers"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	assessmentAreaController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/controllers"
	conceptController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/concept/controllers"
	curriculumController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/importCurriculum/controllers"
	learningAreaController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/learningarea/controllers"
	strandController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/strand/controllers"
	subjectClusterController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/subjectcluster/controllers"
	subTopicController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/subtopic/controllers"
	topicController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/topic/controllers"
	enrollmentController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/enrollment/controllers"
	eosController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/controllers"
	interviewController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/interview/controllers"
	invoiceController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/controllers"
	productController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/item/product/controllers"
	serviceController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/item/service/controllers"
	lessonPlanControllers "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/lessonPlan/controllers"
	parentController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/parent/controllers"
	paymentsController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/payments/controllers"
	permissionController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/permissions/controllers"
	preInterviewController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/preAdmissionInterviewTest/controllers"
	previousSchoolController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/previousSchool/controllers"
	previousSchoolResultsController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/previousSchoolResults/controllers"
	progressReportController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/progressReport/controllers"
	promotionController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/promotion/controllers"
	reportSettingsController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/controllers"
	reportTemplate "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportTemplate/api/controllers"
	schemeOfWorkControllers "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/schemeOfWork/controllers"
	schoolController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/controllers"
	learningMaterialsControllers "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/learningMaterials/controllers"
	learningLevelControllers "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/level/controllers"
	sourceOfSchoolController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/sourceOfSchoolInfo/controllers"
	staffController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/staff/controllers"
	statisticsController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/statistics/api/controllers"
	studentController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/controllers"
	studentDetailsController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studentDetails/controllers"
	studyPeriodController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/controllers"
	subjectController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/subject/controllers"
	timetableController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/controllers"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	router := gin.Default()

	//setting cors in the system
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:2054", "https://wholesome-api.noregrets.ug", "https://test-api.wholesome.noregrets.ug"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "OPTIONS", "DELETE"},
		AllowHeaders:     []string{"Access-Control-Allow-Origin", "Origin", "Content-Length", "Content-Type", "User-Agent", "Referrer", "Host", "Token", "Authorization", "Accept-Encoding", "X-CSRF-Token", "Access-Control-Allow-Headers", "Content-Type", "user-number"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			///return origin == "https://github.com"
			return true
		},
		MaxAge: 12 * time.Hour,
	}))

	//router.Use(middlewares.JWTMiddleware())
	router.MaxMultipartMemory = 20 << 20 // 20 MiB
	//router.Static("/", "./fileRepo")
	router.Use(static.Serve("/fileRepo", static.LocalFile("/fileRepo", true)))

	router.NoRoute(commons.RouteNotFound)
	router.GET("ping", commons.APIHealth)

	landingPage := router.Group("/landing-page")
	{
		landingPage.GET("", controllers2.GetLandingPage)
	}

	auth := router.Group("/auth")
	{
		auth.POST("login", authController.Login)
		auth.POST("register", authController.Register)
		auth.POST("set-password", authController.UserPassword)
		auth.POST("reset-password", authController.ResetPassword)
	}

	main := router.Group("/main", middleware.AuthorizeJWT())
	{
		main.GET("profile", authController.Profile)

		curriculum := main.Group("/curriculum")
		{
			curriculum.GET("all", curriculumController.FetchCurriculum, middleware.ValidatePermission("view_curriculum"))
			curriculum.POST("new", curriculumController.NewCurriculum)
			curriculum.GET("/info", curriculumController.CurriculumInfo)
			curriculum.GET("/campus-curriculum-level", curriculumController.FetchCampusCurriculumLevel)
			curriculum.POST("/edit", curriculumController.CurriculumEdit)
			curriculum.POST("/delete", curriculumController.DeleteCurriculum)
			curriculum.POST("upload", curriculumController.ImportCurriculum)
			curriculum.GET("uploaded", curriculumController.FetchUploadedCurriculum)
			curriculum.GET("school-curriculum", curriculumController.FetchSchoolCurriculum)
			curriculum.GET("data", curriculumController.FetchCurriculumInfo)
		}

		school := main.Group("/school")
		{
			school.GET("/all", schoolController.FetchSchools)
			school.POST("/new", schoolController.NewSchool)
			school.POST("/edit", schoolController.EditSchool)
			school.GET("/info", schoolController.FetchSchoolInfo)
			school.DELETE("/delete/:school_number", schoolController.DeleteSchool)
			school.POST("/upload-logo", schoolController.UploadSchoolLogo)

			template := school.Group("/report-templates")
			{
				template.GET("/all", reportTemplate.AllReportTemplates)
				template.POST("/new", reportTemplate.NewReportTemplate)
				template.GET("/by-class", reportTemplate.ClassReportTemplates)
				template.POST("/add-to-classes", reportTemplate.AddTemplateToClass)
				template.POST("/remove-template", reportTemplate.RemoveTemplateFromClass)
				template.DELETE("/delete/:id", reportTemplate.DeleteTemplate)
			}

			modeOfStudy := school.Group("mode-of-study")
			{
				modeOfStudy.POST("", schoolController.NewModeOfStudy)
				modeOfStudy.GET("", schoolController.AllModeOfStudy)
				modeOfStudy.PUT("", schoolController.UpdateModeOfStudy)
				modeOfStudy.DELETE("/:mode_of_study_number", schoolController.DeleteModeOfStudy)
			}

		}

		campus := main.Group("/campus")
		{
			campus.GET("/all", controllers.FetchCampuses)
			campus.GET("/class-list", classController.FetchClassLevel)
			campus.GET("/info", controllers.CampusInfo)
			campus.GET("/level-list", controllers.CampusLevelList)
			campus.GET("/curriculum", controllers.CampusCurriculum)
			campus.POST("/new", controllers.NewCampus)
			campus.POST("/edit", controllers.EditCampus)
			campus.POST("/delete", controllers.DeleteCampus)
			campus.POST("/new-level", controllers.NewCampusLevel)
			campus.POST("/new-class", classController.AddCampusClass)
		}

		level := main.Group("/level")
		{
			level.GET("/info", learningLevelControllers.LevelInfo)
			level.POST("/edit", learningLevelControllers.EditLevel)
			level.POST("/delete", learningLevelControllers.DeleteLevel)
		}

		main.GET("level", controllers.FetchLevel)
		main.GET("campus-level", controllers.FetchCampusLevel)
		main.GET("fetch-campus-level", controllers.FetchCampusLevelList)

		main.POST("new-level", learningLevelControllers.NewLevel)
		main.POST("new-curriculum-level", learningLevelControllers.NewCurriculumLevel)
		main.GET("fetch-level", learningLevelControllers.FetchLearningLevel)
		main.GET("fetch-curriculum-level", learningLevelControllers.FetchCurriculumLevel)

		class := main.Group("/class")
		{
			class.GET("/info", classController.ClassInfo)
			class.GET("/curriculum-class", classController.FetchClassesToAdd)
			class.GET("/campus-class-by-level", classController.FetchCampusClassesByLevel)
			class.GET("/by-campus-number/:campus_number", classController.FetchClassesByCampusNumber)
			class.POST("/edit", classController.EditClass)
			class.POST("/delete", classController.DeleteClass)
			class.GET("/school", classController.FetchClassesBySchoolNumber)
		}
		//main.GET("classes-to-add", classController.FetchClassesToAdd)
		main.POST("new-class", classController.NewClass)
		main.POST("new-curriculum-class", classController.AddCurriculumClass)
		main.GET("classes", classController.FetchClasses)
		main.GET("fetch-curriculum-classes", classController.FetchCurriculumClasses)
		main.GET("filter-curriculum-classes", classController.FilterCurriculumClasses)
		main.GET("scheme-of-work-classes", classController.FetchSchemeOfWorkClasses)
		main.GET("classes-curriculum", classController.FetchClassesCurriculum)
		main.GET("class-subjects", classController.FetchClassSubjects)
		main.POST("subject-to-class", classController.NewClassSubjects)

		//manage subject
		subject := main.Group("/subject")
		{
			subject.GET("/info", subjectController.SubjectInfo)
			subject.GET("/by-curriculum-class", subjectController.SubjectCurriculumClass)
			subject.POST("/edit", subjectController.SubjectEdit)
			subject.POST("/delete", subjectController.DeleteSubject)
		}

		main.POST("new-subject", subjectController.NewSubject)
		main.GET("subjects", subjectController.FetchSubject)
		main.GET("subjects/all", subjectController.FetchAllSubjects)
		main.GET("subject-list", subjectController.FetchSubjectList)
		main.GET("subjects-by-class", subjectController.FetchSubjectByClass)

		//manage Learning Area
		learningArea := main.Group("/learning-area")
		{
			learningArea.GET("/info", learningAreaController.LearningAreaInfo)
			learningArea.POST("/edit", learningAreaController.LearningAreaEdit)
			learningArea.POST("/delete", learningAreaController.DeleteLearningArea)
		}
		main.POST("new-learning-area", learningAreaController.NewLearningArea)
		main.GET("learning-areas", learningAreaController.FetchLearningArea)

		//subject cluster
		subjectCluster := main.Group("/subject-cluster")
		{
			subjectCluster.GET("/info", subjectClusterController.SubjectClusterInfo)
			subjectCluster.POST("/edit", subjectClusterController.SubjectClusterEdit)
			subjectCluster.POST("/delete", subjectClusterController.DeleteSubjectCluster)
		}

		main.POST("new-subject-cluster", subjectClusterController.NewSubjectCluster)
		main.GET("subject-clusters", subjectClusterController.FetchSubjectCluster)
		main.GET("subject-cluster-filter", subjectClusterController.FetchSubjectClusterList)

		//strand
		strand := main.Group("/strand")
		{
			strand.GET("/info", strandController.StrandInfo)
			strand.POST("/edit", strandController.StrandEdit)
			strand.POST("/delete", strandController.DeleteStrand)
		}
		main.POST("new-strand", strandController.NewStrand)
		main.GET("strands", strandController.FetchStrand)
		main.GET("strand-list", strandController.FetchStrandList)

		//topic
		topic := main.Group("/topic")
		{
			topic.GET("/info", topicController.TopicInfo)
			topic.POST("/edit", topicController.TopicEdit)
			topic.POST("/delete", topicController.DeleteTopic)
		}
		main.POST("new-topic", topicController.NewTopic)
		main.GET("topics", topicController.FetchTopic)
		main.GET("topic-list", topicController.FetchTopicList)

		//sub topic
		subTopic := main.Group("/sub-topic")
		{
			subTopic.GET("/info", subTopicController.SubTopicInfo)
			subTopic.POST("/edit", subTopicController.SubTopicEdit)
			subTopic.POST("/delete", subTopicController.DeleteSubTopic)
		}
		main.POST("new-sub-topic", subTopicController.NewSubTopic)
		main.GET("sub-topics", subTopicController.FetchSubTopic)
		main.GET("sub-topic-list", subTopicController.FetchSubTopiList)
		main.GET("sub-topics-by-class", subTopicController.FetchSubTopicByClass)

		//concepts
		concept := main.Group("/concept")
		{
			concept.GET("/info", conceptController.ConceptInfo)
			concept.POST("/edit", conceptController.ConceptEdit)
			concept.POST("/delete", conceptController.DeleteConcept)
		}
		main.POST("new-concept", conceptController.NewConcept)
		main.GET("concepts", conceptController.FetchConcept)

		//adds academic year
		academicYear := main.Group("/academic-year")
		{
			academicYear.GET("all", academicYearController.FetchAcademicYearInfo)
			academicYear.GET("info", academicYearController.FetchAcademicYearData)
			academicYear.GET("by-level", academicYearController.AcademicYearsByLevel)
			academicYear.GET("by-campus-class", academicYearController.AcademicYearsByCampusClass)
			academicYear.POST("/new", academicYearController.NewAcademicYear)
			academicYear.POST("/edit", academicYearController.EditAcademicYear)
			academicYear.POST("/delete", academicYearController.DeleteAcademicYear)
		}
		main.GET("academic-by-class", academicYearController.FetchAcademicByClass)

		//adds study period
		studyPeriod := main.Group("/study-period")
		{
			studyPeriod.GET("all", studyPeriodController.AllStudyPeriods)
			studyPeriod.POST("new", studyPeriodController.NewStudyPeriod)
			studyPeriod.POST("edit", studyPeriodController.EditStudyPeriod)
			studyPeriod.DELETE("delete/:study_period_number", studyPeriodController.RemoveStudyPeriod)
			studyPeriod.GET("info", studyPeriodController.FetchStudyPeriodDetails)
			////studyPeriod.GET("all", studyPeriodController.FetchStudyPeriod)
			//studyPeriod.GET("by-academic-year", studyPeriodController.FetchStudyPeriodByAcademicYear)

			studyPeriodNames := studyPeriod.Group("names")
			{
				studyPeriodNames.POST("", studyPeriodController.NewStudyPeriodName)
				studyPeriodNames.GET("", studyPeriodController.AllStudyPeriodName)
				studyPeriodNames.PUT("", studyPeriodController.UpdateStudyPeriodName)
				studyPeriodNames.DELETE("/:study_period_name_number", studyPeriodController.RemoveStudyPeriodName)
			}

		}

		//adds scheme of work
		schemeOfWork := main.Group("/scheme-of-work")
		{
			schemeOfWork.GET("download-scheme-template", schemeOfWorkControllers.DownloadSchemeOfWork)
			schemeOfWork.POST("upload-scheme-of-work", schemeOfWorkControllers.ImportSchemeOfWork)
			schemeOfWork.PATCH("invalidate/:scheme_of_work_number", schemeOfWorkControllers.InvalidateSchemeOfWork)
			schemeOfWork.PATCH("activate/:scheme_of_work_number", schemeOfWorkControllers.ActivateSchemeOfWork)
			schemeOfWork.GET("fetch-scheme-of-work-by-class-and-subject", schemeOfWorkControllers.FetchSchemeOfWork)
		}
		main.POST("new-scheme-of-work", schemeOfWorkControllers.NewSchemeOfWork)
		main.POST("update-scheme-of-work", schemeOfWorkControllers.UpdateSchemeOfWork)

		//lesson notes
		main.GET("fetch-concepts-by-subtopic", conceptController.FetchConceptBySubTopic)

		lessonPlan := main.Group("/lesson-plan")
		{
			lessonPlan.POST("/new-lesson-Plan", lessonPlanControllers.NewLessonPlan)
			lessonPlan.POST("/add-reference", lessonPlanControllers.AddLessonPlanReference)
			lessonPlan.PATCH("/update-reference/:referenceNumber", lessonPlanControllers.UpdateLessonPlanReference)
			lessonPlan.DELETE("/delete-reference/:referenceNumber", lessonPlanControllers.DeleteLessonPlanReference)
			lessonPlan.GET("/fetch-lesson-plan", lessonPlanControllers.FetchLessonPlan)
			lessonPlan.GET("/references/:id", lessonPlanControllers.FetchLessonPlanReferences)
			lessonPlan.POST("import", lessonPlanControllers.ImportLessonPlan)
			lessonPlan.GET("export-template/:scheme_sub_topic_number", lessonPlanControllers.ExportLessonPlanTemplate)
		}

		studentDetails := main.Group("/student-details")
		{
			studentDetails.POST("/:student_number", studentDetailsController.CreateStudentDetails)
			studentDetails.PUT("/:student_number", studentDetailsController.EditStudentDetails)
			studentDetails.POST("/upload-student-photo", studentDetailsController.UploadStudentPhoto)
		}

		timetable := main.Group("/timetable")
		{
			timetable.GET("all-study-timetables", timetableController.FetchStudyTimetables)
			timetable.POST("create-study-timetable", timetableController.CreateStudyTimetable)
			timetable.POST("add-timetable-period", timetableController.AddPeriod)
			timetable.GET("get-timetable-periods", timetableController.FetchTimetablePeriodsGroupedByDay)
			timetable.GET("get-distinct-days", timetableController.GetDistinctDays)
			timetable.GET("get-periods-by-day/:day", timetableController.GetPeriodsByDay)
			timetable.GET("subjects-by-class-level-number/:class_level_number", timetableController.FetchSubjectsByClassLevelNumber)
			timetable.GET("get-classes-by-campus-number-and-level-number", timetableController.FetchClassLevelsByCampusNumberAndLevelNumber)
			timetable.GET("fetch-levels-by-campus-number/:campus_number", timetableController.FetchLevelsByCampusNumber)
			timetable.GET("studying-timeblocks", timetableController.CreateStudyingTimeBlock)
			timetable.POST("studying-timeblocks", timetableController.FetchStudyingTimeBlocks)
			timetable.POST("studying-timetable", timetableController.CreateStudyingTimeTable)
			timetable.GET("studying-timetable", timetableController.FetchStudyingTimeTables)
			timetable.POST("study-periods", timetableController.CreateStudyPeriod)
			timetable.GET("study-periods", timetableController.FetchStudyPeriods)
		}

		stream := main.Group("/stream")
		{
			stream.POST("new", classController.AddClassStream)
			stream.PATCH("update", classController.UpdateClassStream)
			stream.DELETE(":stream_number", classController.DeleteClassStream)
			stream.GET(":stream_number", classController.GetClassStream)
			stream.GET("fetch", classController.FetchStreams)

		}

		eos := main.Group("/eos")
		{
			entity := eos.Group("/entity")
			{
				entity.POST("add", eosController.AddEntity)
				entity.POST("update", middleware.ValidatePermission(constants.CanUpdateEntity), eosController.UpdateEntity)
				entity.POST("delete", eosController.DeleteEntity, middleware.ValidatePermission(constants.CanDeleteEntity))
				entity.POST("current", eosController.SaveCurrentEntity)
				entity.GET("current", eosController.GetCurrentEntity)
			}

			eos.GET("my-entities", eosController.GetUserEntities)
			eos.GET("calendars", eosController.GetAllCalendars)
			eos.POST("calendars", eosController.AddGeneralCalendar)
			eos.GET("calendars/:calendar_number", eosController.GetCalendarDetails)
			eos.POST("calendars/update", eosController.AddCalendarEvent)
			eos.POST("calendars/add-event", eosController.AddCalendarEvent)
			eos.POST("add-user-to-entity", eosController.AddUserToEntity)
			eos.GET("calendars/events", eosController.AllCalendarEvents)
		}

		application := main.Group("/applications")
		{
			// source of school info
			application.POST("new-source-of-school-info", sourceOfSchoolController.NewSourceOfSchoolInfo)
			application.GET("source-of-school-info", sourceOfSchoolController.GetSourcesOfSchoolInfo)
			application.PUT("edit-source-of-school-info/:source_of_school_info_number", sourceOfSchoolController.EditSourceOfSchoolInfo)
			application.DELETE("source-of-school-info/:source_of_school_info_number", sourceOfSchoolController.DeleteSourceOfSchoolInfo)

			// parents
			application.POST("new-parent", parentController.CreateParent)
			application.GET("parents", parentController.FetchParents)
			application.GET("parents/:id", parentController.FetchParentInfo)
			application.PUT("parents/:id", parentController.UpdateParent)
			application.POST("parents/:id", parentController.DeleteParent)

			// students
			application.POST("/new-student", studentController.CreateStudent)
			application.GET("/students/:student_number", studentController.FetchStudentByNumber)
			application.GET("/students", studentController.FetchAllStudents)
			application.PUT("/students/:student_number", studentController.EditStudent)
			application.DELETE("/students/:student_number", studentController.DeleteStudent)

			// previous schools
			application.POST(":application_number/previous-schools", previousSchoolController.CreatePreviousSchools)
			application.PUT(":application_number/previous-schools", previousSchoolController.UpdatePreviousSchool)

			// application
			application.GET("/:application_number", applicationController.FetchApplication)
			application.GET("all", applicationController.FetchApplications)
			application.PUT("/:application_number", applicationController.EditApplication)
			application.POST("submit-application/:application_number", applicationController.SubmitApplication)
			application.DELETE("/:application_number", applicationController.DeleteApplication)
		}

		learningMaterials := main.Group("/learning-materials")
		{
			learningMaterials.POST("", learningMaterialsControllers.CreateNewLearningMaterial)
			learningMaterials.POST("sub-type", learningMaterialsControllers.CreateMaterialSubType)
			learningMaterials.GET("sub-types/:category", learningMaterialsControllers.FetchMaterialSubTypesByCategory)
			learningMaterials.POST("new-curriculum-reference", learningMaterialsControllers.CreateCurriculumReference)
			learningMaterials.GET("", learningMaterialsControllers.FetchAllLearningMaterials)
		}

		payments := main.Group("/payments")
		{
			payments.POST("", paymentsController.CreatePayment)
		}

		preAdmissionInterviewTestReport := main.Group("/pre-admission-test-interview-reports")
		{
			preAdmissionInterviewTestReport.POST("", preInterviewController.CreatePreAdmissionInterviewReport)
			preAdmissionInterviewTestReport.GET(":number", preInterviewController.FetchPreAdmissionInterviewReportByNumber)
			preAdmissionInterviewTestReport.GET("/all", preInterviewController.GetAllPreAdmissionInterviewReports)
		}

		interview := main.Group("/interviews")
		{
			interview.POST("", interviewController.CreateInterview)
			interview.GET("/all", interviewController.FetchAllInterviews)
			interview.GET("details", interviewController.FetchAnInterview)
			interview.DELETE("delete/:interview_number", interviewController.RemoveInterview)
			interview.PUT("", interviewController.EditInterview)
			interview.GET("/students-by-interview", interviewController.FetchStudentsByInterviewNumber)
			interview.POST("add-student-to-interview", interviewController.AddStudentToInterview)
			interview.GET("/:interview_number/questions", interviewController.FetchQuestionsByInterviewNumber)
			interview.POST("/add-interview-assessment-metric", interviewController.CreateInterviewAssessmentMetric)
			interview.GET("/:interview_number/interview-assessment-metrics", interviewController.FetchInterviewAssessmentMetrics)
			interview.GET("student-interview-assessment-results", interviewController.FetchInterviewAssessmentMetricByStudent)
			interview.POST("remove-student-from-interview", interviewController.RemoveStudentFromInterview)
			interview.POST("remove-interview-assessment-metric", interviewController.RemoveInterviewAssessmentMetric)
		}

		interviewQuestion := main.Group("/interview-questions")
		{
			interviewQuestion.POST("", interviewController.CreateInterviewQuestion)
			interviewQuestion.GET("/:number", interviewController.FetchInterviewQuestion)
			interviewQuestion.GET("", interviewController.FetchAllInterviewQuestions)
			interviewQuestion.PUT("/:number", interviewController.EditInterviewQuestion)
		}

		interviewResult := main.Group("/interview-results")
		{
			interviewResult.POST("/:student_number", interviewController.CreateInterviewResults)
			interviewResult.GET("/:student_number", interviewController.FetchInterviewResultByStudentNumber)
			interviewResult.GET("/all", interviewController.FetchAllInterviewResults)
			interviewResult.DELETE(":number", interviewController.DeleteInterviewResult)
		}

		invoices := main.Group("/invoices")
		{
			invoices.POST("", invoiceController.CreateInvoice)
			invoices.GET("", invoiceController.GetAllInvoices)
			invoices.POST("pay", invoiceController.CreateInvoicePayment)
			invoices.GET("details", invoiceController.GetInvoiceDetails)
			invoices.GET("financial-summary", invoiceController.FinancialSummary)
			invoices.POST("update", invoiceController.UpdateInvoice)
			invoices.GET("general-statement", invoiceController.GeneralStatement)

			payments := invoices.Group("/payments")
			{
				payments.GET("statements", invoiceController.AllStatements)
				payments.GET("credits", invoiceController.TransactionCredits)
				payments.GET("", invoiceController.GetAllPayments)
				payments.POST("cancel", invoiceController.CancelPayment)
				payments.POST("delete", invoiceController.DeletePayment)
			}

			invoices.POST("delete", invoiceController.DeleteInvoice)
			invoices.POST("cancel", invoiceController.CancelInvoice)

			billable := invoices.Group("billable")
			{
				billable.POST("", invoiceController.CreateBillable)
				billable.PUT("", invoiceController.UpdateBillable)
				billable.DELETE("/:billable_number", invoiceController.DeleteBillable)
				billable.GET("", invoiceController.AllBillable)
				billable.GET("details", invoiceController.BillableDetails)
				billable.GET("determine", invoiceController.GetBill)
			}
		}

		services := main.Group("/services")
		{
			services.POST("", serviceController.CreateService)
			services.GET("", serviceController.FetchServices)
			services.DELETE("/:service_number", serviceController.DeleteService)
			services.PUT(":service_number", serviceController.EditService)
		}

		products := main.Group("/products")
		{
			products.POST("", productController.SaveProduct)
			products.PUT(":product_number", productController.EditProduct)
			products.GET("", productController.FetchProducts)
			products.DELETE(":product_number", productController.DeleteProduct)
		}

		items := main.Group("/items")
		{
			items.GET("", serviceController.FetchAllItems)
		}

		admission := main.Group("/admissions")
		{
			admission.POST("", admissionController.CreateAdmission)
			admission.POST("existing-admission", admissionController.CreateExistingAdmission)
			admission.GET("", admissionController.FetchAdmissions)
			admission.GET("/:admission_number", admissionController.FetchAdmissionById)
			admission.GET("items/:admission_number", admissionController.FetchItemsAddedOnAdmission)

		}

		assessmentArea := main.Group("/assessment-areas")
		{
			assessmentArea.POST("", assessmentAreaController.CreateAssessmentArea)
			assessmentArea.GET("", assessmentAreaController.FetchAssessmentAreas)
			assessmentArea.GET("by-campus-class-number/:campus_class_number", assessmentAreaController.FetchAssessmentAreasByCampusClassNumber)
			assessmentArea.PUT("", assessmentAreaController.EditAssessmentArea)
			assessmentArea.DELETE("/:assessment_area_number", assessmentAreaController.DeleteAssessmentArea)
		}

		assessment := main.Group("/assessments")
		{
			assessment.POST("", assessmentController.CreateAssessment)
			assessment.GET("", assessmentController.FetchAssessmentByStudent)

			assessment.POST("assessment-metrics", assessmentController.CreateAssessmentMetric)
			assessment.DELETE("assessment-metrics/:assessment_metric_number", assessmentController.DeleteAssessmentMetric)
			assessment.POST("update", assessmentController.UpdateAssessmentMetric)
			assessment.GET("assessment-metrics", assessmentController.FetchAssessmentMetrics)
			assessment.POST("script-marking-guide", assessmentController.ScriptMarkingGuide)
			assessment.GET("script-marking-guide", assessmentController.GetScriptMarkingGuide)

			assessment.POST("metrics", assessmentController.CreateMetric)
			assessment.GET("metrics", assessmentController.FetchMetrics)
			assessment.DELETE("metrics/:metric_number", assessmentController.DeleteMetric)
			assessment.PUT("metrics/:metric_number", assessmentController.EditMetric)

			aggregateScale := assessment.Group("/aggregate-scales")
			{
				aggregateScale.POST("", aggregateScaleController.CreateAggregateScale)
				aggregateScale.GET("", aggregateScaleController.FetchAllScales)
				aggregateScale.GET("all", aggregateScaleController.FetchAllDetailedScales)
				aggregateScale.POST("update", aggregateScaleController.UpdateAggregateScale)
				aggregateScale.POST("remove/:scale_number", aggregateScaleController.RemoveAggregateScale)
				aggregateScale.POST("upload-illustration", aggregateScaleController.UploadAggregateScaleIllustration)
			}

			assessmentGrade := assessment.Group("/assessment-grades")
			{
				//assessmentGrade.POST("", assessmentGradeController.CreateAssessmentGrade)
				assessmentGrade.GET("", assessmentGradeController.FetchAssessmentGrades)
				//assessmentGrade.PUT("/update", assessmentGradeController.EditAssessmentGrade)
				//assessmentGrade.DELETE("/delete/:assessment_grade_number", assessmentGradeController.DeleteAssessmentGrade)
			}

			gradeGroup := assessment.Group("/grade-groups")
			{
				gradeGroup.POST("", gradesController.CreateGradeGroup)
				gradeGroup.GET("", gradesController.FetchGradeGroups)
				gradeGroup.POST("update", gradesController.UpdateGradeGroup)
				gradeGroup.POST("remove/:number", gradesController.DeleteGradeGroup)
			}

			grades := assessment.Group("/grades")
			{
				grades.POST("", gradesController.CreateGrade)
				grades.GET("", gradesController.FetchGrades)
				grades.POST("update", gradesController.UpdateGrade)
				grades.POST("remove/:grade_number", gradesController.RemoveGrade)
			}

		}

		attendance := main.Group("/attendances")
		{
			attendance.POST("", controllers4.CreateAttendance)
			attendance.GET("", controllers4.FetchAttendance)
			attendance.POST("checkout", controllers4.CheckOut)
		}

		settings := main.Group("settings")
		{
			permissions := settings.Group("permissions")
			{
				permissions.GET("groups", authController.FetchPermissionGroups)
			}

			roles := settings.Group("roles")
			{
				roles.POST("", permissionController.CreateRole)
				roles.GET("", permissionController.FetchRoles)
				roles.GET(":role_number", permissionController.FetchRoleDetails)
				roles.PUT("assign", permissionController.AssignRolesToAUser)
				roles.POST("remove", permissionController.RemoveRoleAssignment)
				roles.POST("update/:role_number", permissionController.UpdateRole)
				roles.DELETE("delete/:role_number", permissionController.DeleteRole)

			}

			scopes := settings.Group("scopes")
			{
				scopes.POST("", permissionController.CreateScope)
				scopes.GET("", permissionController.FetchScopes)
			}
		}

		users := main.Group("users")
		{
			users.GET("", authController.FetchUsers)
			users.GET("all", authController.FetchAllUsers)
			users.GET(":user_number", authController.FetchUser)
			users.POST("", authController.CreateUser)
			users.GET("meals", cateringController.FetchUserMeals)
		}

		reportCards := main.Group("report-cards")
		{
			reportCards.POST("", progressReportController.CreateProgressReport)
			reportCards.GET("", progressReportController.FetchAllProgressReports)
		}

		gradingStandard := main.Group("grading-standard")
		{
			gradingStandard.GET("", aggregateScaleController.FetchGradingStandards)
			gradingStandard.GET("classes", aggregateScaleController.ClassGradingStandards)
			gradingStandard.POST("", aggregateScaleController.CreateGradingStandard)
			gradingStandard.POST("update", aggregateScaleController.UpdateGradingStandard)
			gradingStandard.POST("status", aggregateScaleController.UpdateStatusGradingStandard)
			gradingStandard.POST("remove/:grading_standard_number", aggregateScaleController.RemoveGradingStandard)
			gradingStandard.POST("attach-assessment-area", aggregateScaleController.AttachAssessmentAreas)
			gradingStandard.GET("assessment-areas", aggregateScaleController.FetchAssessmentAreasByGradingStandard)
			gradingStandard.GET("details", aggregateScaleController.DetailsGradingStandard)
			gradingStandard.POST("default-for-class", aggregateScaleController.MakeGradingStandardDefault)
		}

		enrollment := main.Group("enrollment")
		{
			enrollment.POST("", enrollmentController.CreateEnrollment)
			enrollment.GET("", enrollmentController.FetchEnrolledStudents)
			enrollment.POST("terminate-student/:student_number", enrollmentController.TerminateStudentEnrollment)
			enrollment.POST("re-enroll-student/:student_number", enrollmentController.ReEnrollStudent)
		}

		staff := main.Group("staff-management")
		{
			staff.POST("attach-staff-to-stream", staffController.AttachStaffToStream)
			staff.GET("", staffController.GetStaffOnStream)
			staff.DELETE("", staffController.DetachStaffFromStream)

			salaryScale := staff.Group("salary-scale")
			{
				salaryScale.POST("", staffController.CreateSalaryScale)
				salaryScale.GET("", staffController.FetchSalaryScales)
				salaryScale.PUT("", staffController.EditSalaryScale)
			}

			tasks := staff.Group("tasks")
			{
				tasks.POST("", staffController.CreateTask)
				tasks.PUT("", staffController.EditTask)
				tasks.DELETE("", staffController.DeleteTask)
				tasks.GET("", staffController.FetchTasks)
			}

			taskSupervisions := staff.Group("task-supervisions")
			{
				taskSupervisions.POST("", staffController.CreateTaskSupervision)
				taskSupervisions.GET("", staffController.FetchTaskSUpervisions)
			}

		}

		reportSettings := main.Group("report-settings")
		{
			reportSettings.POST("", reportSettingsController.CreateReportSettings)
			reportSettings.DELETE("", reportSettingsController.DeleteReportSetting)
			reportSettings.PUT("update", reportSettingsController.UpdateReportSettings)
			reportSettings.GET("", reportSettingsController.FetchReportSettings)
			reportSettings.POST("columns", reportSettingsController.CreatReportColumn)
			reportSettings.GET("columns", reportSettingsController.FetchReportColumns)
			reportSettings.DELETE("columns/:report_column_number", reportSettingsController.DeleteReportColumns)
			reportSettings.POST("report-column-assessment", reportSettingsController.CreateReportColumnAssessment)
			reportSettings.PUT("report-column-assessment/:report_column_assessment_number", reportSettingsController.UpdateReportColumnAssessment)
			reportSettings.DELETE("report-column-assessment", reportSettingsController.DeleteReportColumnAssessment)
			reportSettings.GET("report-column-assessment", reportSettingsController.AllReportColumnAssessments)
			reportSettings.GET("report-summary", reportSettingsController.FetchReportSummaries)
			reportSettings.POST("report-score", reportSettingsController.CreateReportScore)
			reportSettings.GET("report-score-summary", reportSettingsController.FetchReportScoreSummaries)
			reportSettings.GET("report-settings-by-student", reportSettingsController.FetchReportSummaryByStudentAndStudyPeriodAndClass)
			reportSettings.POST("generate-report", reportSettingsController.GenerateReportSummary)
			reportSettings.GET("report-summary-details/:report_summary_number", reportSettingsController.FetchReportSummaryDetails)
			reportSettings.GET("report-summary-promotion-details/:report_summary_number", reportSettingsController.FetchReportPromotionDetails)
			reportSettings.POST("add-head-teachers-comment/:report_setting_number", reportSettingsController.CreateHeadTeacherComment)
			reportSettings.GET("head-teachers-comments", reportSettingsController.FetchHeadTeacherComments)
			reportSettings.POST("add-teacher-comment", reportSettingsController.CreateTeachersComment)
			reportSettings.POST("learning-area-comments", reportSettingsController.CreateLearningAreaComment)
			reportSettings.DELETE("head-teacher-comment/:id", reportSettingsController.DeleteHeadTeacherComment)

			narrative := reportSettings.Group("narrative")
			{
				narrative.GET("", assessmentController.GetNarrative)
				narrative.PATCH(":narrative_number", assessmentController.UpdateNarrative)
				narrative.GET(":narrative_number", assessmentController.GetDetailedNarrative)
				narrative.DELETE(":narrative_number", assessmentController.DeleteNarrative)
				narrative.POST("", assessmentController.SaveNarrative)

				scheduled := narrative.Group("scheduled")
				{
					scheduled.PATCH(":schedule_number", assessmentController.UpdateNarrativeSchedule)
					scheduled.POST("", assessmentController.SaveNarrativeSchedule)
					scheduled.GET("", assessmentController.NarrativeSchedule)
					scheduled.DELETE(":schedule_number", assessmentController.DeleteNarrativeSchedule)
				}

			}

		}

		promotions := main.Group("promotions")
		{
			promotions.POST("", promotionController.PromoteStudent)
			promotions.GET("", promotionController.FetchAllPromotions)
			promotions.PUT(":promotion_number", promotionController.EditPromotion)

			promotionSettings := promotions.Group("settings")
			{
				promotionSettings.POST("", promotionController.CreatePromotionSettings)
				promotionSettings.PUT(":promotion_setting_number", promotionController.UpdatePromotionSettings)
				promotionSettings.GET("", promotionController.FetchPromotionSettings)
				promotionSettings.GET("details", promotionController.FetchPromotionSettingDetails)
				promotionSettings.GET("grades", promotionController.FetchPromotionSettingGrades)
			}

			promotionSettingsValues := promotions.Group("settings-values")
			{
				promotionSettingsValues.POST("", promotionController.CreatePromotionSettingValue)
				promotionSettingsValues.PUT("", promotionController.UpdatePromotionSettingValue)
				promotionSettingsValues.GET("", promotionController.FetchPromotionSettingsValues)
			}
		}

		previousSchoolResults := main.Group("previous-school-results")
		{
			previousSchoolResults.POST("", previousSchoolResultsController.CreatePreviousSchoolResults)
			previousSchoolResults.GET("all", previousSchoolResultsController.FetchAllPreviousSchoolResults)
			previousSchoolResults.GET("", previousSchoolResultsController.FetchPreviousSchoolResultsByStudent)
		}

		students := main.Group("students")
		{
			students.GET("/by-assessment-metric/:assessment_metric_number", studentController.FetchStudentsByAssessmentMetricNumber)
			students.POST("medical-records", studentController.CreateMedicalRecord)
			students.POST("family-structure", studentController.CreateFamilyStructure)
			students.POST("emergency-information", studentController.CreateEmergencyInformation)
		}

		catering := main.Group("catering")
		{
			// Feeding sessions
			feedingSession := catering.Group("feeding-sessions")
			{
				feedingSession.POST("", cateringController.CreateFeedingSession)
				feedingSession.GET("", cateringController.FetchFeedingSessions)
				feedingSession.POST("update", cateringController.UpdateFeedingSession)
				feedingSession.POST("delete", cateringController.DeleteFeedingSession)
			}

			// Feeding session per day
			catering.POST("feeding-sessions-per-day", cateringController.CreateFeedingSessionPerDay)
			catering.GET("feeding-sessions-per-day", cateringController.FetchFeedingSessionsPerDay)

			// Feeding session items
			catering.POST("feeding-sessions-items", cateringController.CreateFeedingSessionItemController)
			catering.GET("feeding-sessions-items", cateringController.FetchFeedingSessionItems)

			// Feeding groups
			feedingGroups := catering.Group("feeding-groups")
			{
				feedingGroups.POST("", cateringController.CreateFeedingGroup)
				feedingGroups.GET("", cateringController.FetchFeedingGroups)
				feedingGroups.POST("update", cateringController.UpdateFeedingGroup)
				feedingGroups.POST("delete", cateringController.DeleteFeedingGroup)
			}

			// Feeding groups per day
			catering.POST("feeding-groups-per-day", cateringController.CreateFeedingGroupPerDay)
			catering.GET("feeding-groups-per-day", cateringController.FetchFeedingGroupsPerDay)

			// Menu items

			menuItems := catering.Group("menu-items")
			{
				menuItems.POST("", cateringController.CreateMenuItem)
				menuItems.GET("", cateringController.FetchMenuItems)
				menuItems.POST("update", cateringController.EditMenuItem)
				menuItems.POST("delete", cateringController.DeleteMenuItem)
			}

			catering.GET("menu-items-paginated", cateringController.FetchMenuItemsPaginated)

			// Units of measure
			catering.POST("units-of-measure", cateringController.CreateUnitOfMeasure)
			catering.GET("units-of-measure", cateringController.FetchUnitsOfMeasure)
			catering.PUT("units-of-measure", cateringController.EditUnitOfMeasure)
			catering.DELETE("units-of-measure/:unit_of_measure_number", cateringController.RemoveUnitOfMeasure)

			// Item assembly
			itemAssembly := catering.Group("item-assembly")
			{
				itemAssembly.POST("", cateringController.CreateItemAssembly)
				itemAssembly.GET("", cateringController.FetchItemAssemblies)
				itemAssembly.POST("update", cateringController.EditItemAssembly)
			}

			// Feedings
			catering.POST("feedings", cateringController.CreateFeeding)
			catering.GET("feedings", cateringController.FetchFeedings)

			// Stocking
			catering.POST("stocking", cateringController.CreateStocking)
			catering.GET("stocking", cateringController.FetchStockings)
			catering.GET("stocking/:stock_item_number", cateringController.FetchItemStock)

			// Feeding group users

			feedingGroupUsers := catering.Group("feeding-group-users")
			{
				feedingGroupUsers.POST("", cateringController.CreateFeedingGroupUser)
				feedingGroupUsers.GET("", cateringController.FetchFeedingGroupUsers)
				feedingGroupUsers.GET("available", cateringController.FetchAvailableUsers)
			}

			// Feeding group session
			catering.POST("feeding-group-sessions", cateringController.CreateFeedingGroupSession)
			catering.POST("feeding-group-calendar", cateringController.ActivateFeedingGroupOnCalendars)
			catering.GET("feeding-group-calendar", cateringController.FetchFeedingGroupOnCalendars)
			catering.GET("feeding-group-sessions", cateringController.FetchFeedingGroupSessions)

			// Feeding per group
			catering.POST("meal-courses-per-group", cateringController.AddMealCoursesToFeedingGroup)
			catering.GET("meal-courses-per-group", cateringController.FetchFeedingsPerGroup)

			// Meals
			catering.POST("meals", cateringController.AddMeal)
			catering.GET("meals", cateringController.FetUsersMeals)
			catering.GET("meal-projections", cateringController.MealsProjections)
			catering.GET("meal-calendar", cateringController.MealCourseCalendar)
			catering.POST("meal-course-feeding-sessions", cateringController.FeedingSessionsToMealCourse)
			catering.GET("meal-course-feeding-sessions", cateringController.AllMealCourseFeedingSessions)
			catering.POST("activate-meal-course", cateringController.ActivateMealCourse)

			cateringUpdate := catering.Group("meal-courses")
			{
				cateringUpdate.POST("", cateringController.SaveMealCourse)
				cateringUpdate.GET("", cateringController.GetMealCourses)
				cateringUpdate.POST("update", cateringController.UpdateSaveMealCourse)
				cateringUpdate.POST("update/item", cateringController.UpdateMealCourseItem)
			}

		}

		statistics := main.Group("statistics")
		{
			statistics.GET("summary", statisticsController.FetchStatistics)
			statistics.POST("print", statisticsController.PrintReport)
			statistics.POST("test", statisticsController.PrintTestReport)
		}

		activityReporting := main.Group("activity-reporting")
		{
			indicators := activityReporting.Group("indicators")
			{
				indicators.POST("", activityReportingControllers.SaveIndicator)
				indicators.GET("", activityReportingControllers.GetIndicators)
			}
		}

		landingPage := main.Group("landing-page")
		{
			landingPage.POST("new", controllers2.CreateLandingPage)
			landingPage.GET("all", controllers2.GetLandingPage)
		}

		curriculumDelivery := main.Group("curriculum-delivery")
		{
			curriculumDelivery.POST("add-session", controllers3.AddSession)
			curriculumDelivery.GET("sessions", controllers3.GetSessions)
			curriculumDelivery.GET("lessons", controllers3.GetLessons)
			curriculumDelivery.GET("students/:curriculum_delivery_number", controllers3.GetStudentsBySession)
			curriculumDelivery.PUT("update-session", controllers3.UpdateSession)
			curriculumDelivery.GET("percentage-by-stream/:stream_number", controllers3.GetCurriculumDeliveryProgressByStream)
		}

	}
	attachments := router.Group("attachments")
	{
		attachments.GET("download/:attachment_number", controller.DownloadAttachment)
		attachments.POST("delete", controller.DeleteAttachment)
	}

	return router
}
