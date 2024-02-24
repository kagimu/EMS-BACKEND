package routes

import (
	"time"

	productController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/item/product/controllers"
	serviceController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/item/service/controllers"

	academicYearController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/academicYear/controllers"
	admissionController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/controllers"
	applicationController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/application/controllers"
	assessmentController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/controllers"
	authController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/controllers"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/middleware"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/controllers"
	classController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/controllers"
	assessmentAreaController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/controllers"
	conceptController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/concept/controllers"
	curricullumController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/importCurriculum/controllers"
	learningAreaController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/learningarea/controllers"
	strandController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/strand/controllers"
	subjectClusterController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/subjectcluster/controllers"
	subTopicController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/subtopic/controllers"
	topicController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/topic/controllers"
	eosController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/controllers"
	interviewController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/interview/controllers"
	invoiceController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/controllers"
	invoicePaymentController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoicePayment/controllers"
	lessonPlanControllers "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/lessonPlan/controllers"
	parentController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/parent/controllers"
	paymentsController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/payments/controllers"
	preInterviewController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/preAdmissionInterviewTest/controllers"
	previousSchoolController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/previousSchool/controllers"
	schemeOfWorkControllers "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/schemeOfWork/controllers"
	schoolController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/controllers"
	learningMaterialsControllers "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/learningMaterials/controllers"
	learningLevelControllers "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/level/controllers"
	sourceOfSchoolController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/sourceOfSchoolInfo/controllers"
	studentController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/controllers"
	studentDetailsController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studentDetails/controllers"
	studyPeriodController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/controllers"
	subjectController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/subject/controllers"
	timetableController "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/controllers"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func SetupRoute() *gin.Engine {
	router := gin.Default()

	//setting cors in the system
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:2054", "https://wholesome-api.noregrets.ug"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Access-Control-Allow-Origin", "Origin", "Content-Length", "Content-Type", "User-Agent", "Referrer", "Host", "Token", "Authorization", "Accept-Encoding", "X-CSRF-Token", "Access-Control-Allow-Headers", "Content-Type", "user-number"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			///return origin == "https://github.com"
			return true
		},
		MaxAge: 12 * time.Hour,
	}))

	//this is to validate user toke before accessing any route
	//router.Use(middlewares.JWTMiddleware())
	router.MaxMultipartMemory = 20 << 20 // 20 MiB
	//router.Static("/", "./fileRepo")
	router.Use(static.Serve("/fileRepo", static.LocalFile("/fileRepo", true)))

	auth := router.Group("/auth")
	{
		auth.POST("login", authController.Login)
		auth.POST("register", authController.Register)
		auth.POST("set-password", authController.UserPassword)
		auth.POST("reset-password", authController.ResetPassword)
	}

	main := router.Group("/main", middleware.AuthorizeJWT())
	{
		//main.GET("profile/:id", middlewares.JWTMiddleware(), controllers.Profile)
		main.GET("profile", authController.Profile)

		//manage curriculum
		curriculum := main.Group("/curriculum")
		{
			curriculum.GET("/info", curricullumController.CurriculumInfo)
			curriculum.POST("/edit", curricullumController.CurriculumEdit)
			curriculum.POST("/delete", curricullumController.DeleteCurriculum)
		}
		main.POST("new-curriculum", curricullumController.NewCurriculum)
		main.GET("curriculum-list", curricullumController.FetchCurriculum)
		main.GET("curriculum-info", curricullumController.FetchCurriculumInfo)
		main.POST("upload-curriculum", curricullumController.ImportCurriculum)
		main.GET("uploaded-curriculum", curricullumController.FetchUploadedCurriculum)
		main.GET("fetch-school-curriculum", curricullumController.FetchSchoolCurriculum)
		main.GET("fetch-school-curriculum-level", curricullumController.FetchSchoolCurriculumLevel)

		//manage School
		school := main.Group("/school")
		{
			school.GET("/all", schoolController.FetchSchools)
			school.POST("/new", schoolController.NewSchool)
			school.POST("/edit", schoolController.EditSchool)
			school.GET("/info", schoolController.FetchSchoolInfo)
			school.POST("/delete", schoolController.DeleteSchool)
		}

		//manage campus
		campus := main.Group("/campus")
		{
			campus.GET("/all", controllers.FetchCampuses)
			campus.GET("/info", controllers.CampusInfo)
			campus.GET("/curriculum", controllers.CampusCurriculum)
			campus.POST("/new", controllers.NewCampus)
			campus.POST("/edit", controllers.EditCampus)
			campus.POST("/delete", controllers.DeleteCampus)
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

		//level
		main.POST("new-level", learningLevelControllers.NewLevel)
		main.POST("new-curriculum-level", learningLevelControllers.NewCurriculumLevel)
		main.GET("fetch-level", learningLevelControllers.FetchLearningLevel)
		main.GET("fetch-curriculum-level", learningLevelControllers.FetchCurriculumLevel)

		//	manage class
		class := main.Group("/class")
		{
			class.GET("/info", classController.ClassInfo)
			class.GET("by-campus-number/:campus_number", classController.FetchClassesByCampusNumber)
			class.POST("/edit", classController.EditClass)
			class.POST("/delete", classController.DeleteClass)
		}

		main.POST("new-class", classController.NewClass)
		main.POST("new-campus-class", classController.AddCampusClass)
		main.POST("new-curriculum-class", classController.AddCurriculumClass)
		main.GET("classes", classController.FetchClasses)
		main.GET("fetch-curriculum-classes", classController.FetchCurriculumClasses)
		main.GET("filter-curriculum-classes", classController.FilterCurriculumClasses)
		main.GET("scheme-of-work-classes", classController.FetchSchemeOfWorkClasses)
		main.GET("classes-curriculum", classController.FetchClassesCurriculum)
		main.GET("classes-to-add", classController.FetchClassesToAdd)
		main.GET("class-level", classController.FetchClassLevel)
		main.GET("class-subjects", classController.FetchClassSubjects)
		main.POST("subject-to-class", classController.NewClassSubjects)

		//manage subject
		subject := main.Group("/subject")
		{
			subject.GET("/info", subjectController.SubjectInfo)
			subject.POST("/edit", subjectController.SubjectEdit)
			subject.POST("/delete", subjectController.DeleteSubject)
		}

		main.POST("new-subject", subjectController.NewSubject)
		main.GET("subjects", subjectController.FetchSubject)
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
		main.POST("new-academic-year", academicYearController.NewAcademicYear)
		main.GET("fetch-academic-year", academicYearController.FetchAcademicYearInfo)
		main.GET("academic-by-class", academicYearController.FetchAcademicByClass)

		//adds study period
		main.POST("new-study-period", studyPeriodController.NewStudyPeriod)
		main.GET("fetch-study-period", studyPeriodController.FetchStudyPeriod)

		//adds scheme of work
		schemeOfWork := main.Group("/scheme-of-work")
		{
			schemeOfWork.GET("download-scheme-template", schemeOfWorkControllers.DownloadSchemeOfWork)
			schemeOfWork.POST("upload-scheme-of-work", schemeOfWorkControllers.ImportSchemeOfWork)
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
			lessonPlan.GET("/fetch-lesson-plan", lessonPlanControllers.FetchLessonPlan)
			lessonPlan.GET("/references/:id", lessonPlanControllers.FetchLessonPlanReferences)
		}

		// Student details
		studentDetails := main.Group("/student-details")
		{
			studentDetails.POST("/:student_number", studentDetailsController.CreateStudentDetails)
			studentDetails.PUT("/:student_number", studentDetailsController.EditStudentDetails)
		}

		//timetable routes
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
		}

		stream := main.Group("/stream")
		{
			stream.POST("new", classController.AddClassStream)
			stream.GET("fetch", classController.FetchStreamsByClassLevelNumber)

		}

		eos := main.Group("/eos")
		{
			eos.POST("entity/add", eosController.AddEntity)
			eos.GET("calendars", eosController.GetAllCalendars)
			eos.POST("calendars", eosController.AddGeneralCalendar)
			eos.GET("calendars/:calendar_number", eosController.GetCalendarDetails)
			eos.POST("calendars/update", eosController.AddCalendarEvent)
			eos.POST("calendars/add-event", eosController.AddCalendarEvent)
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
			application.PATCH("parents/:id", parentController.UpdateParent)
			application.POST("parents/:id", parentController.DeleteParent)

			// students
			application.POST("/new-student", studentController.CreateStudent)
			application.GET("/students/:apn", studentController.FetchStudentByAPN)
			application.GET("/students", studentController.FetchAllStudents)
			application.PUT("/students/:apn", studentController.EditStudent)
			application.POST("/students/:id", studentController.DeleteStudent)

			// previous schools
			application.POST(":application_number/previous-schools", previousSchoolController.CreatePreviousSchools)

			// application
			application.GET("/:application_number", applicationController.FetchApplication)
			application.GET("all", applicationController.FetchApplications)
			application.GET("", applicationController.FetchApplicationsByStatus)
			application.PUT("/:application_number", applicationController.EditApplication)
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

		preAdmissionInterviewTest := main.Group("/pre-admission-interview-tests")
		{
			preAdmissionInterviewTest.POST("", preInterviewController.CreatePreAdmissionInterviewTest)
		}

		interview := main.Group("/interviews")
		{
			interview.POST("", interviewController.CreateInterview)
			interview.GET("/all", interviewController.FetchAllInterviews)
			interview.GET("/:interview_number", interviewController.FetchAnInterview)
			interview.PUT("/:interview_number", interviewController.EditInterview)
			interview.GET("/:interview_number/students", interviewController.FetchStudentsByInterviewNumber)
			interview.POST("add-student-to-interview", interviewController.AddStudentToInterview)
			interview.GET("/:interview_number/questions", interviewController.FetchQuestionsByInterviewNumber)
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
			invoices.PUT(":student_number", invoiceController.UpdateInvoice)
		}
		invoicePayments := main.Group("/invoice-payments")
		{
			invoicePayments.POST("", invoicePaymentController.CreateInvoicePayment)
		}

		services := main.Group("/services")
		{
			services.POST("", serviceController.CreateService)
			services.GET("", serviceController.FetchServices)
			services.DELETE("", serviceController.DeleteService)
			services.PUT(":service_number", serviceController.EditService)
		}

		products := main.Group("/products")
		{
			products.POST("", productController.CreateProduct)
			products.PUT(":product_number", productController.EditProduct)
			products.GET("", productController.FetchProducts)
			products.DELETE(":product_number", productController.DeleteProduct)
		}

		admission := main.Group("/admissions")
		{
			admission.POST("", admissionController.CreateAdmission)
			admission.PUT("/:number/admit", admissionController.AdmitStudent)
			admission.PUT("/admit", admissionController.BulkAdmitStudents)

			offer := admission.Group("/admission-offers")
			{
				offer.POST("", admissionController.CreateAdmissionOffer)
				offer.GET("", admissionController.FetchAdmissionOffers)
				offer.GET("/:admission_offer_number", admissionController.FetchAnAdmissionOffer)
				offer.GET("/student/:student_number", admissionController.FetchAdmissionOfferByStudentNumber)

			}

			conditions := admission.Group("/admission-offer-conditions")
			{
				conditions.POST("", admissionController.CreateAdmissionOfferCondition)
				conditions.GET("", admissionController.FetchAdmissionOfferConditions)

			}

		}

		assessmentArea := main.Group("/assessment-areas")
		{
			assessmentArea.POST("", assessmentAreaController.CreateAssessmentArea)
			assessmentArea.GET("", assessmentAreaController.FetchAssessmentAreas)
			assessmentArea.PUT("", assessmentAreaController.EditAssessmentArea)
		}

		assessment := main.Group("/assessments")
		{
			assessment.POST("", assessmentController.CreateAssessment)
			assessment.GET(":student_number", assessmentController.FetchAssessmentByStudentNumber)
			assessment.GET("", assessmentController.FetchAssessments)

			assessment.POST("examination-types", assessmentController.CreateExaminationType)
		}

		settings := main.Group("settings/permissions")
		{
			settings.POST("", authController.CreatePermission)
			settings.GET("", authController.FetchPermissions)
			settings.GET("groups", authController.FetchPermissionGroups)
		}

		roles := main.Group("settings/roles")
		{
			roles.POST("", authController.CreateRole)
			roles.GET("", authController.FetchRoles)
			roles.GET(":role_number", authController.FetchRoleDetails)
			roles.PUT("assign", authController.AssignRolesToAUser)
			roles.PUT("update-assignment", authController.RemoveRoleAssignment)
		}
	}

	return router
}
