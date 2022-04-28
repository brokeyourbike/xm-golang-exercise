package middlewares

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/brokeyourbike/xm-golang-exercise/configs"
	log "github.com/sirupsen/logrus"
)

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

type Cache interface {
	Get([]byte) ([]byte, error)
	Set([]byte, []byte, int) error
}

// Ipapi is a middleware that fetches the country code for request IP address,
// and verifies if it's allowed to proceed.
type Ipapi struct {
	config     *configs.Config
	httpClient HTTPClient
	cache      Cache
}

// NewIpapi creates an instance of Ipapi middleware.
func NewIpapi(config *configs.Config, httpClient HTTPClient, cache Cache) *Ipapi {
	return &Ipapi{config: config, httpClient: httpClient, cache: cache}
}

// Handle used to fetch the country code for request IP address,
// and verifies if it's allowed to proceed.
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

// getCountryCode handles fetching and storing country code in cache.
// If cache is empty it will call ipapi service for the result.
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

// fetchCountryCode performs HTTP request to the ipapi service
// and returns country code on success.
func (i *Ipapi) fetchCountryCode(ip string) (string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/country", i.config.Ipapi.BaseURL, ip), nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "xm-golang-exercise/0.0.0")

	resp, err := i.httpClient.Do(req)
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

// isCodeAllowed verifies if country code is allowed.
func (i *Ipapi) isCodeAllowed(code string) bool {
	for _, c := range i.config.Ipapi.AllowedCountries {
		if c == code {
			return true
		}
	}
	return false
}
