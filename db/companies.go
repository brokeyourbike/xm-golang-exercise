package db

import (
	"errors"

	"github.com/brokeyourbike/xm-golang-exercise/api/requests"
	"github.com/brokeyourbike/xm-golang-exercise/models"
	"gorm.io/gorm"
)

// CompaniesRepo allows to store and retrieve companies from the database.
type CompaniesRepo struct {
	db *gorm.DB
}

// NewCompaniesRepo creates a new instance of the CompaniesRepo.
func NewCompaniesRepo(db *gorm.DB) *CompaniesRepo {
	return &CompaniesRepo{db: db}
}

// Create adds a new company to the database.
func (c *CompaniesRepo) Create(company models.Company) (models.Company, error) {
	err := c.db.Create(&company).Error
	return company, err
}

// Get returns a single company from the database.
// If a company with the given id does not exist in the database
// this function returns a ErrCompanyNotFound error.
func (c *CompaniesRepo) Get(id uint64) (models.Company, error) {
	var company models.Company

	err := c.db.Where("id = ?", id).First(&company).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return company, models.ErrCompanyNotFound
	}

	if err != nil {
		return company, err
	}

	return company, nil
}

// GetAll returns all companies from the database.
func (c *CompaniesRepo) GetAll(p requests.CompanyPayload) ([]models.Company, error) {
	var companies []models.Company

	err := c.db.Where(p.ToCompany()).Find(&companies).Error
	if err != nil {
		return companies, err
	}

	return companies, nil
}

// Delete deletes a company from the database.
func (c *CompaniesRepo) Delete(id uint64) error {
	return c.db.Delete(&models.Company{ID: id}).Error
}

// Update updates a company from the database.
func (c *CompaniesRepo) Update(company models.Company) error {
	return c.db.Model(&company).Updates(&models.Company{
		Name:    company.Name,
		Code:    company.Code,
		Country: company.Country,
		Website: company.Website,
		Phone:   company.Phone,
	}).Error
}
