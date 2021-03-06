package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/brokeyourbike/xm-golang-exercise/api/handlers"
	"github.com/brokeyourbike/xm-golang-exercise/api/middlewares"
	"github.com/brokeyourbike/xm-golang-exercise/api/server"
	"github.com/brokeyourbike/xm-golang-exercise/configs"
	"github.com/brokeyourbike/xm-golang-exercise/db"
	"github.com/brokeyourbike/xm-golang-exercise/models"
	"github.com/brokeyourbike/xm-golang-exercise/pkg/validator"
	"github.com/caarlos0/env/v6"
	"github.com/coocood/freecache"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/schema"
	gorm_logrus "github.com/onrik/gorm-logrus"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	if err := run(); err != nil {
		log.Fatalf("%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	cfg := configs.Config{}
	if err := env.Parse(&cfg); err != nil {
		return fmt.Errorf("cannot parse config: %v", err)
	}

	orm, err := gorm.Open(mysql.Open(cfg.Database.Dsn), &gorm.Config{Logger: gorm_logrus.New()})
	if err != nil {
		return fmt.Errorf("cannot connect to DB: %v", err)
	}

	cache := freecache.NewCache(cfg.CacheSizeMb * 1024 * 1024)
	httpClient := http.Client{Timeout: time.Second * time.Duration(cfg.Ipapi.TimeoutSeconds)}

	orm.AutoMigrate(&models.Company{})
	companiesRepo := db.NewCompaniesRepo(orm)

	c := handlers.NewCompanies(companiesRepo)
	cmw := middlewares.NewCompanyCtx(companiesRepo)
	pmw := middlewares.NewCompanyPayloadCtx(validator.NewValidation(), schema.NewDecoder())
	ipmw := middlewares.NewIpapi(&cfg, &httpClient, cache)

	mux := server.NewServer(chi.NewRouter(), c, cmw, pmw, ipmw)
	mux.ListenAndServe(&cfg)

	return nil
}
