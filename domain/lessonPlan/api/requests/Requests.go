package requests

import "errors"

type ClassInfo struct {
	Name         string `json:"name"`
	Abbreviation string `json:"abbreviation"`
	Level        string `json:"level"`
}

type StreamRequest struct {
	Name             string `json:"name"`
	ClassLevelNumber string `json:"class_level_number"`
}

type LessonPlanReferenceRequest struct {
	LessonPlanNumber string     `json:"lesson_plan_number"`
	References       []Material `json:"references"`
}

type Material struct {
	LearningMaterialNumber string `json:"learning_material_number"`
	Name                   string `json:"name"`
	Code                   string `json:"code"`
}

func (lpr *LessonPlanReferenceRequest) Validate() error {
	if len(lpr.References) < 1 {
		return errors.New("please select at-least one material")
	}

	/*if lpr.Code == "" {
		return errors.New("please enter the code")
	}

	if lpr.LearningMaterialNumber == "" {
		return errors.New("please select a learning material to continue")
	}

	if lpr.Name == "" {
		return errors.New("reference name is required")
	}
	*/
	if lpr.LessonPlanNumber == "" {
		return errors.New("please select a lesson plan to add a reference")
	}

	return nil
}
