package services

import (
	"errors"
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	repositories4 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/academicYear/repositories"
	aggregateSaceltModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/aggregateScale/models"
	assesmentGradeModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/assessmentGrade/models"
	gradesModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/grades/models"
	repository2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/grades/repository"
	assesmentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/models"
	repositories2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/repositories"
	repositories6 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/repositories"
	classService "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/services"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	assessmentAreaRepo "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/repositories"
	assessmentAreaRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/repositories"
	repositories8 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/learningarea/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/enrollment/repository"
	models2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/promotion/models"
	repositories10 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/promotion/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/reporsitory"
	reportTemplateModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportTemplate/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportTemplate/services"
	repositories7 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/repositories"
	services4 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/services"
	repositories5 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/repositories"
	repositories9 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/repositories"
	repositories3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/repositories"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"strconv"
	"time"
)

var enrolmentRepository = repository.NewEnrolmentRepository()

func SaveReportSettings(reportSettingsRequest requests.ReportSettingsRequest, currentUserNumber string) error {

	campusClass, err := classService.ResolveCampusClass(reportSettingsRequest.ClassNumber, reportSettingsRequest.CampusNumber, reportSettingsRequest.CurriculumNumber, reportSettingsRequest.CampusCurriculumNumber)
	if err != nil {
		return err
	}

	existingReportSettings, err := reporsitory.FindDuplicateReportSettings(reportSettingsRequest.StudyPeriodNumber, reportSettingsRequest.AcademicYearNumber, campusClass.Number, constants.StatusActive)

	if err != nil {
		return errors.New("error checking report duplicate: " + err.Error())
	}

	if existingReportSettings.Number != "" {
		return errors.New("Report settings for this class  already exists")
	}

	// Save report settings
	reportSettings := models.ReportSetting{
		Number:             uuid.NewV4().String(),
		StudyPeriodNumber:  reportSettingsRequest.StudyPeriodNumber,
		AcademicYearNumber: reportSettingsRequest.AcademicYearNumber,
		CampusClassNumber:  campusClass.Number,
		Columns:            reportSettingsRequest.Columns,
		Status:             constants.StatusActive,
		UserNumber:         currentUserNumber,
		Code:               reportSettingsRequest.Code,
	}

	if err = reporsitory.SaveReportSettings(reportSettings); err != nil {
		return errors.New("Failed to create report settings" + err.Error())
	}

	return nil
}

func UpdateReportSettings(reportSettingsRequest requests.UpdateReportSettingsRequest, currentUserNumber string) error {

	reportSettings, err := reporsitory.GetReportSettingByNumber(reportSettingsRequest.ReportSettingNumber)
	if err != nil {
		return errors.New("error checking report setting by number error: " + err.Error())
	}

	if reportSettings.Number == "" {
		return errors.New("Invalid Report Settings Number ")
	}

	campusClass, err := classService.ResolveCampusClass(reportSettingsRequest.ClassNumber, reportSettingsRequest.CampusNumber, reportSettingsRequest.CurriculumNumber, reportSettingsRequest.CampusCurriculumNumber)
	if err != nil {
		return err
	}

	existingReportSettings, err := reporsitory.FindDuplicateReportSettings(reportSettingsRequest.StudyPeriodNumber, reportSettingsRequest.AcademicYearNumber, campusClass.Number, constants.StatusActive)

	if err != nil {
		return errors.New("error checking report duplicate: " + err.Error())
	}

	if existingReportSettings.Number != "" && existingReportSettings.Number != reportSettings.Number {
		return errors.New("Report settings for this class  already exists")
	}

	reportColums, err := reporsitory.FindReportColumnsByReportSettingsNumber(reportSettings.Number)
	if err != nil {
		return err
	}

	if (reportSettingsRequest.Columns < reportSettings.Columns) && (len(reportColums) > 0) {
		return errors.New("please delete a column you dont need and adjust the percentages instead of trying to reduce the colum number")
	}

	reportSettings.StudyPeriodNumber = reportSettingsRequest.StudyPeriodNumber
	reportSettings.AcademicYearNumber = reportSettingsRequest.AcademicYearNumber
	reportSettings.CampusClassNumber = campusClass.Number
	reportSettings.Columns = reportSettingsRequest.Columns
	reportSettings.Status = constants.StatusActive
	reportSettings.UserNumber = currentUserNumber
	reportSettings.Code = reportSettingsRequest.Code

	if err = reporsitory.SaveReportSettings(reportSettings); err != nil {
		return errors.New("Failed to create report settings" + err.Error())
	}

	return nil
}

