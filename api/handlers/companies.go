package handlers

import (
	"net/http"

	"github.com/brokeyourbike/xm-golang-exercise/api/requests"
	"github.com/brokeyourbike/xm-golang-exercise/api/responses"
	"github.com/brokeyourbike/xm-golang-exercise/models"
	"github.com/go-chi/render"
)

// CompanyCtxKey is a key used for the Company object in the context
type CompanyCtxKey struct{}

// CompanyPayloadCtxKey is a key used for the Company payload object in the context
type CompanyPayloadCtxKey struct{}

type CompaniesRepo interface {
	Create(models.Company) (models.Company, error)
	GetAll(requests.CompanyPayload) ([]models.Company, error)
	Get(uint64) (models.Company, error)
	Delete(uint64) error
	Update(models.Company) error
}

type companies struct {
	companiesRepo CompaniesRepo
}

func NewCompanies(c CompaniesRepo) *companies {
	return &companies{companiesRepo: c}
}

// HandleCompanyCreate handles POST requests to create companies
func (c *companies) HandleCompanyCreate(w http.ResponseWriter, r *http.Request) {
	data := r.Context().Value(CompanyPayloadCtxKey{}).(requests.CompanyPayload)

	company, err := c.companiesRepo.Create(data.ToCompany())
	if err != nil {
		render.Render(w, r, &responses.ErrResponse{Message: "Company cannot be created", HTTPStatusCode: http.StatusInternalServerError})
		return
	}

	render.Render(w, r, &responses.CompanyResponse{Company: &company, HTTPStatusCode: http.StatusCreated})
}

// HandleCompanyGetOne handles GET requests to display single company
func (c *companies) HandleCompanyGetOne(w http.ResponseWriter, r *http.Request) {
	company := r.Context().Value(CompanyCtxKey{}).(models.Company)

	render.Render(w, r, &responses.CompanyResponse{Company: &company, HTTPStatusCode: http.StatusOK})
}

// HandleCompanyGetAll handles GET requests to view companies
func (c *companies) HandleCompanyGetAll(w http.ResponseWriter, r *http.Request) {
	data := r.Context().Value(CompanyPayloadCtxKey{}).(requests.CompanyPayload)

	companies, err := c.companiesRepo.GetAll(data)
	if err != nil {
		render.Render(w, r, &responses.ErrResponse{Message: "Cannot retrieve companies", HTTPStatusCode: http.StatusInternalServerError})
		return
	}

	render.Render(w, r, &responses.CompaniesResponse{Companies: companies, HTTPStatusCode: http.StatusOK})
}

// HandleCompanyUpdate handles PUT requests to update companies
func (c *companies) HandleCompanyUpdate(w http.ResponseWriter, r *http.Request) {
	data := r.Context().Value(CompanyPayloadCtxKey{}).(requests.CompanyPayload)
	company := r.Context().Value(CompanyCtxKey{}).(models.Company)

	new := data.ToCompany()
	new.ID = company.ID

	if err := c.companiesRepo.Update(new); err != nil {
		render.Render(w, r, &responses.ErrResponse{Message: "Cannot update company", HTTPStatusCode: http.StatusInternalServerError})
		return
	}

	render.Render(w, r, &responses.CompanyResponse{Company: &new, HTTPStatusCode: http.StatusOK})
}

// HandleCompanyDelete handles DELETE requests to delete companies
func (c *companies) HandleCompanyDelete(w http.ResponseWriter, r *http.Request) {
	company := r.Context().Value(CompanyCtxKey{}).(models.Company)

	if err := c.companiesRepo.Delete(company.ID); err != nil {
		render.Render(w, r, &responses.ErrResponse{Message: "Cannot remove company", HTTPStatusCode: http.StatusInternalServerError})
		return
	}

	render.Render(w, r, &responses.ErrResponse{Message: "Company removed", HTTPStatusCode: http.StatusOK})
}
