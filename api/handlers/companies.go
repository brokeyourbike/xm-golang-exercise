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

// ErrResponse renderer type for handling all sorts of errors.
type ErrResponse struct {
	HTTPStatusCode int      `json:"-"`                // http response status code
	Message        string   `json:"message"`          // status message
	Errors         []string `json:"errors,omitempty"` // validation errors
}

func (e ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

// CompanyResponse is the response for the Company data model.
type CompanyResponse struct {
	*models.Company
	HTTPStatusCode int `json:"-"`
}

func (c CompanyResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, c.HTTPStatusCode)
	return nil
}

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

	render.Render(w, r, &CompanyResponse{Company: &company, HTTPStatusCode: http.StatusCreated})
}

// HandleCompanyGetOne handles GET requests to view Company
func (c *companies) HandleCompanyGetOne(w http.ResponseWriter, r *http.Request) {
	company := r.Context().Value(CompanyCtxKey{}).(models.Company)

	render.Render(w, r, &CompanyResponse{Company: &company, HTTPStatusCode: http.StatusOK})
}

// HandleCompanyGetAll handles GET requests to view all companies
func (c *companies) HandleCompanyGetAll(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
}

// HandleCompanyUpdate handles PUT requests to update companies
func (c *companies) HandleCompanyUpdate(w http.ResponseWriter, r *http.Request) {
	data := r.Context().Value(CompanyPayloadCtxKey{}).(models.Company)
	company := r.Context().Value(CompanyCtxKey{}).(models.Company)

	data.ID = company.ID

	if err := c.companiesRepo.Update(data); err != nil {
		render.Render(w, r, &ErrResponse{Message: "Cannot update company", HTTPStatusCode: http.StatusInternalServerError})
		return
	}

	render.Render(w, r, &CompanyResponse{Company: &data, HTTPStatusCode: http.StatusOK})
}

// HandleCompanyDelete handles DELETE requests to delete companies
func (c *companies) HandleCompanyDelete(w http.ResponseWriter, r *http.Request) {
	company := r.Context().Value(CompanyCtxKey{}).(models.Company)

	if err := c.companiesRepo.Delete(company.ID); err != nil {
		render.Render(w, r, &ErrResponse{Message: "Cannot remove company", HTTPStatusCode: http.StatusInternalServerError})
		return
	}

	render.Render(w, r, &ErrResponse{Message: "Company removed", HTTPStatusCode: http.StatusOK})
}
