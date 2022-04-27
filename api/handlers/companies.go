package handlers

import (
	"net/http"

	"github.com/brokeyourbike/xm-golang-exercise/models"
	"github.com/go-chi/render"
)

// CompanyCtxKey is a key used for the Company object in the context
type CompanyCtxKey struct{}

// CompanyPayloadCtxKey is a key used for the Company payload object in the context
type CompanyPayloadCtxKey struct{}

type CompaniesRepo interface {
	Create(company models.Company) (models.Company, error)
	Get(id uint64) (models.Company, error)
	Delete(id uint64) error
	Update(company models.Company) error
}

type companies struct {
	companiesRepo CompaniesRepo
}

func NewCompanies(c CompaniesRepo) *companies {
	return &companies{companiesRepo: c}
}

// HandleCompanyCreate handles POST requests to create companies
func (c *companies) HandleCompanyCreate(w http.ResponseWriter, r *http.Request) {
	data := r.Context().Value(CompanyPayloadCtxKey{}).(models.Company)

	company, err := c.companiesRepo.Create(data)
	if err != nil {
		render.Render(w, r, &ErrResponse{Message: "Company cannot be created", HTTPStatusCode: http.StatusInternalServerError})
		return
	}

	render.Render(w, r, &Payload{Company: &company, HTTPStatusCode: http.StatusCreated})
}

// HandleCompanyUpdate handles PUT requests to update companies
func (c *companies) HandleCompanyUpdate(w http.ResponseWriter, r *http.Request) {
	data := r.Context().Value(CompanyPayloadCtxKey{}).(models.Company)
	company := r.Context().Value(CompanyCtxKey{}).(models.Company)

	data.ID = company.ID

	if err := c.companiesRepo.Update(data); err != nil {
		return
	}

	render.Render(w, r, &Payload{Company: &data, HTTPStatusCode: http.StatusCreated})
}

// HandleCompanyDelete handles DELETE requests to delete companies
func (c *companies) HandleCompanyDelete(w http.ResponseWriter, r *http.Request) {
	company := r.Context().Value(CompanyCtxKey{}).(models.Company)

	if err := c.companiesRepo.Delete(company.ID); err != nil {
		render.Render(w, r, &ErrResponse{Message: "Cannot removed company", HTTPStatusCode: http.StatusInternalServerError})
		return
	}

	render.Render(w, r, &ErrResponse{Message: "Company removed", HTTPStatusCode: http.StatusOK})
}
