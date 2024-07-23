package services

import (
	"errors"
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	aggregateSaceltModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/aggregateScale/models"
	assesmentGradeModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/assessmentGrade/models"
	gradesModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/grades/models"
	repository2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/grades/repository"
	assesmentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/models"
	repositories2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	assessmentAreaModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/models"
	assessmentAreaRepo "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/enrollment/repository"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/reporsitory"
	reportTemplateModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportTemplate/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportTemplate/services"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func SaveReportSettings(reportSettingsRequest requests.ReportSettingsRequest, currentUserNumber string) error {

	existingReportSettings, err := reporsitory.FindDuplicateReportSettings(reportSettingsRequest.StudyPeriodNumber, reportSettingsRequest.AcademicYearNumber, reportSettingsRequest.CampusClassNumber, constants.StatusActive)

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
		CampusClassNumber:  reportSettingsRequest.CampusClassNumber,
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

func GenerateReportSummary(payload requests.GenerateReportSummaryRequest, userNumber string) error {

	//check if report settings exists
	reportSettings, err := reporsitory.GetReportSettingByCampusClassNumberAndStatus(payload.CampusClassNumber, constants.StatusActive)

	if err != nil {
		return errors.New("Failed to fetch settings  for this class")
	}

	//check if report columns exists
	reportColumns, err := reporsitory.FindReportColumnsByReportSettingNumber(reportSettings.Number)
	if err != nil {
		return errors.New("failed to fetch columns  for this class")
	}

	if len(reportColumns) != int(reportSettings.Columns) {
		return errors.New("report columns set for the settings are not yet fully set up")

	}

	tx := config.DB.Begin()
	enrollment, err := repository.GetEnrollmentByStudentNumber(payload.StudentNumber, tx)

	if err != nil {
		return errors.New("failed to fetch enrollment for this student to pick stream")

	}

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

				scale, err = GetProperAggregateScale(scoreAverage, assessmentArea, campusClass.Number)
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
				}

				if score > 0 {
					if err := reporsitory.CreateReportScore(reportScore, tx); err != nil {
						tx.Rollback()
						return errors.New("failed to create report score")
					}

				}

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

		scale, err = GetProperAggregateScale(summation, assessmentArea, campusClass.Number)
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
			reportSummary.Grade = grade.Name
			reportScoreSummary.Summation = summation
			reportScoreSummary.UserNumber = userNumber

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
				Grade:                grade.Name,
				Summation:            summation,
				UserNumber:           userNumber,
			}

			if err := reporsitory.CreateReportScoreSummary(reportScoreSummary, tx); err != nil {
				tx.Rollback()
				return errors.New("failed to create report score summary")
			}

		}

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

	if err := reporsitory.SaveReportSummary(reportSummary, tx); err != nil {
		tx.Rollback()
		return errors.New("failed to update edit summary" + err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
	}
	return nil

}

func GetProperAggregateScale(score float64, assessmentArea assessmentAreaModels.AssessmentArea, campusClassNumber string) (aggregateSaceltModel.AggregateScale, error) {
	var scale aggregateSaceltModel.AggregateScale
	var aggregateScales []aggregateSaceltModel.AggregateScale

	class, err := repositories.GetClassByCampusClassNumber(campusClassNumber)
	gradingStandard, err := reporsitory.GetGradingStandardNumberByAssessmentAreaNumber(assessmentArea.Number)
	if err != nil {
		return scale, errors.New("failed to fetch grading standard by assessment area number")

	}

	if gradingStandard.Number == "" {
		return scale, errors.New(fmt.Sprintf("no grading standard for this assessment area (%s)", assessmentArea.Name))

	}

	scaleNumber, err := reporsitory.GetScaleNumberByGradingStandardNumberAndClassNumber(gradingStandard.Number, class.Number)
	if err != nil {
		return scale, errors.New("failed to fetch grading standard number")
	}

	if scaleNumber == "" {
		return scale, errors.New(fmt.Sprintf("There is no scale configured for %s and class %s ", gradingStandard.Name, class.Name))
	}

	aggregateScales, err = repositories2.FindAggregateScaleByScaleNumber(scaleNumber)

	//TODO get proper aggregate scale via  grading standard and class scales then scale then .. var aags assessmentAreaModel.AssessmentAreaGradingStandard
	for _, val := range aggregateScales {
		if score >= val.MinValue && score <= val.MaxValue {
			scale = val

			return scale, nil
		}
	}

	return scale, errors.New("no matching aggregate scale found")
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

	return matchedAssessmentGradeScale, errors.New("no matching points")

}
