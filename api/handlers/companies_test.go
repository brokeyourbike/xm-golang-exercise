package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brokeyourbike/xm-golang-exercise/api/server"
	"github.com/brokeyourbike/xm-golang-exercise/mocks"
	"github.com/stretchr/testify/assert"
)

func TestHandleCompanyDelete(t *testing.T) {
	cases := map[string]struct {
		companyID  uint64
		statusCode int
		message    string
		setupMock  func(companiesRepo *mocks.CompaniesRepo)
	}{}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			companiesRepo := new(mocks.CompaniesRepo)
			companies := NewCompanies(companiesRepo)
			c.setupMock(companiesRepo)

			req := httptest.NewRequest(http.MethodPost, "/user", nil)
			w := httptest.NewRecorder()

			srv := server.NewServer(companies, nil, nil, nil)
			srv.DeleteCompany(w, req)

			assert.Equal(t, c.statusCode, w.Result().StatusCode)

			companiesRepo.AssertExpectations(t)
		})
	}
}
