# xm-golang-exercise

[![Go Reference](https://pkg.go.dev/badge/github.com/brokeyourbike/xm-golang-exercise.svg)](https://pkg.go.dev/github.com/brokeyourbike/xm-golang-exercise)
[![Go Report Card](https://goreportcard.com/badge/github.com/brokeyourbike/xm-golang-exercise)](https://goreportcard.com/report/github.com/brokeyourbike/xm-golang-exercise)
[![Maintainability](https://api.codeclimate.com/v1/badges/314c92377d5671930b71/maintainability)](https://codeclimate.com/github/brokeyourbike/xm-golang-exercise/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/314c92377d5671930b71/test_coverage)](https://codeclimate.com/github/brokeyourbike/xm-golang-exercise/test_coverage)

XM Golang Exercise - v21.0.0

## How to run

```bash
go run main.go
```

## Configuration

| Enviroment Variable | Description | Default  |
| ------------- |:-------------| :-----|
| `HOST` | application host | `127.0.0.1` |
| `PORT` | application port |   `9090` |
| `DATABASE_DSN` | MySQL database DSN (data shource name) | `u:p@tcp(127.0.0.1:3306)/db?charset=utf8mb4` |
| `CACHE_SIZE_MB` | applicaiton cache size, in MB | `100` |
| `IPAPI_BASE_URL` | base URL for the [ipapi.co](https://ipapi.co) service | `https://ipapi.co` |
| `IPAPI_TTL_SECONDS` | how long to cache information about IP address, in seconds | `10` |
| `IPAPI_TIMEOUT_SECONDS` | HTTP client timeout, in seconds | `10` |
| `IPAPI_ALLOWED_COUNTRIES` | list of allowed countries (comma separated) | `CY` |

## Authors
- [Ivan Stasiuk](https://github.com/brokeyourbike) | [Twitter](https://twitter.com/brokeyourbike) | [LinkedIn](https://www.linkedin.com/in/brokeyourbike) | [stasi.uk](https://stasi.uk)

## License
[MIT License](https://github.com/glocurrency/xm-golang-exercise/blob/main/LICENSE)