func GenerateReportSummary(payload requests.GenerateReportSummaryRequest, userNumber string) error {

	//check if report settings exists
	reportSettings, err := reporsitory.GetReportSettingByCampusClassNumberAndStudyPeriodAndStatus(payload.CampusClassNumber, payload.StudyPeriodNumber, constants.StatusActive)

	if err != nil {
		return errors.New("Failed to fetch settings  for this class")
	}

	//check if report columns exists
	reportColumns, err := reporsitory.FindReportColumnsByReportSettingNumber(reportSettings.Number)
	if err != nil {
		return errors.New("failed to fetch columns  for this class")
	}

	configuredColumnsTotal := len(reportColumns)
	columnsTotal := int(reportSettings.Columns)

	if configuredColumnsTotal != columnsTotal {
		return errors.New(fmt.Sprintf("You have configured %d and the report settings expects %d", configuredColumnsTotal, columnsTotal))
	}

	//Check if all has assessment schedules
	for _, reportColumn := range reportColumns {

		checkAssessments, err := reporsitory.FindReportColumnAssessmentsByReportColumnNumber(reportColumn.Number)
		if err != nil {
			return errors.New("Checking column assessments error " + err.Error())
		}

		if len(checkAssessments) == 0 {
			return errors.New(fmt.Sprintf("Report column '%s' has no assessment schedules set", reportColumn.Name))
		}

	}

	enrollment, err := enrolmentRepository.GetStudentEnrolment(payload.StudentNumber, payload.StudyPeriodNumber, payload.CampusClassNumber, constants.StatusActive)
	if err != nil {
		return errors.New("Getting active student enrolment error: " + err.Error())
	}

	if enrollment.Number == "" {
		return errors.New("The selected student is not yet enrolled ")
	}

	tx := config.DB.Begin()

	var reportTemplate reportTemplateModel.ReportTemplate

	if reportTemplate, err = services.DetermineReportTemplateForClass(payload.CampusClassNumber); err != nil {
		return err
	}

	if reportTemplate.Number == "" {
		tx.Rollback()
		return errors.New("this class has no report template set up")
	}

	headTeacherComment, err := reporsitory.GetHeadTeacherCommentByReportSettingNumber(reportSettings.Number)

	if err != nil {
		tx.Rollback()
		return errors.New("failed to fetch head teacher comment")

	}

	reportSummary, err := reporsitory.GetReportSummaryByStudentNumberAndStudyPeriodNumberAndCampusClassNumber(payload.StudentNumber, payload.StudyPeriodNumber, payload.CampusClassNumber)
	if err != nil {
		return err
	}
	if reportSummary.Number != "" {

		// update the existing report summary with data from the payload
		reportSummary.ReportSettingNumber = reportSettings.Number
		reportSummary.Grade = ""
		reportSummary.Aggregate = -1
		reportSummary.Comment = ""
		reportSummary.UserNumber = userNumber
		reportSummary.StreamNumber = enrollment.StreamNumber
		reportSummary.ReportTemplateName = reportTemplate.Name
		reportSummary.HeadTeacherComment = headTeacherComment.Comment

		if err := reporsitory.SaveReportSummary(reportSummary, tx); err != nil {
			tx.Rollback()
			return errors.New("failed to update report summary" + err.Error())

		}
	} else {

		reportSummary = models.ReportSummary{
			Number:              uuid.NewV4().String(),
			ReportSettingNumber: reportSettings.Number,
			StudentNumber:       payload.StudentNumber,
			CampusClassNumber:   payload.CampusClassNumber,
			StudyPeriodNumber:   payload.StudyPeriodNumber,
			Grade:               "",
			Aggregate:           -1,
			Comment:             "",
			UserNumber:          userNumber,
			Code:                0,
			StreamNumber:        enrollment.StreamNumber,
			ReportTemplateName:  reportTemplate.Name,
			HeadTeacherComment:  headTeacherComment.Comment,
		}

		//if err := reporsitory.SaveReportSummary(reportSummary, tx); err != nil {
		//	tx.Rollback()
		//	return errors.New("failed to create report summary")
		//
		//}
	}

	campusClass, err := repositories.GetCampusClassByNumber(payload.CampusClassNumber)

	if err != nil {
		tx.Rollback()
		return errors.New("failed to fetch campus class")
	}

	assessmentAreas, err := assessmentAreaRepo.FindAssessmentAreasByCurriculumClassNumber(campusClass.CurriculumClassNumber, tx)

	if err != nil {
		tx.Rollback()
		return errors.New("failed to fetch assessment areas")

	}
	var finalAggregateScale aggregateSaceltModel.AggregateScale
	var totalAggregates float64
	for _, assessmentArea := range assessmentAreas {
		var scale aggregateSaceltModel.AggregateScale
		var grade gradesModels.Grade

		for _, reportColum := range reportColumns {

			rcas, err := reporsitory.FindReportColumnAssessmentByReportColumnNumberAndAssessmentAreaNumber(reportColum.Number, assessmentArea.Number)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					continue
				}
			}

			var score float64
			for _, rca := range rcas {
				var assessment assesmentModel.Assessment

				if err := tx.Where("assessment_metric_number = ? AND student_number = ? ", rca.AssessmentMetricNumber, payload.StudentNumber).First(&assessment).Error; err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						continue
					}
				}
				score += assessment.ActualScore

			}
			length := float64(len(rcas))

			var scoreAverage = score / length
			if length == 0 || score == 0 {
				scoreAverage = -1
			}

			// find the aggregate scale with from grading standard of the  assessment metric
			var percentageScore float64

			if length > 0 && score > 0 {
				percentageScore = (scoreAverage / 100) * reportColum.PercentageShare

				scale, err = GetProperAggregateScale(scoreAverage, campusClass.Number)
				if err != nil {
					tx.Rollback()
					return err
				}
				grade, err = repository2.FindGradeByNumber(scale.GradeNumber)
				if err != nil {
					tx.Rollback()
					return errors.New("failed to fetch grade for this scale")
				}
			}

			//check if the report score already exists with the same report summary number, assessment area number and report column number
			var reportScore models.ReportScore
			if err := tx.Where("report_summary_number = ? AND assessment_area_number = ? AND report_column_number = ?", reportSummary.Number, assessmentArea.Number, reportColum.Number).First(&reportScore).Error; err == nil {
				//update the report score
				reportScore.Score = commons.RoundFloat(percentageScore, 2)
				reportScore.AverageScore = commons.RoundFloat(scoreAverage, 2)
				reportScore.Aggregate = scale.Weight
				reportScore.Grade = grade.Name
				reportScore.Illustration = scale.Illustration
				reportScore.UserNumber = userNumber
				reportScore.Remark = scale.Description

				//kkk

				if score > 0 {
					if err := reporsitory.SaveReportScore(reportScore); err != nil {
						tx.Rollback()
						return errors.New("failed to update report score")
					}
				}

			} else {

				//create a new report score

				reportScore = models.ReportScore{
					Number:               uuid.NewV4().String(),
					ReportSummaryNumber:  reportSummary.Number,
					AssessmentAreaNumber: assessmentArea.Number,
					ReportColumnNumber:   reportColum.Number,
					Score:                commons.RoundFloat(percentageScore, 2),
					AverageScore:         commons.RoundFloat(scoreAverage, 2),
					Aggregate:            scale.Weight,
					Grade:                grade.Name,
					UserNumber:           userNumber,
					Remark:               scale.Description,
				}

				if score > 0 {
					if err := reporsitory.CreateReportScore(reportScore, tx); err != nil {
						tx.Rollback()
						return errors.New("failed to create report score")
					}

				}

				finalAggregateScale = scale

			}

		}

		var rs []models.ReportScore

		if err := tx.Where("report_summary_number = ? AND assessment_area_number = ?", reportSummary.Number, assessmentArea.Number).Find(&rs).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				continue
			}
		}
		// sum report scores

		var summation float64
		for _, r := range rs {
			summation += r.Score
			fmt.Println(fmt.Sprintf("score %.1f", r.Score))
		}

		fmt.Println(fmt.Sprintf("summation %.1f", summation))
		scale, err = GetProperAggregateScale(summation, campusClass.Number)
		if err != nil {
			tx.Rollback()
			return err
		}

		totalAggregates += scale.Weight

		grade, err := repository2.FindGradeByNumber(scale.GradeNumber)

		if err != nil {
			tx.Rollback()
			return errors.New("failed to fetch grade for this scale")
		}

		//check if the report score summary already exists with the same report summary number and assessment area number
		var reportScoreSummary models.ReportScoreSummary
		if err := tx.Where("report_summary_number = ? AND assessment_area_number = ?", reportSummary.Number, assessmentArea.Number).First(&reportScoreSummary).Error; err == nil {
			//update the report score summary
			reportScoreSummary.Aggregate = scale.Weight
			reportScoreSummary.Grade = grade.ShortName
			reportScoreSummary.Summation = commons.RoundFloat(summation, 0)
			reportScoreSummary.UserNumber = userNumber
			reportScoreSummary.Comment = scale.Description
			reportScoreSummary.Illustration = scale.Illustration

			if err := reporsitory.SaveReportScoreSummary(reportScoreSummary, tx); err != nil {
				tx.Rollback()
				return errors.New("failed to update report score summary")
			}

		} else {
			//create a new report score summary
			reportScoreSummary = models.ReportScoreSummary{
				Number:               uuid.NewV4().String(),
				ReportSummaryNumber:  reportSummary.Number,
				AssessmentAreaNumber: assessmentArea.Number,
				Aggregate:            scale.Weight,
				Grade:                grade.ShortName,
				Summation:            commons.RoundFloat(summation, 0),
				UserNumber:           userNumber,
				Comment:              scale.Description,
				Illustration:         scale.Illustration,
			}

			if err := reporsitory.CreateReportScoreSummary(reportScoreSummary, tx); err != nil {
				tx.Rollback()
				return errors.New("failed to create report score summary")
			}

		}

		finalAggregateScale = scale

	}
	assessmentGrade, err := GetReportSummaryGrade(totalAggregates, payload.CampusClassNumber)

	if err != nil {
		tx.Rollback()
		reportSummary.Grade = "N/A"
		return err

	} else {
		fmt.Println("assessmentGrade", assessmentGrade.Name)
		reportSummary.Grade = assessmentGrade.Name
		reportSummary.AssessmentGradeScaleNumber = assessmentGrade.Number
	}

	reportSummary.Aggregate = totalAggregates

	headTeacherComment, err = GetProperHeadTeacherCommentByReportSettingNumberAndPoints(reportSettings.Number, totalAggregates)
	if err != nil {
		tx.Rollback()
		return err

	}
	reportSummary.HeadTeacherComment = headTeacherComment.Comment
	reportSummary.AggregateScaleUsed = finalAggregateScale.Number

	if err := reporsitory.SaveReportSummary(reportSummary, tx); err != nil {
		tx.Rollback()
		return errors.New("failed to update edit summary" + err.Error())
	}

	reportSummaryNumber := reportSummary.Number

	finalStatus, err := checkStatus(campusClass.Number, reportSummary.StudyPeriodNumber, reportSummary.AssessmentGradeScaleNumber, reportSummary.Aggregate, reportSummary.Aggregate, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if finalStatus.Number != "" {
		fmt.Println("Final learner status: " + finalStatus.Type)
		//FindPromotionByReportSummaryNumber
		currentPromotion, err := repositories10.FindPromotionByReportSummaryNumber(reportSummaryNumber, tx)
		if err != nil {
			tx.Rollback()
			return errors.New("failed to check promotion error: " + err.Error())
		}

		promotion := models2.Promotion{
			Number:                 uuid.NewV4().String(),
			ReportSummaryNumber:    reportSummaryNumber,
			PromotionSettingNumber: finalStatus.Number,
			Status:                 finalStatus.Type,
			UserNumber:             finalStatus.UserNumber,
			PromotionType:          "Auto",
			CreatedAt:              time.Now(),
			UpdatedAt:              time.Now(),
		}

		if currentPromotion.Number != "" {
			promotion.Number = currentPromotion.Number
		}

		err = repositories10.SavePromotion(promotion, tx)
		if err != nil {
			tx.Rollback()
			return errors.New("failed to save promotion " + err.Error())
		}

	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
	}
	return nil

}

