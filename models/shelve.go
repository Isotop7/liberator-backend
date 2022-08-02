package models

import "github.com/jinzhu/gorm"

// Shelve
type Shelve struct {
	gorm.Model
	Location string `json:"location"`
	Content  []Book `json:"content"`
}
