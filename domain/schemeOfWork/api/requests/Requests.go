package requests

import (
	"errors"
	"fmt"
)

type ExcelSchemeOfWorkInfo struct {
	SubTopic      string
	SortingCode   string
	EstimatedTime string
}

type SubtopicListInfo struct {
	Number         string `json:"number"`
	Name           string `json:"name"`
	Code           int64  `json:"sort_code"`
	UserNumber     string `json:"user_number"`
	EstimatedTime  int64  `json:"estimated_time"`
	SubTopicNumber string `json:"sub_topic_number"`
}

type SchemeOfWorkInfo struct {
	Name               string             `json:"name"`
	ClassNumber        string             `json:"class_number"`
	SubjectNumber      string             `json:"subject_number"`
	AcademicYearNumber string             `json:"academic_year_number"`
	SchemeNumber       string             `json:"scheme_number"`
	SubtopicList       []SubtopicListInfo `json:"subtopicList"`
}

func (i SchemeOfWorkInfo) Validate() error {

	if i.Name == "" {
		return errors.New("name is required")
	}

	if i.SchemeNumber == "" {
		return errors.New("scheme number is required")
	}

	if i.ClassNumber == "" {
		return errors.New("class number is required")
	}

	if i.SubjectNumber == "" {
		return errors.New("subject number is required")
	}

	items := i.SubtopicList
	if len(items) == 0 {
		return errors.New("choose subtopics list")
	}

	isDup, code := hasDuplicateCodes(items)
	if isDup {
		return errors.New(fmt.Sprintf("You have a duplicate of '%d' in your code", code))
	}

	return nil
}

func hasDuplicateCodes(subtopics []SubtopicListInfo) (bool, int64) {
	seen := make(map[int64]bool)
	for _, s := range subtopics {
		if seen[s.Code] {
			return true, s.Code
		}
		seen[s.Code] = true
	}
	return false, 0
}
