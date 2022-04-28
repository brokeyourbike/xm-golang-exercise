package responses

import (
	"net/http"

	"github.com/brokeyourbike/xm-golang-exercise/models"
	"github.com/go-chi/render"
)

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

type CompaniesResponse struct {
	Companies      []models.Company `json:"companies"`
	HTTPStatusCode int              `json:"-"`
}

func (c CompaniesResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, c.HTTPStatusCode)
	return nil
}
