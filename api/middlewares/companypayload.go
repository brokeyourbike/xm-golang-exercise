package middlewares

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/brokeyourbike/xm-golang-exercise/api/handlers"
	"github.com/brokeyourbike/xm-golang-exercise/api/requests"
	"github.com/brokeyourbike/xm-golang-exercise/api/responses"
	"github.com/brokeyourbike/xm-golang-exercise/pkg/validator"
	"github.com/go-chi/render"
	"github.com/gorilla/schema"
)

// CompanyPayloadCtx is used to validate incoming Company payload data.
// In case payload is invalid, we return formatted errors.
type CompanyPayloadCtx struct {
	validator    *validator.Validation
	queryDecoder *schema.Decoder
}

func NewCompanyPayloadCtx(v *validator.Validation, queryDecoder *schema.Decoder) *CompanyPayloadCtx {
	return &CompanyPayloadCtx{validator: v, queryDecoder: queryDecoder}
}

func (c *CompanyPayloadCtx) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var data requests.CompanyPayload

		if r.Method == http.MethodGet {
			err := c.queryDecoder.Decode(&data, r.URL.Query())
			if err != nil {
				render.Render(w, r, &responses.ErrResponse{
					Message:        "Invalid query params",
					HTTPStatusCode: http.StatusBadRequest,
				})
				return
			}
		} else {
			if json.NewDecoder(r.Body).Decode(&data) != nil {
				render.Render(w, r, &responses.ErrResponse{
					Message:        "Invalid JSON",
					HTTPStatusCode: http.StatusBadRequest,
				})
				return
			}

			if errs := c.validator.Validate(&data); len(errs) != 0 {
				render.Render(w, r, &responses.ErrResponse{
					Message:        "Invalid request data",
					Errors:         errs.Errors(),
					HTTPStatusCode: http.StatusBadRequest,
				})
				return
			}
		}

		ctx := context.WithValue(r.Context(), handlers.CompanyPayloadCtxKey{}, data)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
