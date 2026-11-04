package requests

type EducationItemBasic struct {
	Institution string `json:"institution"`
	StartYear   string `json:"startYear"`
	EndYear     string `json:"endYear"`
	Award       string `json:"award"`
}

type EducationItemAdvanced struct {
	Institution string `json:"institution"`
	StartDate   string `json:"startDate"`
	EndDate     string `json:"endDate"`
	Award       string `json:"award"`
	AwardDate   string `json:"awardDate"`
}

type StaffBioDataEducationRequest struct {
	PrimaryEducation        []EducationItemBasic    `json:"primaryEducation"`
	SecondaryEducation      []EducationItemBasic    `json:"secondaryEducation"`
	HighSchoolEducation     []EducationItemBasic    `json:"highSchoolEducation"`
	PostHighSchoolEducation []EducationItemAdvanced `json:"postHighSchoolEducation"`
}

func (r *StaffBioDataEducationRequest) Validate() error {
	return nil
}
