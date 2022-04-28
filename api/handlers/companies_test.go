package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/brokeyourbike/xm-golang-exercise/api/requests"
	"github.com/brokeyourbike/xm-golang-exercise/api/server"
	"github.com/brokeyourbike/xm-golang-exercise/mocks"
	"github.com/brokeyourbike/xm-golang-exercise/models"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type dymmyMw struct{}

func (d dymmyMw) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func TestCompanies(t *testing.T) {
	cases := map[string]struct {
		method     string
		path       string
		statusCode int
		response   string
		setupCtx   func(*http.Request) context.Context
		setupMock  func(companiesRepo *mocks.CompaniesRepo)
	}{
		"view company": {
			method:     http.MethodGet,
			path:       "/companies/30",
			statusCode: http.StatusOK,
			response:   `{"id":30,"name":"tes","code":"tt","country":"US","website":"example.com","phone":"+1234"}`,
			setupCtx: func(req *http.Request) context.Context {
				company := models.Company{ID: 30, Name: "tes", Code: "tt", Country: "US", Website: "example.com", Phone: "+1234"}
				return context.WithValue(req.Context(), CompanyCtxKey{}, company)
			},
			setupMock: func(companiesRepo *mocks.CompaniesRepo) {},
		},
		"create company": {
			method:     http.MethodPost,
			path:       "/companies",
			statusCode: http.StatusCreated,
			response:   `{"id":1,"name":"tes","code":"123","country":"US","website":"example.com","phone":"+1234"}`,
			setupCtx: func(req *http.Request) context.Context {
				data := requests.CompanyPayload{Name: "tes", Code: "123", Country: "US", Website: "example.com", Phone: "+1234"}
				return context.WithValue(req.Context(), CompanyPayloadCtxKey{}, data)
			},
			setupMock: func(companiesRepo *mocks.CompaniesRepo) {
				company := models.Company{ID: 1, Name: "tes", Code: "123", Country: "US", Website: "example.com", Phone: "+1234"}
				companiesRepo.On("Create", mock.AnythingOfType("Company")).Return(company, nil)
			},
		},
		"company cannot be created": {
			method:     http.MethodPost,
			path:       "/companies",
			statusCode: http.StatusInternalServerError,
			response:   `{"message":"Company cannot be created"}`,
			setupCtx: func(req *http.Request) context.Context {
				return context.WithValue(req.Context(), CompanyPayloadCtxKey{}, requests.CompanyPayload{})
			},
			setupMock: func(companiesRepo *mocks.CompaniesRepo) {
				companiesRepo.On("Create", mock.AnythingOfType("Company")).Return(models.Company{}, errors.New("cannot create"))
			},
		},
		"cannot remove company": {
			method:     http.MethodDelete,
			path:       "/companies/30",
			statusCode: http.StatusInternalServerError,
			response:   `{"message":"Cannot remove company"}`,
			setupCtx: func(req *http.Request) context.Context {
				company := models.Company{ID: 30}
				return context.WithValue(req.Context(), CompanyCtxKey{}, company)
			},
			setupMock: func(companiesRepo *mocks.CompaniesRepo) {
				companiesRepo.On("Delete", uint64(30)).Return(errors.New("cannot remove"))
			},
		},
		"company removed": {
			method:     http.MethodDelete,
			path:       "/companies/30",
			statusCode: http.StatusOK,
			response:   `{"message":"Company removed"}`,
			setupCtx: func(req *http.Request) context.Context {
				company := models.Company{ID: 30}
				return context.WithValue(req.Context(), CompanyCtxKey{}, company)
			},
			setupMock: func(companiesRepo *mocks.CompaniesRepo) {
				companiesRepo.On("Delete", uint64(30)).Return(nil)
			},
		},
		"cannot update company": {
			method:     http.MethodPut,
			path:       "/companies/30",
			statusCode: http.StatusInternalServerError,
			response:   `{"message":"Cannot update company"}`,
			setupCtx: func(req *http.Request) context.Context {
				company := models.Company{ID: 30}
				data := requests.CompanyPayload{}

				ctx := context.WithValue(req.Context(), CompanyCtxKey{}, company)
				return context.WithValue(ctx, CompanyPayloadCtxKey{}, data)
			},
			setupMock: func(companiesRepo *mocks.CompaniesRepo) {
				companiesRepo.On("Update", mock.AnythingOfType("Company")).Return(errors.New("cannot update"))
			},
		},
		"company updated": {
			method:     http.MethodPut,
			path:       "/companies/30",
			statusCode: http.StatusOK,
			response:   `{"id":30,"name":"after","code":"a123","country":"AU","website":"after.com","phone":"+56789"}`,
			setupCtx: func(req *http.Request) context.Context {
				company := models.Company{ID: 30, Name: "before", Code: "b123", Country: "BE", Website: "before.com", Phone: "+1234"}
				data := requests.CompanyPayload{Name: "after", Code: "a123", Country: "AU", Website: "after.com", Phone: "+56789"}

				ctx := context.WithValue(req.Context(), CompanyCtxKey{}, company)
				return context.WithValue(ctx, CompanyPayloadCtxKey{}, data)
			},
			setupMock: func(companiesRepo *mocks.CompaniesRepo) {
				companiesRepo.On("Update", mock.AnythingOfType("Company")).Return(nil)
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			companiesRepo := new(mocks.CompaniesRepo)
			companies := NewCompanies(companiesRepo)
			c.setupMock(companiesRepo)

			req := httptest.NewRequest(c.method, c.path, nil)
			w := httptest.NewRecorder()

			mw := dymmyMw{}

			srv := server.NewServer(chi.NewMux(), companies, mw, mw, mw)
			srv.ServeHTTP(w, req.WithContext(c.setupCtx(req)))

			assert.Equal(t, c.statusCode, w.Result().StatusCode)
			assert.Equal(t, c.response, strings.Trim(w.Body.String(), "\n"))

			companiesRepo.AssertExpectations(t)
		})
	}
}
