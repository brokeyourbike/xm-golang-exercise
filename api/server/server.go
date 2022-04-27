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
	router    *chi.Mux
	companies CompaniesHandler
	cmw       Middleware
	pmw       Middleware
	ipmw      Middleware
}

func NewServer(r *chi.Mux, c CompaniesHandler, cmw Middleware, pmw Middleware, ipmw Middleware) *server {
	s := server{router: r, companies: c, cmw: cmw, pmw: pmw, ipmw: ipmw}
	return &s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.routes()
	s.router.ServeHTTP(w, r)
}

func (s *server) Handle(config *configs.Config) {
	s.routes()
	http.ListenAndServe(config.Host+":"+config.Port, s.router)
}

func (s *server) routes() {
	s.router.Use(middleware.RealIP)
	s.router.Use(middleware.Recoverer)
	s.router.Use(render.SetContentType(render.ContentTypeJSON))

	s.router.Route("/companies", func(r chi.Router) {
		r.With(s.ipmw.Handle, s.pmw.Handle).Post("/", s.companies.HandleCompanyCreate)
		r.Route("/{id:[0-9]+}", func(r chi.Router) {
			r.Use(s.cmw.Handle)
			r.With(s.pmw.Handle).Put("/", s.companies.HandleCompanyUpdate)
			r.Delete("/", s.companies.HandleCompanyDelete)
		})
	})
}
