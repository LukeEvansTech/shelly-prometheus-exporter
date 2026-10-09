FROM golang:1.27-alpine@sha256:738d1cf061836894ff6bb8c33881080ac66de8cf0586615012a0c8f592649cfa AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o shelly-prometheus-exporter ./cmd/exporter

FROM scratch

WORKDIR /

COPY --from=build /app/shelly-prometheus-exporter .

ENTRYPOINT ["/shelly-prometheus-exporter"]