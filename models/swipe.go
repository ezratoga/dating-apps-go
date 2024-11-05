package models

import (
	"github.com/lib/pq"
)

// Swipe models
type Swipe struct {
	ID				uint				`gorm:"unique:not null"`
	Username       	string				`gorm:"unique;not null"`
	Like			pq.StringArray		`gorm:"type:text[]"`
	Pass			pq.StringArray		`gorm:"type:text[]"`
	Match			pq.StringArray		`gorm:"type:text[]"`
}

// Swipe payload
type SwipePayload struct {
	TargetId		uint				`json:"target_id"	validate:"required"`
	Direction		string				`json:"direction"	validate:"oneof=left right"`
}