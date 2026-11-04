package responses

type SchoolStructureSummary struct {
	SchoolsCount   int `json:"schools_count"`
	CampusesCount  int `json:"campuses_count"`
	CurriculaCount int `json:"curricula_count"`
	ClassesCount   int `json:"classes_count"`
}

type SchoolFilterOption struct {
	Number string `json:"number"`
	Name   string `json:"name"`
}

type CampusFilterOption struct {
	Number       string `json:"number"`
	SchoolNumber string `json:"school_number"`
	Name         string `json:"name"`
}

type SchoolStructureFilterOptions struct {
	Schools  []SchoolFilterOption `json:"schools"`
	Campuses []CampusFilterOption `json:"campuses"`
}

type SchoolReportClassNode struct {
	ID           string `json:"id"`
	Number       string `json:"number"`
	ClassNumber  string `json:"class_number"`
	Code         string `json:"code"`
	Name         string `json:"name"`
	Abbreviation string `json:"abbreviation"`
	Type         string `json:"type"`
	LevelName    string `json:"level_name"`
	Status       string `json:"status"`
}

type SchoolReportCurriculumNode struct {
	ID               string                  `json:"id"`
	Number           string                  `json:"number"`
	CurriculumNumber string                  `json:"curriculum_number"`
	Code             string                  `json:"code"`
	Name             string                  `json:"name"`
	Type             string                  `json:"type"`
	Level            string                  `json:"level"`
	Status           string                  `json:"status"`
	ClassesCount     int                     `json:"classes_count"`
	Classes          []SchoolReportClassNode `json:"classes"`
}

type SchoolReportCampusNode struct {
	ID             string                       `json:"id"`
	Number         string                       `json:"number"`
	Code           string                       `json:"code"`
	Name           string                       `json:"name"`
	ShortName      string                       `json:"short_name"`
	Type           string                       `json:"type"`
	Location       string                       `json:"location"`
	Phone          string                       `json:"phone"`
	Email          string                       `json:"email"`
	Status         string                       `json:"status"`
	CurriculaCount int                          `json:"curricula_count"`
	ClassesCount   int                          `json:"classes_count"`
	Curricula      []SchoolReportCurriculumNode `json:"curricula"`
}

type SchoolReportSchoolCurriculumItem struct {
	Number           string `json:"number"`
	CurriculumNumber string `json:"curriculum_number"`
	Name             string `json:"name"`
}

type SchoolReportSchoolNode struct {
	ID               string                             `json:"id"`
	Number           string                             `json:"number"`
	Code             string                             `json:"code"`
	Name             string                             `json:"name"`
	Type             string                             `json:"type"`
	Location         string                             `json:"location"`
	Head             string                             `json:"head"`
	Motto            string                             `json:"motto"`
	Status           string                             `json:"status"`
	CampusesCount    int                                `json:"campuses_count"`
	CurriculaCount   int                                `json:"curricula_count"`
	ClassesCount     int                                `json:"classes_count"`
	SchoolCurriculum []SchoolReportSchoolCurriculumItem `json:"school_curriculum"`
	Campuses         []SchoolReportCampusNode           `json:"campuses"`
}

type SchoolStructureReportResponse struct {
	Status        int                          `json:"status"`
	Message       string                       `json:"message"`
	FilterOptions SchoolStructureFilterOptions `json:"filter_options"`
	Summary       SchoolStructureSummary       `json:"summary"`
	Data          []SchoolReportSchoolNode     `json:"data"`
}

