package models

import (
	"errors"
	"time"
)

// ErrCompanyNotFound is an error raised when a product can not be found in the database
var ErrCompanyNotFound = errors.New("company not found")

// Company defines the structure for an API company
type Company struct {
	ID        uint64    `json:"id" gorm:"primary_key" validate:"excluded_with_all"`
	Name      string    `json:"name" gorm:"type:varchar(255);index" validate:"required,gt=0,max=255"`
	Code      string    `json:"code" gorm:"type:varchar(255);index" validate:"required,gt=0,max=255"`
	Country   string    `json:"country" gorm:"type:varchar(2);index" validate:"required,iso3166_1_alpha2"`
	Website   string    `json:"website" gorm:"type:varchar(255);index" validate:"required,fqdn"`
	Phone     string    `json:"phone" gorm:"type:varchar(255);index" validate:"required,e164"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
