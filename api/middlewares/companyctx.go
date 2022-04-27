package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/brokeyourbike/xm-golang-exercise/api/handlers"
	"github.com/brokeyourbike/xm-golang-exercise/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	log "github.com/sirupsen/logrus"
)

type CompanyCtx struct {
	companiesRepo handlers.CompaniesRepo
}

func NewCompanyCtx(r handlers.CompaniesRepo) *CompanyCtx {
	return &CompanyCtx{companiesRepo: r}
}

// CompanyCtx is used to load an Company object from
// the URL parameters passed through as the request. In case
// the Company could not be found, we stop here and return a 404.
func (c *CompanyCtx) Handle(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 0)
		if err != nil {
			log.WithFields(log.Fields{"id": chi.URLParam(r, "id")}).Warn("ID URL param invalid")
			render.Render(w, r, &handlers.ErrResponse{Message: "Invalid ID", HTTPStatusCode: http.StatusBadRequest})
			return
		}

		company, err := c.companiesRepo.Get(uint64(id))

		if errors.Is(err, models.ErrCompanyNotFound) {
			render.Render(w, r, &handlers.ErrResponse{Message: "Resource not found", HTTPStatusCode: http.StatusNotFound})
			return
		}

		if err != nil {
			render.Render(w, r, &handlers.ErrResponse{Message: "Cannot query company", HTTPStatusCode: http.StatusInternalServerError})
			return
		}

		ctx := context.WithValue(r.Context(), handlers.CompanyCtxKey{}, company)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
