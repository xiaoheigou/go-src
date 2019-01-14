package models

import (
	"time"
)

// ID - ID of table
type ID struct {
	Id int64 `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
}

// DumbID - ID key but dump output in json
type DumbID struct {
	ID uint `gorm:"primary_key;AUTO_INCREMENT" json:"-"`
}

// Timestamp - database timestamp
type Timestamp struct {
	CreatedAt time.Time  `gorm:"type:timestamp(6)" json:"created_at"`
	UpdatedAt time.Time  `gorm:"type:timestamp(6);null;" json:"updated_at"`
	DeletedAt *time.Time `gorm:"timestamp" json:"-"`
}
