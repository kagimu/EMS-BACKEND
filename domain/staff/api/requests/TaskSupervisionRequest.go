package requests

import "errors"

type TaskSupervisionRequest struct {
	DateOfSupervision   string `json:"date_of_supervision"`
	DescriptiveAnalysis string `json:"descriptive_analysis"`
	NameOfSupervisor    string `json:"name_of_supervisor"`
	TaskNumber          string `json:"task_number"`
}

func (r *TaskSupervisionRequest) Validate() error {
	if r.DateOfSupervision == "" {
		return errors.New("date of supervision")
	}
	if r.DescriptiveAnalysis == "" {
		return errors.New("descriptive analysis")
	}
	if r.NameOfSupervisor == "" {
		return errors.New("name of supervisor")
	}

	return nil
}
