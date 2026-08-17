FROM golang:1.23-alpine@sha256:383395b794dffa5b53012a212365d40c8e37109a626ca30d6151c8348d380b5f AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o shelly-prometheus-exporter ./cmd/exporter

FROM scratch

COPY --from=build /app/shelly-prometheus-exporter .

ENTRYPOINT ["/shelly-prometheus-exporter"]