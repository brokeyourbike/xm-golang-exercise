package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type CompaniesHandler interface {
	HandleCompanyCreate(w http.ResponseWriter, r *http.Request)
	HandleCompanyGetOne(w http.ResponseWriter, r *http.Request)
	HandleCompanyGetAll(w http.ResponseWriter, r *http.Request)
	HandleCompanyUpdate(w http.ResponseWriter, r *http.Request)
	HandleCompanyDelete(w http.ResponseWriter, r *http.Request)
}

type Middleware interface {
	Handle(next http.Handler) http.Handler
}

type server struct {
	router    *chi.Mux
	companies CompaniesHandler
	mwCompany Middleware
	mwPayload Middleware
	mwIp      Middleware
}

func NewServer(r *chi.Mux, c CompaniesHandler, mwCompany Middleware, mwPayload Middleware, mwIp Middleware) *server {
	s := server{router: r, companies: c, mwCompany: mwCompany, mwPayload: mwPayload, mwIp: mwIp}
	return &s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.routes()
	s.router.ServeHTTP(w, r)
}

func (s *server) routes() {
	s.router.Use(middleware.RealIP)
	s.router.Use(middleware.Recoverer)
	s.router.Use(render.SetContentType(render.ContentTypeJSON))

	s.router.Route("/companies", func(r chi.Router) {
		r.With(s.mwIp.Handle, s.mwPayload.Handle).Post("/", s.companies.HandleCompanyCreate)
		r.With(s.mwPayload.Handle).Get("/", s.companies.HandleCompanyGetAll)
		r.Route("/{id:[0-9]+}", func(r chi.Router) {
			r.With(s.mwCompany.Handle).Get("/", s.companies.HandleCompanyGetOne)
			r.With(s.mwCompany.Handle, s.mwPayload.Handle).Put("/", s.companies.HandleCompanyUpdate)
			r.With(s.mwIp.Handle, s.mwCompany.Handle).Delete("/", s.companies.HandleCompanyDelete)
		})
	})
}
