package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type OnlineTest struct {
	gorm.Model
	TestID    string     `gorm:"size:255;not null;"`
	Duration  int        `gorm:"not null;"`
	Questions []Question `gorm:"foreignKey:OnlineTestID"`
}

type Question struct {
	gorm.Model
	Type   string `gorm:"size:50;not null;"`
	Text   string `gorm:"not null;type:text;"`
	Answer string
	Meta   datatypes.JSON
	Order  int `gorm:"not null;"`
}

// Multiple-choice question metadata
type MCQMeta struct {
	Options []string `json:"options"`
}

// True/False metadata
type TOFMeta struct {
	Statements []string `json:"statements"`
}
