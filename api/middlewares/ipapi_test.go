package middlewares

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/brokeyourbike/xm-golang-exercise/configs"
	"github.com/brokeyourbike/xm-golang-exercise/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestIpapi(t *testing.T) {
	cases := map[string]struct {
		remoteAddr string
		statusCode int
		setupMock  func(cfg *configs.Config, httpClient *mocks.HTTPClient, cache *mocks.Cache)
	}{
		"invalid RemoteAddr": {
			remoteAddr: "127.0.0.1",
			statusCode: http.StatusInternalServerError,
			setupMock:  func(cfg *configs.Config, httpClient *mocks.HTTPClient, cache *mocks.Cache) {},
		},
		"already in cache not allowed value": {
			remoteAddr: "127.0.0.1:1234",
			statusCode: http.StatusUnauthorized,
			setupMock: func(cfg *configs.Config, httpClient *mocks.HTTPClient, cache *mocks.Cache) {
				cfg.Ipapi.AllowedCountries = []string{"US"}
				cache.On("Get", []byte("127.0.0.1")).Return([]byte("GB"), nil)
			},
		},
		"already in cache allowed value": {
			remoteAddr: "127.0.0.1:1234",
			statusCode: http.StatusOK,
			setupMock: func(cfg *configs.Config, httpClient *mocks.HTTPClient, cache *mocks.Cache) {
				cfg.Ipapi.AllowedCountries = []string{"US"}
				cache.On("Get", []byte("127.0.0.1")).Return([]byte("US"), nil)
			},
		},
		"http client unable to finish request": {
			remoteAddr: "127.0.0.1:1234",
			statusCode: http.StatusInternalServerError,
			setupMock: func(cfg *configs.Config, httpClient *mocks.HTTPClient, cache *mocks.Cache) {
				cache.On("Get", []byte("127.0.0.1")).Return([]byte{}, errors.New("nothing in cache"))
				httpClient.On("Do", mock.Anything).Return(nil, errors.New("unable to perform request"))
			},
		},
		"http client status code not OK": {
			remoteAddr: "127.0.0.1:1234",
			statusCode: http.StatusInternalServerError,
			setupMock: func(cfg *configs.Config, httpClient *mocks.HTTPClient, cache *mocks.Cache) {
				cache.On("Get", []byte("127.0.0.1")).Return([]byte{}, errors.New("nothing in cache"))

				resp := http.Response{StatusCode: http.StatusUnauthorized, Body: io.NopCloser(strings.NewReader("US"))}
				httpClient.On("Do", mock.Anything).Return(&resp, nil)
			},
		},
		"http response allowed value": {
			remoteAddr: "127.0.0.1:1234",
			statusCode: http.StatusOK,
			setupMock: func(cfg *configs.Config, httpClient *mocks.HTTPClient, cache *mocks.Cache) {
				cfg.Ipapi.AllowedCountries = []string{"US"}
				cfg.Ipapi.TTLSeconds = 5

				cache.On("Get", []byte("127.0.0.1")).Return([]byte{}, errors.New("nothing in cache"))

				resp := http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader("US"))}
				httpClient.On("Do", mock.Anything).Return(&resp, nil)

				cache.On("Set", []byte("127.0.0.1"), []byte("US"), 5).Return(nil)
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			cfg := configs.Config{}
			httpClient := new(mocks.HTTPClient)
			cache := new(mocks.Cache)
			c.setupMock(&cfg, httpClient, cache)

			mw := NewIpapi(&cfg, httpClient, cache)
			h := mw.Handle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.RemoteAddr = c.remoteAddr

			w := httptest.NewRecorder()
			h.ServeHTTP(w, req)

			assert.Equal(t, c.statusCode, w.Result().StatusCode)

			httpClient.AssertExpectations(t)
			cache.AssertExpectations(t)
		})
	}
}
