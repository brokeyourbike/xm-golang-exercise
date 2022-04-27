package middlewares

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/brokeyourbike/xm-golang-exercise/api/handlers"
	"github.com/brokeyourbike/xm-golang-exercise/models"
	"github.com/brokeyourbike/xm-golang-exercise/pkg/validator"
	"github.com/go-chi/render"
)

type CompanyPayloadCtx struct{}

func NewCompanyPayloadCtx() *CompanyPayloadCtx {
	return &CompanyPayloadCtx{}
}

// CompanyPayloadCtx is used to validate incoming Company payload data.
// In case payload is invalid, we return formatted errors.
func (c *CompanyPayloadCtx) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var data models.Company

		if json.NewDecoder(r.Body).Decode(&data) != nil {
			render.Render(w, r, &handlers.ErrResponse{
				Message:        "Invalid JSON",
				HTTPStatusCode: http.StatusBadRequest,
			})
			return
		}

		if errs := validator.NewValidation().Validate(&data); len(errs) != 0 {
			render.Render(w, r, &handlers.ErrResponse{
				Message:        "Invalid request data",
				Errors:         errs.Errors(),
				HTTPStatusCode: http.StatusBadRequest,
			})
			return
		}

		ctx := context.WithValue(r.Context(), handlers.CompanyPayloadCtxKey{}, data)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
