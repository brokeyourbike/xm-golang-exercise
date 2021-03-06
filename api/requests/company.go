package requests

import "github.com/brokeyourbike/xm-golang-exercise/models"

// CompanyPayload is a data structure used to decode incoming company data.
type CompanyPayload struct {
	Name    string `json:"name" shema:"name" validate:"required,gt=0,max=255"`
	Code    string `json:"code" shema:"code" validate:"required,gt=0,max=255"`
	Country string `json:"country" shema:"country" validate:"required,iso3166_1_alpha2"`
	Website string `json:"website" shema:"website" validate:"required,fqdn"`
	Phone   string `json:"phone" shema:"phone" validate:"required,e164"`
}

// ToCompany creates a new Company from the CompanyPayload.
func (c *CompanyPayload) ToCompany() models.Company {
	return models.Company{
		Name:    c.Name,
		Code:    c.Code,
		Country: c.Country,
		Website: c.Website,
		Phone:   c.Phone,
	}
}
