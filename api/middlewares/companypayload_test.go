package middlewares

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/brokeyourbike/xm-golang-exercise/api/handlers"
	"github.com/brokeyourbike/xm-golang-exercise/models"
	"github.com/stretchr/testify/assert"
)

func TestCompanyPayloadCtx(t *testing.T) {
	cases := map[string]struct {
		body       string
		statusCode int
		response   string
	}{
		"json should be valid": {
			body:       "not-a-json",
			statusCode: http.StatusBadRequest,
			response:   `{"message":"Invalid JSON"}`,
		},
		"valid json will call next": {
			body:       `{"name":"john","code":"123","country":"US","website":"example.com","phone":"+1234567898"}`,
			statusCode: http.StatusOK,
			response:   `the end.`,
		},
		"invalid data should return description": {
			body:       `{"name":"john","country":"US","website":"example.com","phone":"+1234567898"}`,
			statusCode: http.StatusBadRequest,
			response:   `{"message":"Invalid request data","errors":["Field validation for 'Code' failed on the 'required' tag"]}`,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			mw := NewCompanyPayloadCtx()
			h := mw.Handle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("the end."))
				w.WriteHeader(http.StatusOK)

				data := r.Context().Value(handlers.CompanyPayloadCtxKey{})
				assert.IsType(t, models.Company{}, data)
			}))

			req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(c.body))
			w := httptest.NewRecorder()
			h.ServeHTTP(w, req)

			assert.Equal(t, c.statusCode, w.Result().StatusCode)
			assert.Equal(t, c.response, strings.Trim(w.Body.String(), "\n"))
		})
	}
}
