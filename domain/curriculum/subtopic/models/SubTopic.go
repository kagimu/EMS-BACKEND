package models

import (
	"github.com/hamsoftug/wholesomeEd/domain/auth/models"
	models2 "github.com/hamsoftug/wholesomeEd/domain/curriculum/topic/models"
	"time"
)

type (
	SubTopic struct {
		ID            uint          `gorm:"primaryKey"`
		Number        string        `json:"number" gorm:"size:100;unique;not null"`
		Name          string        `json:"name" form:"name" binding:"required" gorm:"size:60;unique;not null"`
		UserNumber    string        `json:"user_number" gorm:"size:40;not null"`
		EstimatedTime string        `json:"estimated_time" gorm:"size:10;not null"`
		TopicNumber   string        `json:"topic_number" gorm:"size:40;not null"`
		CreatedAt     time.Time     `json:"created_at"`
		UpdatedAt     time.Time     `json:"updated_at"`
		DeletedAt     *time.Time    `sql:"index"`
		Topic         models2.Topic `gorm:"foreignKey:TopicNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
		User          models.User   `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	}
)

func (u *SubTopic) TableName() string {
	return "sub_topics"
}
