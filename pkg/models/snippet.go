package models

import (
	"time"
)

// ID - ID of table
type ID struct {
	ID uint `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
}

// DumbID - ID key but dump output in json
type DumbID struct {
	ID uint `gorm:"primary_key;AUTO_INCREMENT" json:"-"`
}

// Timestamp - database timestamp
type Timestamp struct {
	CreatedAt time.Time  `json:"create_at"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
}
