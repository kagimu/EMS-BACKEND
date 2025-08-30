package requests

import (
	"errors"
	"html"
	"strings"
)

type SubTopicInfo struct {
	Code                 string `json:"code"`
	TeachingSequenceCode string `json:"teaching_sequence_code"`
	Name                 string `json:"name"`
	TopicNumber          string `json:"topic_number"`
	EstimatedTime        string `json:"estimated_time"`
	CurriculumNumber     string `json:"curriculum_number"`
	SubTopicNumber       string `json:"sub_topic_number"`
}

type SubTopicDelRequest struct {
	Number string `json:"number"`
}

func (p *SubTopicInfo) Prepare() {
	p.Name = html.EscapeString(strings.TrimSpace(p.Name))
}

func (p *SubTopicInfo) Validate() error {

	if p.Code == "" {
		return errors.New("code is required")
	}

	if p.Code == "0" {
		return errors.New("code must be greater than zero")
	}

	if p.Name == "" {
		return errors.New("name is required")
	}

	if p.CurriculumNumber == "" {
		return errors.New("curriculum is required")
	}

	if p.TopicNumber == "" {
		return errors.New("topic is required")
	}

	if p.EstimatedTime == "" {
		return errors.New("estimated time is required")
	}

	if p.TeachingSequenceCode == "" {
		return errors.New("teaching sequence code is required")
	}

	return nil
}
