FROM golang:1.11.4-alpine AS builder

WORKDIR /go/src/github.com/chonla/oddsvr-api

COPY . .

ARG APP_VERSION

RUN apk add --no-cache musl-dev git gcc \
    && go get ./... \
    && go test ./... \
    && GOOS=linux GOARCH=amd64 go build -o oddsvr-api -ldflags "-X main.AppVersion=${APP_VERSION}"



FROM alpine:latest

ARG ODDSVR_FRONT_BASEURL
ARG ODDSVR_DB
ARG ODDSVR_STRAVA_CLIENT_ID
ARG ODDSVR_STRAVA_CLIENT_SECRET
ARG ODDSVR_JWT_SECRET
ARG ODDSVR_ADDR

ENV ODDSVR_FRONT_BASEURL=${ODDSVR_FRONT_BASEURL}
ENV ODDSVR_DB=${ODDSVR_DB}
ENV ODDSVR_STRAVA_CLIENT_ID=${ODDSVR_STRAVA_CLIENT_ID}
ENV ODDSVR_STRAVA_CLIENT_SECRET=${ODDSVR_STRAVA_CLIENT_SECRET}
ENV ODDSVR_JWT_SECRET=${ODDSVR_JWT_SECRET}
ENV ODDSVR_ADDR=${ODDSVR_ADDR}

EXPOSE 8080

COPY --from=builder /go/src/github.com/chonla/oddsvr-api/oddsvr-api .

CMD ["./oddsvr-api"]
