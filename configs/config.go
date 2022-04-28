package configs

// Config is a single storage of configuration options.
type Config struct {
	Host     string `env:"HOST" envDefault:"127.0.0.1"`
	Port     string `env:"PORT" envDefault:"9090"`
	Database struct {
		Dsn string `env:"DATABASE_DSN" envDefault:"gorm:pa55_worD@tcp(localhost:3306)/xm?charset=utf8mb4&parseTime=True&loc=Local"`
	}
	CacheSizeMb int `env:"CACHE_SIZE_MB" envDefault:"100"`
	Ipapi       struct {
		BaseURL          string   `env:"IPAPI_BASE_URL" envDefault:"https://ipapi.co"`
		TTLSeconds       uint     `env:"IPAPI_TTL" envDefault:"10"`
		TimeoutSeconds   uint     `env:"IPAPI_TIMEOUT" envDefault:"10"`
		AllowedCountries []string `env:"IPAPI_ALLOWED_COUNTRIES" envDefault:"CY,Undefined"`
	}
}
