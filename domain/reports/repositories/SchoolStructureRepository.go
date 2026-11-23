package repositories

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	campusModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	classModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reports/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reports/api/responses"
	schoolModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/models"
	levelModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/level/models"
	"github.com/gin-gonic/gin"
)

func FetchSchoolStructureReport(ctx *gin.Context, filters requests.ReportsFilterRequest) (responses.SchoolStructureReportResponse, error) {
	resp := responses.SchoolStructureReportResponse{
		Status:  http.StatusOK,
		Message: "School structure report fetched successfully",
		FilterOptions: responses.SchoolStructureFilterOptions{
			Schools:  []responses.SchoolFilterOption{},
			Campuses: []responses.CampusFilterOption{},
		},
		Data: []responses.SchoolReportSchoolNode{},
	}

	db := config.DB

	// 1. Resolve current entity number from user header / context
	entityNumber, _ := commons.GetUserEntityNumber(ctx)

	// 2. Fetch all entity Schools & Campuses first to populate FilterOptions
	var allEntitySchools []schoolModel.School
	baseSchoolQuery := db.Model(&schoolModel.School{})
	if entityNumber != "" {
		baseSchoolQuery = baseSchoolQuery.Where("entity_number = ?", entityNumber)
	}
	if err := baseSchoolQuery.Order("created_at asc").Find(&allEntitySchools).Error; err != nil {
		return resp, err
	}

	if len(allEntitySchools) == 0 {
		return resp, nil
	}

	allSchoolNumbers := make([]string, 0, len(allEntitySchools))
	for _, s := range allEntitySchools {
		allSchoolNumbers = append(allSchoolNumbers, s.Number)
		resp.FilterOptions.Schools = append(resp.FilterOptions.Schools, responses.SchoolFilterOption{
			Number: s.Number,
			Name:   s.Name,
		})
	}

	var allEntityCampuses []campusModel.Campus
	if err := db.Model(&campusModel.Campus{}).
		Where("school_number IN ?", allSchoolNumbers).
		Order("created_at asc").
		Find(&allEntityCampuses).Error; err != nil {
		return resp, err
	}

	for _, c := range allEntityCampuses {
		resp.FilterOptions.Campuses = append(resp.FilterOptions.Campuses, responses.CampusFilterOption{
			Number:       c.Number,
			SchoolNumber: c.SchoolNumber,
			Name:         c.Name,
		})
	}

	// 3. Apply optional SchoolNumber filter
	schools := make([]schoolModel.School, 0, len(allEntitySchools))
	for _, s := range allEntitySchools {
		if filters.SchoolNumber != "" && s.Number != filters.SchoolNumber {
			continue
		}
		schools = append(schools, s)
	}

	if len(schools) == 0 {
		return resp, nil
	}

	// 4. Apply optional CampusNumber filter
	campusNumbers := make([]string, 0, len(allEntityCampuses))
	campusesBySchool := make(map[string][]campusModel.Campus)
	for _, c := range allEntityCampuses {
		if filters.SchoolNumber != "" && c.SchoolNumber != filters.SchoolNumber {
			continue
		}
		if filters.CampusNumber != "" && c.Number != filters.CampusNumber {
			continue
		}
		campusNumbers = append(campusNumbers, c.Number)
		campusesBySchool[c.SchoolNumber] = append(campusesBySchool[c.SchoolNumber], c)
	}

	// 5. Fetch CampusCurriculums under the filtered campuses
	var campusCurriculums []campusModel.CampusCurriculum
	if len(campusNumbers) > 0 {
		ccQuery := db.Model(&campusModel.CampusCurriculum{}).
			Where("campus_number IN ?", campusNumbers).
			Preload("SchoolCurriculum.Curriculum").
			Preload("SchoolCurriculum")
		if err := ccQuery.Find(&campusCurriculums).Error; err != nil {
			return resp, err
		}
	}

	campusCurrNumbers := make([]string, 0, len(campusCurriculums))
	curriculumsByCampus := make(map[string][]campusModel.CampusCurriculum)
	for _, cc := range campusCurriculums {
		if filters.Curriculum != "" {
			currName := strings.ToLower(cc.SchoolCurriculum.Curriculum.Name)
			filterCurr := strings.ToLower(filters.Curriculum)
			if cc.Number != filters.Curriculum &&
				cc.SchoolCurriculum.CurriculumNumber != filters.Curriculum &&
				!strings.Contains(currName, filterCurr) {
				continue
			}
		}
		campusCurrNumbers = append(campusCurrNumbers, cc.Number)
		curriculumsByCampus[cc.CampusNumber] = append(curriculumsByCampus[cc.CampusNumber], cc)
	}

	// 6. Fetch CampusLevels to map campus_level_number -> campus_curriculum_number & level name
	var campusLevels []levelModel.CampusLevel
	if len(campusCurrNumbers) > 0 {
		db.Model(&levelModel.CampusLevel{}).
			Where("campus_curriculum_number IN ?", campusCurrNumbers).
			Preload("CurriculumLevel.Level").
			Preload("CurriculumLevel").
			Order("code asc").
			Find(&campusLevels)
	}

	campusLevelToCampusCurr := make(map[string]string)
	levelNamesByCampusCurr := make(map[string][]string)
	seenLevelByCampusCurr := make(map[string]map[string]bool)

	for _, cl := range campusLevels {
		campusLevelToCampusCurr[cl.Number] = cl.CampusCurriculumNumber
		lvlName := cl.CurriculumLevel.Level.Name
		if lvlName != "" {
			if seenLevelByCampusCurr[cl.CampusCurriculumNumber] == nil {
				seenLevelByCampusCurr[cl.CampusCurriculumNumber] = make(map[string]bool)
			}
			if !seenLevelByCampusCurr[cl.CampusCurriculumNumber][lvlName] {
				seenLevelByCampusCurr[cl.CampusCurriculumNumber][lvlName] = true
				levelNamesByCampusCurr[cl.CampusCurriculumNumber] = append(levelNamesByCampusCurr[cl.CampusCurriculumNumber], lvlName)
			}
		}
	}

	// 7. Fetch CampusClasses under these campuses
	var campusClasses []classModel.CampusClass
	if len(campusNumbers) > 0 {
		db.Model(&classModel.CampusClass{}).
			Where("campus_number IN ?", campusNumbers).
			Preload("CurriculumClass.Class").
			Preload("CurriculumClass.CurriculumLevel.Level").
			Preload("CurriculumClass").
			Order("code asc").
			Find(&campusClasses)
	}

	// Group classes by campus_curriculum_number (or fallback by campus_number + curriculum_number)
	classesByCampusCurr := make(map[string][]responses.SchoolReportClassNode)
	classesByCampusAndCurrNum := make(map[string][]responses.SchoolReportClassNode)

	for _, cls := range campusClasses {
		className := cls.CurriculumClass.Class.Name
		if className == "" {
			className = cls.ClassName
		}
		abbr := cls.CurriculumClass.Class.Abbreviation
		if abbr == "" {
			abbr = cls.ClassAbbreviation
		}
		if abbr == "" && cls.Code > 0 {
			abbr = fmt.Sprintf("%d", cls.Code)
		}
		lvlName := cls.CurriculumClass.CurriculumLevel.Level.Name

		node := responses.SchoolReportClassNode{
			ID:           cls.Number,
			Number:       cls.Number,
			ClassNumber:  cls.CurriculumClass.Class.Number,
			Code:         abbr,
			Name:         className,
			Abbreviation: abbr,
			Type:         "Class",
			LevelName:    lvlName,
			Status:       "Active",
		}

		if ccNum, ok := campusLevelToCampusCurr[cls.CampusLevelNumber]; ok && ccNum != "" {
			classesByCampusCurr[ccNum] = append(classesByCampusCurr[ccNum], node)
		} else if cls.CurriculumClass.CurriculumLevel.CurriculumNumber != "" {
			key := cls.CampusNumber + ":" + cls.CurriculumClass.CurriculumLevel.CurriculumNumber
			classesByCampusAndCurrNum[key] = append(classesByCampusAndCurrNum[key], node)
		}
	}

	// 8. Assemble hierarchical tree and rollup counts
	totalCampuses := 0
	totalCurricula := 0
	totalClasses := 0

	for idx, sch := range schools {
		schCampuses := campusesBySchool[sch.Number]
		// If a specific campus filter was requested and this school doesn't own that campus, skip this school
		if filters.CampusNumber != "" && len(schCampuses) == 0 {
			continue
		}

		campusNodes := make([]responses.SchoolReportCampusNode, 0, len(schCampuses))
		schCurriculaCount := 0
		schClassesCount := 0
		defaultSchoolLocation := ""

		for cIdx, cmp := range schCampuses {
			if defaultSchoolLocation == "" && cmp.Location != "" {
				defaultSchoolLocation = cmp.Location
			}

			cmpCurrs := curriculumsByCampus[cmp.Number]
			currNodes := make([]responses.SchoolReportCurriculumNode, 0, len(cmpCurrs))
			cmpClassesCount := 0

			for curIdx, cc := range cmpCurrs {
				clsList := classesByCampusCurr[cc.Number]
				if len(clsList) == 0 {
					fallbackKey := cmp.Number + ":" + cc.SchoolCurriculum.CurriculumNumber
					clsList = classesByCampusAndCurrNum[fallbackKey]
				}
				if clsList == nil {
					clsList = []responses.SchoolReportClassNode{}
				}

				lvlSummary := strings.Join(levelNamesByCampusCurr[cc.Number], ", ")
				if lvlSummary == "" && len(clsList) > 0 {
					lvlSummary = clsList[0].LevelName
				}
				if lvlSummary == "" {
					lvlSummary = "General"
				}

				currCode := fmt.Sprintf("CUR-%02d", curIdx+1)
				if cc.SchoolCurriculum.Curriculum.Code > 0 {
					currCode = fmt.Sprintf("%d", cc.SchoolCurriculum.Curriculum.Code)
				}

				currNodes = append(currNodes, responses.SchoolReportCurriculumNode{
					ID:               cc.Number,
					Number:           cc.Number,
					CurriculumNumber: cc.SchoolCurriculum.CurriculumNumber,
					Code:             currCode,
					Name:             cc.SchoolCurriculum.Curriculum.Name,
					Type:             "Curriculum",
					Level:            lvlSummary,
					Status:           "Active",
					ClassesCount:     len(clsList),
					Classes:          clsList,
				})

				cmpClassesCount += len(clsList)
			}

			campusNodes = append(campusNodes, responses.SchoolReportCampusNode{
				ID:             cmp.Number,
				Number:         cmp.Number,
				Code:           fmt.Sprintf("CMP-%02d", cIdx+1),
				Name:           cmp.Name,
				ShortName:      cmp.Name,
				Type:           "Campus",
				Location:       cmp.Location,
				Phone:          cmp.Phone,
				Email:          cmp.Email,
				Status:         "Active",
				CurriculaCount: len(currNodes),
				ClassesCount:   cmpClassesCount,
				Curricula:      currNodes,
			})

			schCurriculaCount += len(currNodes)
			schClassesCount += cmpClassesCount
		}

		headLabel := sch.PrincipalName
		if sch.PrincipalTitle != "" && sch.PrincipalName != "" {
			headLabel = fmt.Sprintf("%s (%s)", sch.PrincipalName, sch.PrincipalTitle)
		}

		resp.Data = append(resp.Data, responses.SchoolReportSchoolNode{
			ID:             sch.Number,
			Number:         sch.Number,
			Code:           fmt.Sprintf("SCH-%02d", idx+1),
			Name:           sch.Name,
			Type:           "School",
			Location:       defaultSchoolLocation,
			Head:           headLabel,
			Motto:          sch.Motto,
			Status:         "Active",
			CampusesCount:  len(campusNodes),
			CurriculaCount: schCurriculaCount,
			ClassesCount:   schClassesCount,
			Campuses:       campusNodes,
		})

		totalCampuses += len(campusNodes)
		totalCurricula += schCurriculaCount
		totalClasses += schClassesCount
	}

	resp.Summary = responses.SchoolStructureSummary{
		SchoolsCount:   len(resp.Data),
		CampusesCount:  totalCampuses,
		CurriculaCount: totalCurricula,
		ClassesCount:   totalClasses,
	}

	return resp, nil
}
