package middlewares

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/brokeyourbike/xm-golang-exercise/api/handlers"
	"github.com/brokeyourbike/xm-golang-exercise/api/requests"
	"github.com/brokeyourbike/xm-golang-exercise/pkg/validator"
	"github.com/gorilla/schema"
	"github.com/stretchr/testify/assert"
)

func TestCompanyPayloadCtx(t *testing.T) {
	cases := map[string]struct {
		method     string
		path       string
		body       string
		statusCode int
		response   string
	}{
		"json should be valid": {
			method:     http.MethodPost,
			path:       "/",
			body:       "not-a-json",
			statusCode: http.StatusBadRequest,
			response:   `{"message":"Invalid JSON"}`,
		},
		"valid json will call next": {
			method:     http.MethodPost,
			path:       "/",
			body:       `{"name":"john","code":"123","country":"US","website":"example.com","phone":"+1234567898"}`,
			statusCode: http.StatusOK,
			response:   `the end.`,
		},
		"invalid data should return description": {
			method:     http.MethodPost,
			path:       "/",
			body:       `{"name":"john","country":"US","website":"example.com","phone":"+1234567898"}`,
			statusCode: http.StatusBadRequest,
			response:   `{"message":"Invalid request data","errors":["Field validation for 'Code' failed on the 'required' tag"]}`,
		},
		"query should be valid": {
			method:     http.MethodGet,
			path:       "/?name=1234",
			body:       "",
			statusCode: http.StatusBadRequest,
			response:   ``,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			mw := NewCompanyPayloadCtx(validator.NewValidation(), schema.NewDecoder())
			h := mw.Handle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("the end."))
				w.WriteHeader(http.StatusOK)

				data := r.Context().Value(handlers.CompanyPayloadCtxKey{})
				assert.IsType(t, requests.CompanyPayload{}, data)
			}))

			req := httptest.NewRequest(c.method, c.path, strings.NewReader(c.body))
			w := httptest.NewRecorder()
			h.ServeHTTP(w, req)

			assert.Equal(t, c.statusCode, w.Result().StatusCode)
			assert.Equal(t, c.response, strings.Trim(w.Body.String(), "\n"))
		})
	}
}