func GetProperAggregateScale(score float64, campusClassNumber string) (aggregateSaceltModel.AggregateScale, error) {
	var scale aggregateSaceltModel.AggregateScale

	var aggregateScales []aggregateSaceltModel.AggregateScale

	score = commons.RoundFloat(score, 0)

	class, err := repositories.GetClassByCampusClassNumber(campusClassNumber)
	if err != nil {
		return scale, errors.New("Getting class from campus error " + err.Error())
	}

	classGradingStandard, err := repositories2.GetDefaultGradingStandardClassScaleByClassNumber(class.Number)
	if err != nil {
		return scale, errors.New("Checking for default grading standard error " + err.Error())
	}

	if classGradingStandard.GradingStandardNumber == "" {
		return scale, errors.New(fmt.Sprintf("%s has no grading standard set", class.Name))
	}

	scaleNumber, err := reporsitory.GetScaleNumberByGradingStandardNumberAndClassNumber(classGradingStandard.GradingStandardNumber, class.Number)
	if err != nil {
		return scale, errors.New("failed to fetch grading standard number")
	}

	if scaleNumber == "" {
		gradingStandard, _ := repositories2.GetGradingStandardByNumber(classGradingStandard.GradingStandardNumber)
		return scale, errors.New(fmt.Sprintf("There is no scale configured for %s and class %s ", gradingStandard.Name, class.Name))
	}

	aggregateScales, err = repositories2.FindAggregateScaleByScaleNumber(scaleNumber)

	for _, val := range aggregateScales {
		if score >= val.MinValue && score <= val.MaxValue {
			scale = val

			return scale, nil
		}
	}

	return scale, errors.New(fmt.Sprintf("no matching aggregate scale found for %.1f and %s", score, campusClassNumber))
}

