package server

import (
	"net/http"

	"github.com/brokeyourbike/xm-golang-exercise/configs"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type CompaniesHandler interface {
	HandleCompanyCreate(w http.ResponseWriter, r *http.Request)
	HandleCompanyUpdate(w http.ResponseWriter, r *http.Request)
	HandleCompanyDelete(w http.ResponseWriter, r *http.Request)
}

type Middleware interface {
	Handle(next http.Handler) http.Handler
}

type server struct {
	companies CompaniesHandler
	cmw       Middleware
	pmw       Middleware
	ipmw      Middleware
}

func NewServer(c CompaniesHandler, cmw Middleware, pmw Middleware, ipmw Middleware) *server {
	s := server{companies: c, cmw: cmw, pmw: pmw, ipmw: ipmw}
	return &s
}

func (s *server) CreateCompany(w http.ResponseWriter, r *http.Request) {
	s.companies.HandleCompanyCreate(w, r)
}

func (s *server) UpdateCompany(w http.ResponseWriter, r *http.Request) {
	s.companies.HandleCompanyUpdate(w, r)
}

func (s *server) DeleteCompany(w http.ResponseWriter, r *http.Request) {
	s.companies.HandleCompanyDelete(w, r)
}

func (s *server) Handle(config *configs.Config, router *chi.Mux) {
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(render.SetContentType(render.ContentTypeJSON))

	router.Route("/companies", func(r chi.Router) {
		r.With(s.ipmw.Handle, s.pmw.Handle).Post("/", s.CreateCompany)
		r.Route("/{id:[0-9]+}", func(r chi.Router) {
			r.Use(s.cmw.Handle)
			r.With(s.pmw.Handle).Put("/", s.UpdateCompany)
			r.Delete("/", s.DeleteCompany)
		})
	})

	http.ListenAndServe(config.Host+":"+config.Port, router)
}
