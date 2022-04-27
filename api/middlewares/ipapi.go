package middlewares

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/brokeyourbike/xm-golang-exercise/configs"
	log "github.com/sirupsen/logrus"
)

type Ipapi struct {
	config *configs.Config
}

func NewIpapi(config *configs.Config) *Ipapi {
	return &Ipapi{config: config}
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

		code, err := i.findCountryCode(parts[0])
		if err != nil {
			log.WithFields(log.Fields{"ip": parts[0]}).Error("Cannot fetch country for IP")
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

// TODO: cache result?
func (i *Ipapi) findCountryCode(ip string) (string, error) {
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
