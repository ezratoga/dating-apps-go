package models

import (
	"time"
)

type Invoice struct {
	ID			uint		`gorm:"primaryKey"`
	InvoiceID	string		`gorm:"unique;type:text;not null"	json:"invoiceId"`
	UserID		uint		`gorm:"unique;not null"	json:"userId"`
	Status 		string 		`gorm:"type:text"	json:"status"`
	CreatedAt	time.Time
	UpdatedAt	time.Time
}

type InvoicePayload struct {
	InvoiceID	string		`gorm:"unique;type:text;not null"	json:"invoiceId"`
}