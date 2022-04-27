package handlers

import (
	"net/http"

	"github.com/brokeyourbike/xm-golang-exercise/models"
	"github.com/go-chi/render"
)

// ErrResponse renderer type for handling all sorts of errors.
type ErrResponse struct {
	Err            error    `json:"-"`                // low-level runtime error
	HTTPStatusCode int      `json:"-"`                // http response status code
	Message        string   `json:"message"`          // user-level status message
	Errors         []string `json:"errors,omitempty"` // application-level errors
}

func (e ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

// Payload is the request/response payload for the Company data model.
type Payload struct {
	*models.Company
	HTTPStatusCode int `json:"-"`
}

func (c Payload) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, c.HTTPStatusCode)
	return nil
}
