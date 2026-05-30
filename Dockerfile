FROM golang:1.22-alpine AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOFLAGS=-mod=mod go build -o api ./cmd/api

FROM alpine:3.19
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /build/api .
COPY public/GeoLite2-City.mmdb /app/data/
COPY public/GeoLite2-ASN.mmdb /app/data/
COPY public/GeoLite2-Country.mmdb /app/data/
EXPOSE 8080
CMD ["./api"]
