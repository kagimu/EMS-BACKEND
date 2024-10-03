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
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/repositories"
	classService "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/services"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	assessmentAreaRepo "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/enrollment/repository"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/reporsitory"
	reportTemplateModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportTemplate/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportTemplate/services"
	repositories3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/repositories"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"strconv"
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
		return errors.New("please delete a column you dont need adn adjust the percentages instead of trying to reduce the colum size")
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
		}

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

	return scale, errors.New(fmt.Sprintf("no matching aggregate scale found for %d", scale))
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

	err = reporsitory.DeleteReportColumn(reportColumn)
	if err != nil {
		return errors.New("failed to delete report column: " + err.Error())
	}

	return nil
}
