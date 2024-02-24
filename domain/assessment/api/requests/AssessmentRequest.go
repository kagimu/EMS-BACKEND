package requests

import "errors"

type AssessmentRequest struct {
	AssessmentAreaNumber  string `json:"assessment_area_number"`
	Subject               string `json:"subject"`
	Attendance            uint64 `json:"attendance"`
	Observations          string `json:"observations"`
	GuidedExercise        uint64 `json:"guided_exercise"`
	UnguidedExercise      uint64 `json:"unguided_exercise"`
	ExplorativeAssignment uint64 `json:"explorative_assignment"`
	GroupWork             uint64 `json:"group_work"`
	TopicalText           uint64 `json:"topical_text"`
	Project               uint64 `json:"project"`
	CreativeExperiment    uint64 `json:"creative_experiment"`
	LevelNumber           string `json:"level_number"`
	ExaminationTypeNumber string `json:"examination_type_number"`
	StudentNumber         string `json:"student_number"`
}

func (a *AssessmentRequest) Validate() error {
	if a.Subject == "" {
		return errors.New("subject is required")
	}
	if a.LevelNumber == "" {
		return errors.New("level is required")
	}
	if a.AssessmentAreaNumber == "" {
		return errors.New("assessment area is required")
	}
	if a.ExaminationTypeNumber == "" {
		return errors.New("examination type is required")
	}
	if a.StudentNumber == "" {
		return errors.New("student number is required")
	}

	return nil
}
