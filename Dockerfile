FROM golang:1.26-alpine@sha256:3889b425f035be855a72fb4755265311293b6d414521f0a519d819df32222d83 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o shelly-prometheus-exporter ./cmd/exporter

FROM scratch

COPY --from=build /app/shelly-prometheus-exporter .

ENTRYPOINT ["/shelly-prometheus-exporter"]