package requests

import "errors"

type ScopeUserCampusRequest struct {
	CampusNumber  string `json:"campus_number"`
	AddedByNumber string `json:"added_by_number"`
}

func (s *ScopeUserCampusRequest) ValidateCampus() error {
	if s.CampusNumber == "" {
		return errors.New("campus number is required")
	}
	if s.AddedByNumber == "" {
		return errors.New("added by number is required")
	}
	return nil
}

type ScopeUserClassRequest struct {
	ClassNumber   string `json:"class_number"`
	AddedByNumber string `json:"added_by_number"`
}

func (s *ScopeUserClassRequest) ValidateClass() error {
	if s.ClassNumber == "" {
		return errors.New("class number is required")
	}
	if s.AddedByNumber == "" {
		return errors.New("added by number is required")
	}
	return nil
}

type ScopeUserCurriculumRequest struct {
	CurriculumNumber string `json:"curriculum_number"`
	AddedByNumber    string `json:"added_by_number"`
}

func (s *ScopeUserCurriculumRequest) ValidateCurriculum() error {
	if s.CurriculumNumber == "" {
		return errors.New("curriculum number is required")
	}
	if s.AddedByNumber == "" {
		return errors.New("added by number is required")
	}
	return nil
}

type ScopeUserLearningAreaRequest struct {
	LearningAreaNumber string `json:"learning_area_number"`
	AddedByNumber      string `json:"added_by_number"`
}

func (s *ScopeUserLearningAreaRequest) ValidateLearningArea() error {
	if s.LearningAreaNumber == "" {
		return errors.New("learning area number is required")
	}
	if s.AddedByNumber == "" {
		return errors.New("added by number is required")
	}
	return nil
}

type ScopeUserLevelRequest struct {
	LevelNumber   string `json:"level_number"`
	AddedByNumber string `json:"added_by_number"`
}

func (s *ScopeUserLevelRequest) ValidateLevel() error {
	if s.LevelNumber == "" {
		return errors.New("level number is required")
	}
	if s.AddedByNumber == "" {
		return errors.New("added by number is required")
	}
	return nil
}

type ScopeUserSchoolRequest struct {
	SchoolNumber  string `json:"school_number"`
	AddedByNumber string `json:"added_by_number"`
}

func (s *ScopeUserSchoolRequest) ValidateSchool() error {
	if s.SchoolNumber == "" {
		return errors.New("school number is required")
	}
	if s.AddedByNumber == "" {
		return errors.New("added by number is required")
	}
	return nil
}