func GetReportSummaryGrade(points float64, campusClassNumber string) (assesmentGradeModel.AssessmentGradeScale, error) {
	var matchedAssessmentGradeScale assesmentGradeModel.AssessmentGradeScale

	class, err := repositories.GetClassByCampusClassNumber(campusClassNumber)
	if err != nil {
		return matchedAssessmentGradeScale, errors.New("getting class from campus class number error: " + err.Error())
	}

	if class.Number == "" {
		return matchedAssessmentGradeScale, errors.New("no class from campus class number")
	}

	classAssessmentGrade, err := repositories2.GetAssessmentGradeByClassNumberAndStatus(class.Number, constants.StatusActive)
	if err != nil {
		return matchedAssessmentGradeScale, errors.New("getting class assessment grade error: " + err.Error())
	}

	if classAssessmentGrade.Number == "" {
		return matchedAssessmentGradeScale, errors.New(class.Name + " does not have assessment grade set")
	}

	assessmentGrades, err := repositories2.FindAssessmentGradeScaleByGradeNumber(classAssessmentGrade.Number)
	if err != nil {
		return matchedAssessmentGradeScale, errors.New("getting assessments grade scale error: " + err.Error())
	}

	if len(assessmentGrades) == 0 {
		return matchedAssessmentGradeScale, errors.New("assessment grade scale not set up for " + class.Name)
	}

	for _, val := range assessmentGrades {
		if points >= val.MinimumPoints && points <= val.MaximumPoints {
			matchedAssessmentGradeScale = val

			return matchedAssessmentGradeScale, nil
		}
	}

	return matchedAssessmentGradeScale, errors.New(fmt.Sprintf("%s points dont macth any grading catery for %s", strconv.FormatFloat(points, 'f', -1, 64), class.Name))

}

