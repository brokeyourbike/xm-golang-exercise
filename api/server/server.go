package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

// CompaniesHandler defines a set of handlers required
// to be used with the server.
type CompaniesHandler interface {
	HandleCompanyCreate(http.ResponseWriter, *http.Request)
	HandleCompanyGetOne(http.ResponseWriter, *http.Request)
	HandleCompanyGetAll(http.ResponseWriter, *http.Request)
	HandleCompanyUpdate(http.ResponseWriter, *http.Request)
	HandleCompanyDelete(http.ResponseWriter, *http.Request)
}

// Middleware defines a requirements for the middlewares
// to be used with the server.
type Middleware interface {
	Handle(http.Handler) http.Handler
}

// server represents mux.
type server struct {
	router    *chi.Mux
	companies CompaniesHandler
	mwCompany Middleware
	mwPayload Middleware
	mwIp      Middleware
}

// NewServer creates a new server with the given router and handlers.
func NewServer(r *chi.Mux, c CompaniesHandler, mwCompany Middleware, mwPayload Middleware, mwIp Middleware) *server {
	s := server{router: r, companies: c, mwCompany: mwCompany, mwPayload: mwPayload, mwIp: mwIp}
	return &s
}

// ServeHTTP handles serving the routes.
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.routes()
	s.router.ServeHTTP(w, r)
}

// routes defines routes and middlewares.
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
