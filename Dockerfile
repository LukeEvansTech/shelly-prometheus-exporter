FROM golang:1.27-alpine@sha256:26402d86be3d72e6a9410afa0108f03529f51f0c1b5eb7f503d0bc44cc7857ac AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o shelly-prometheus-exporter ./cmd/exporter

FROM scratch

WORKDIR /

COPY --from=build /app/shelly-prometheus-exporter .

ENTRYPOINT ["/shelly-prometheus-exporter"]