package middlewares

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/brokeyourbike/xm-golang-exercise/configs"
	"github.com/coocood/freecache"
	log "github.com/sirupsen/logrus"
)

type Ipapi struct {
	config *configs.Config
	cache  *freecache.Cache
}

func NewIpapi(config *configs.Config, cache *freecache.Cache) *Ipapi {
	return &Ipapi{config: config, cache: cache}
}

// Ipapi is a middleware that logs the start and end of each request.
func (i *Ipapi) Handle(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.RemoteAddr, ":")
		if len(parts) < 2 {
			log.WithFields(log.Fields{"RemoteAddr": r.RemoteAddr}).Error("RemoteAddr format invalid")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		code, err := i.getCountryCode(parts[0])
		if err != nil {
			log.WithFields(log.Fields{"ip": parts[0]}).Error("Cannot find country for IP")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if !i.isCodeAllowed(code) {
			log.WithFields(log.Fields{"code": code}).Warn("Country code not allowed")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func (i *Ipapi) getCountryCode(ip string) (string, error) {
	v, err := i.cache.Get([]byte(ip))
	if err == nil {
		return string(v), nil
	}

	code, err := i.fetchCountryCode(ip)
	if err == nil {
		i.cache.Set([]byte(ip), []byte(code), int(i.config.Ipapi.TTLSeconds))
	}

	return code, err
}

func (i *Ipapi) fetchCountryCode(ip string) (string, error) {
	c := http.Client{Timeout: time.Second * time.Duration(i.config.Ipapi.TimeoutSeconds)}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/country", i.config.Ipapi.BaseURL, ip), nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "xm-golang-exercise/0.0.0")

	resp, err := c.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.WithFields(log.Fields{"statusCode": resp.StatusCode}).Warn("Response status code is not OK")
		return "", fmt.Errorf("response code is not OK: %d", resp.StatusCode)
	}

	code, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(code), nil
}

func (i *Ipapi) isCodeAllowed(code string) bool {
	for _, c := range i.config.Ipapi.AllowedCountries {
		if c == code {
			return true
		}
	}
	return false
}