func GetProperHeadTeacherCommentByReportSettingNumberAndPoints(reportSettingNumber string, points float64) (models.HeadTeacherComment, error) {
	var headTeacherComment models.HeadTeacherComment
	var headTeacherComments []models.HeadTeacherComment

	headTeacherComments, err := reporsitory.FindHeadTeacherCommentsByReportSettingNumber(reportSettingNumber)
	if err != nil {
		return headTeacherComment, errors.New("failed to fetch head teacher comments")
	}

	for _, val := range headTeacherComments {
		if points >= val.MinAggregate && points <= val.MaxAggregate {
			headTeacherComment = val

			return headTeacherComment, nil
		}
	}

	return headTeacherComment, errors.New(fmt.Sprintf("No matching head teacher comment found for '%.1f' aggregate points", points))
}

func FetchAllReportSettings(campusNumber, classNumber, campusClassNumber, academicYearNumber, studyPeriodNumber string) ([]responses.ReportSettingsResponse, error) {
	var reportSettingsResponse = make([]responses.ReportSettingsResponse, 0)
	all, err := reporsitory.FetchAllReportSettings(campusNumber, classNumber, campusClassNumber, academicYearNumber, studyPeriodNumber)
	if err != nil {
		return reportSettingsResponse, errors.New("Fetching report settings error: " + err.Error())
	}

	for _, reportSetting := range all {
		reportSettingsResponse = append(reportSettingsResponse, createReportSettingsResponse(reportSetting))
	}

	return reportSettingsResponse, nil
}

