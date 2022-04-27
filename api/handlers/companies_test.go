package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brokeyourbike/xm-golang-exercise/api/server"
	"github.com/brokeyourbike/xm-golang-exercise/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleCompanyDelete(t *testing.T) {
	cases := map[string]struct {
		companyID  string
		statusCode int
		message    string
		setupMock  func(companiesRepo *mocks.CompaniesRepo, mw *mocks.Middleware)
	}{
		"username is required": {
			companyID:  "10",
			statusCode: http.StatusNotFound,
			message:    "UserName is too short\n",
			setupMock: func(companiesRepo *mocks.CompaniesRepo, mw *mocks.Middleware) {
				mw.On("Handle", mock.Anything).Return(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			companiesRepo := new(mocks.CompaniesRepo)
			mw := new(mocks.Middleware)
			companies := NewCompanies(companiesRepo)
			c.setupMock(companiesRepo, mw)

			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/user/%s", c.companyID), nil)
			w := httptest.NewRecorder()

			srv := server.NewServer(chi.NewMux(), companies, mw, mw, mw)
			srv.ServeHTTP(w, req)

			assert.Equal(t, c.statusCode, w.Result().StatusCode)

			companiesRepo.AssertExpectations(t)
		})
	}
}
