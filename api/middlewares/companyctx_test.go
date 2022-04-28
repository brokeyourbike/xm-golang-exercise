package middlewares

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/brokeyourbike/xm-golang-exercise/api/handlers"
	"github.com/brokeyourbike/xm-golang-exercise/mocks"
	"github.com/brokeyourbike/xm-golang-exercise/models"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCompanyCtx(t *testing.T) {
	cases := map[string]struct {
		companyID  string
		statusCode int
		response   string
		setupMock  func(companiesRepo *mocks.CompaniesRepo)
	}{
		"out of bounds ID": {
			companyID:  "100000000000000000000000000000000000000000000000000000000000000",
			statusCode: http.StatusBadRequest,
			response:   `{"message":"Invalid ID"}`,
			setupMock:  func(companiesRepo *mocks.CompaniesRepo) {},
		},
		"company not found": {
			companyID:  "10",
			statusCode: http.StatusNotFound,
			response:   `{"message":"Company not found"}`,
			setupMock: func(companiesRepo *mocks.CompaniesRepo) {
				companiesRepo.On("Get", uint64(10)).Return(models.Company{}, models.ErrCompanyNotFound)
			},
		},
		"cannot query company": {
			companyID:  "10",
			statusCode: http.StatusInternalServerError,
			response:   `{"message":"Cannot query company"}`,
			setupMock: func(companiesRepo *mocks.CompaniesRepo) {
				companiesRepo.On("Get", uint64(10)).Return(models.Company{}, gorm.ErrInvalidData)
			},
		},
		"it can call next": {
			companyID:  "10",
			statusCode: http.StatusOK,
			response:   `the end.`,
			setupMock: func(companiesRepo *mocks.CompaniesRepo) {
				companiesRepo.On("Get", uint64(10)).Return(models.Company{}, nil)
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			companiesRepo := new(mocks.CompaniesRepo)
			c.setupMock(companiesRepo)

			mw := NewCompanyCtx(companiesRepo)
			h := mw.Handle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("the end."))
				w.WriteHeader(http.StatusOK)

				data := r.Context().Value(handlers.CompanyCtxKey{})
				assert.IsType(t, models.Company{}, data)
			}))

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s", c.companyID), nil)
			w := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Get("/{id:[0-9]+}", h.ServeHTTP)
			r.ServeHTTP(w, req)

			assert.Equal(t, c.statusCode, w.Result().StatusCode)
			assert.Equal(t, c.response, strings.Trim(w.Body.String(), "\n"))

			companiesRepo.AssertExpectations(t)
		})
	}
}