func createReportSettingsResponse(reportSetting models.ReportSetting) responses.ReportSettingsResponse {

	studyPeriod, _ := repositories3.GetStudyPeriodByNumber(reportSetting.StudyPeriodNumber)

	studyPeriodName, _ := repositories3.GetStudyPeriodNameByNumber(studyPeriod.StudyPeriodNameNumber)

	academicYear, _ := repositories4.GetAcademicYearByNumber(reportSetting.AcademicYearNumber)
	campusClass, _ := repositories.GetCampusClassByNumber(reportSetting.CampusClassNumber)
	class, _ := repositories.GetClassByCampusClassNumber(campusClass.Number)

	return responses.ReportSettingsResponse{
		Name:                academicYear.AcademicYearName + " - " + class.Name + " - " + studyPeriodName.Name,
		ReportSettingNumber: reportSetting.Number,
		Status:              reportSetting.Status,
		Columns:             reportSetting.Columns,
		Code:                reportSetting.Code,
		AcademicYear:        academicYear.AcademicYearName,
		StudyPeriod:         studyPeriodName.Name,
		Class:               class.Name,
		StudyPeriodNumber:   studyPeriod.Number,
		AcademicYearNumber:  academicYear.Number,
		ClassNumber:         class.Number,
		LevelNumber:         class.LevelNumber,
		CampusClassNumber:   reportSetting.CampusClassNumber,
		CampusLevelNumber:   campusClass.CampusLevelNumber,
	}
}

