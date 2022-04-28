package models

import (
	"errors"
)

// ErrCompanyNotFound is an error raised when a product can not be found in the database
var ErrCompanyNotFound = errors.New("company not found")

// Company defines the structure for an API company
type Company struct {
	ID      uint64 `json:"id" gorm:"primary_key"`
	Name    string `json:"name" gorm:"type:varchar(255);index"`
	Code    string `json:"code" gorm:"type:varchar(255);index"`
	Country string `json:"country" gorm:"type:varchar(2);index"`
	Website string `json:"website" gorm:"type:varchar(255);index"`
	Phone   string `json:"phone" gorm:"type:varchar(255);index"`
}
