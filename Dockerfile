# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.18.1 as build
WORKDIR /build
ADD . /build/
RUN go build -o app

##
## Run
##
FROM gcr.io/distroless/base
WORKDIR /
COPY --from=build /build/app /bin/app
EXPOSE 9090
USER nonroot:nonroot
ENTRYPOINT ["/bin/app"]