func DeleteReportColumn(reportColumnNumber string) error {
	reportColumn, err := reporsitory.GetReportColumnByNumber(reportColumnNumber)
	if err != nil {
		return errors.New("failed to fetch report column by number: " + err.Error())
	}

	if reportColumn.Number == "" {
		return errors.New("Invalid report column number ")
	}

	tx := config.DB.Begin()
	err = reporsitory.DeleteReportColumn(reportColumn, tx)
	if err != nil {
		return errors.New("failed to delete report column: " + err.Error())
	}

	// update report settings columns
	reportSettings, err := reporsitory.GetReportSettingByNumber(reportColumn.ReportSettingNumber, tx)
	if err != nil {
		return errors.New("failed to fetch report settings by number: " + err.Error())
	}

	reportSettings.Columns = reportSettings.Columns - 1
	err = reporsitory.SaveReportSettings(reportSettings, tx)
	if err != nil {
		return errors.New("failed to update report settings: " + err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return errors.New("failed to commit transaction: " + err.Error())
	}

	return nil
}

func GetReportSummary(reportSummaryNumber string) (responses.FullReportSummary, error) {
	reportSummary, err := reporsitory.FindReportSummaryByNumber(reportSummaryNumber)
	if err != nil {
		return responses.FullReportSummary{}, errors.New("failed to find report summary error: " + err.Error())
	}

	campusClass, err := repositories.GetCampusClassByNumber(reportSummary.CampusClassNumber)
	if err != nil {
		return responses.FullReportSummary{}, errors.New("failed to find campus classes error: " + err.Error())
	}

	curriculumClass, err := repositories.GetCurriculumClass(campusClass.CurriculumClassNumber)
	if err != nil {
		return responses.FullReportSummary{}, errors.New("failed to find curriculum classes error: " + err.Error())
	}

	curriculumLevel, err := repositories5.GetCurriculumLevelByNumber(curriculumClass.CurriculumLevelNumber)
	if err != nil {
		return responses.FullReportSummary{}, errors.New("failed to find curriculum level error: " + err.Error())
	}

	campus, err := repositories6.GetCampusByNumber(campusClass.CampusNumber)
	if err != nil {
		return responses.FullReportSummary{}, errors.New("failed to find campus error: " + err.Error())
	}

	school, err := repositories7.GetSchoolByNumber(campus.SchoolNumber)
	if err != nil {
		return responses.FullReportSummary{}, errors.New("failed to find school error: " + err.Error())
	}

	learningAreas, err := repositories8.FindLearningAreasByCurriculumNumber(curriculumLevel.CurriculumNumber)
	if err != nil {
		return responses.FullReportSummary{}, errors.New("failed to find learning areas error: " + err.Error())
	}

	reportSetting, err := reporsitory.GetReportSettingByNumber(reportSummary.ReportSettingNumber)
	if err != nil {
		return responses.FullReportSummary{}, errors.New("failed to fetch report settings by number: " + err.Error())
	}

	class, err := repositories.GetClassByNumber(campusClass.ClassNumber)
	if err != nil {
		return responses.FullReportSummary{}, errors.New("failed to find class error: " + err.Error())
	}

	levelName, err := repositories5.GetLevelByNumber(class.LevelNumber)
	if err != nil {
		return responses.FullReportSummary{}, errors.New("failed to find level error: " + err.Error())
	}

	var rcs []string

	reportColumns, err := reporsitory.FindReportColumnsByReportSettingsNumber(reportSetting.Number)

	for _, reportColumn := range reportColumns {
		rcs = append(rcs, reportColumn.Name)
	}

	var la []responses.LearningArea

	for _, learningArea := range learningAreas {
		assessments := make([]responses.Assessments, 0)

		assessmentAreas, err := assessmentAreaRepositories.FindAllAssessmentAreasBySubjectCluster(campusClass.CurriculumClassNumber, curriculumLevel.CurriculumNumber, learningArea.Number)
		if err != nil {
			return responses.FullReportSummary{}, errors.New("failed to find assessment areas error: " + err.Error() + " for " + learningArea.Name)
		}

		for _, assessmentArea := range assessmentAreas {

			reportScoreSummary, err := reporsitory.FindReportScoreSummaryByReportSummaryAndAssessmentArea(reportSummaryNumber, assessmentArea.Number)
			if err != nil {
				return responses.FullReportSummary{}, errors.New("failed to find report score summary error: " + err.Error())
			}

			reportScores, err := reporsitory.FindAllReportScore(reportSummaryNumber, assessmentArea.Number)

			if len(reportScores) < len(rcs) {
				//  get the difference and add dummy report scores to the report scores
				var diff = len(rcs) - len(reportScores)
				for i := 0; i < diff; i++ {
					reportScores = append(reportScores, models.ReportScore{
						Score: -1,
					})
				}

			}

			assessments = append(assessments, responses.Assessments{
				Name:               assessmentArea.Name,
				ReportScores:       reportScores,
				ReportScoreSummary: reportScoreSummary,
			})
		}

		la = append(la, responses.LearningArea{
			Name:        learningArea.Name,
			Assessments: assessments,
		})

	}

	scalesUsed, err := reporsitory.GetAppliedScaleByScaleNumber(reportSummary.AggregateScaleUsed)
	if err != nil {
		return responses.FullReportSummary{}, errors.New("failed to fetch scale error: " + err.Error())
	}

	scales := reporsitory.BuildAppliedScaleResponse(scalesUsed)

	//GetStudyPeriodNameByNumber

	studyPeriod, err := repositories3.GetStudyPeriodByNumber(reportSummary.StudyPeriodNumber)
	if err != nil {
		return responses.FullReportSummary{}, errors.New("failed to find study period error: " + err.Error())
	}

	studyPeriodName, _ := repositories3.GetStudyPeriodNameByNumber(studyPeriod.StudyPeriodNameNumber)

	student, err := repositories9.FindStudentByNumber(reportSummary.StudentNumber)
	if err != nil {
		return responses.FullReportSummary{}, errors.New("failed to find student error: " + err.Error())
	}

	stream, err := repositories.GetStreamByNumber(reportSummary.StreamNumber)
	if err != nil {
		return responses.FullReportSummary{}, errors.New("failed to find stream error: " + err.Error())
	}

	reportTemplate, err := services.DetermineReportTemplateForClass(campusClass.Number)
	if err != nil {
		return responses.FullReportSummary{}, errors.New("failed to find report template error: " + err.Error())
	}

	response := responses.FullReportSummary{
		Student: responses.LeanStudent{
			Number:      student.Number,
			Surname:     student.Surname,
			SecondName:  student.SecondName,
			OtherName:   student.OtherName,
			Gender:      student.Gender,
			DateOfBirth: student.DateOfBirth,
		},
		Grade:               reportSummary.Grade,
		LearningAreas:       la,
		Stream:              stream.Name,
		Class:               campusClass.ClassName,
		ClassAbbreviation:   campusClass.ClassAbbreviation,
		School:              services4.BuildLeanSchoolResponse(&school),
		Campus:              campus,
		StudyPeriod:         studyPeriodName.Name,
		ReportColumns:       rcs,
		ReportTemplateName:  reportSummary.ReportTemplateName,
		TeacherComment:      reportSummary.TeacherComment,
		HeadTeachersComment: reportSummary.HeadTeacherComment,
		AggregateScale:      scales,
		LevelName:           levelName.Name,
		NextTermBegins:      studyPeriod.NextTermBeginsOn.String(),
		ReportTemplate:      reportTemplate.Name,
	}

	promo, err := repositories10.FindPromotionByReportSummaryNumber(reportSummaryNumber)
	if err != nil {
		return responses.FullReportSummary{}, errors.New("failed to find promotion status: " + err.Error())
	}

	if promo.Number != "" {
		response.PromotionStatus = promo.Status
	}

	return response, nil
}

func checkStatus(campusClassNumber, studyPeriodNumber, assessmentGradeScaleNumber string, totalAggregate, totalScore float64, db *gorm.DB) (models2.PromotionSetting, error) {

	finalStatus, err := checkIfPromoted(campusClassNumber, studyPeriodNumber, assessmentGradeScaleNumber, db)
	if err != nil {
		return models2.PromotionSetting{}, err
	}

	if finalStatus.Number != "" {
		return finalStatus, nil
	}

	finalStatus, err = checkIfRepeat(campusClassNumber, studyPeriodNumber, assessmentGradeScaleNumber, db)
	if err != nil {
		return models2.PromotionSetting{}, err
	}

	if finalStatus.Number != "" {
		return finalStatus, nil
	}

	return models2.PromotionSetting{}, errors.New("Failed to find promotion/demotion/repeat configuration criteria, error: ")
}

func checkIfPromoted(campusClassNumber string, studyPeriodNumber string, assessmentGradeScaleNumber string, db *gorm.DB) (models2.PromotionSetting, error) {

	//Check promotion
	promotion, err := repositories10.FindPromotionSettingsByStudyPeriod(studyPeriodNumber, campusClassNumber, "promotion", db)
	if err != nil {
		return models2.PromotionSetting{}, errors.New("Failed to find promotion configuration, error: " + err.Error())
	}

	if promotion.Number == "" {
		return models2.PromotionSetting{}, nil
	}

	scoredAssessmentGrade, err := repositories2.FindAssessmentGradeScaleByNumber(assessmentGradeScaleNumber, db)
	if err != nil {
		return models2.PromotionSetting{}, errors.New("Failed to find scored assessment grade, error: " + err.Error())
	}

	promotionValue, err := repositories10.FindPromotionSettingValueBySettingNumber(promotion.Number, db)
	if err != nil {
		return models2.PromotionSetting{}, errors.New("Failed to find required aggregate scale error: " + err.Error())
	}

	if promotionValue.Number != "" {
		return models2.PromotionSetting{}, errors.New("Promotion value not set for " + promotion.Name)
	}

	minimumAssessmentGrade, err := repositories2.FindAssessmentGradeScaleByNumber(promotionValue.AssessmentGradeScaleNumber, db)
	if err != nil {
		return models2.PromotionSetting{}, errors.New("Failed to find assessment grade, error: " + err.Error())
	}

	fmt.Println("Minimum assessment grade: ", minimumAssessmentGrade.Name)

	if minimumAssessmentGrade.Code >= scoredAssessmentGrade.Code {
		return promotion, nil
	}

	return models2.PromotionSetting{}, nil
}

func checkIfRepeat(campusClassNumber string, studyPeriodNumber string, assessmentGradeScaleNumber string, db *gorm.DB) (models2.PromotionSetting, error) {

	//Check repeat
	promotion, err := repositories10.FindPromotionSettingsByStudyPeriod(studyPeriodNumber, campusClassNumber, "repeat", db)
	if err != nil {
		return models2.PromotionSetting{}, errors.New("Failed to find promotion configuration, error: " + err.Error())
	}

	if promotion.Number == "" {
		return models2.PromotionSetting{}, nil
	}

	scoredAssessmentGrade, err := repositories2.FindAssessmentGradeScaleByNumber(assessmentGradeScaleNumber, db)
	if err != nil {
		return models2.PromotionSetting{}, errors.New("Failed to find scored assessment grade, error: " + err.Error())
	}

	promotionValue, err := repositories10.FindPromotionSettingValueBySettingNumber(promotion.Number, db)
	if err != nil {
		return models2.PromotionSetting{}, errors.New("Failed to find required aggregate scale error: " + err.Error())
	}

	if promotionValue.Number != "" {
		return models2.PromotionSetting{}, errors.New("Promotion value not set for " + promotion.Name)
	}

	minimumAssessmentGrade, err := repositories2.FindAssessmentGradeScaleByNumber(promotionValue.AssessmentGradeScaleNumber, db)
	if err != nil {
		return models2.PromotionSetting{}, errors.New("Failed to find assessment grade, error: " + err.Error())
	}

	if minimumAssessmentGrade.Code >= scoredAssessmentGrade.Code {
		return promotion, nil
	}

	return models2.PromotionSetting{}, nil
}
