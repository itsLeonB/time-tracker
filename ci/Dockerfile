FROM golang:1.23 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /time-tracker ./cmd/app/main.go

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /time-tracker /time-tracker

EXPOSE ${APP_PORT:-8000}

USER nonroot:nonroot

ENTRYPOINT ["/time-tracker"]