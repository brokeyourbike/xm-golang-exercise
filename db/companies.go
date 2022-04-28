package db

import (
	"errors"

	"github.com/brokeyourbike/xm-golang-exercise/api/requests"
	"github.com/brokeyourbike/xm-golang-exercise/models"
	"gorm.io/gorm"
)

type CompaniesRepo struct {
	db *gorm.DB
}

func NewCompaniesRepo(db *gorm.DB) *CompaniesRepo {
	return &CompaniesRepo{db: db}
}

func (c *CompaniesRepo) Create(company models.Company) (models.Company, error) {
	err := c.db.Create(&company).Error
	return company, err
}

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

func (c *CompaniesRepo) GetAll(p requests.CompanyPayload) ([]models.Company, error) {
	var companies []models.Company

	err := c.db.Where(p.ToCompany()).Find(&companies).Error
	if err != nil {
		return companies, err
	}

	return companies, nil
}

func (c *CompaniesRepo) Delete(id uint64) error {
	return c.db.Delete(&models.Company{ID: id}).Error
}

func (c *CompaniesRepo) Update(company models.Company) error {
	return c.db.Model(&company).Updates(&models.Company{
		Name:    company.Name,
		Code:    company.Code,
		Country: company.Country,
		Website: company.Website,
		Phone:   company.Phone,
	}).Error
}
