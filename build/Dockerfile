FROM golang:1.24-alpine AS builder

WORKDIR /src

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /hero-sekai-shop

FROM scratch

COPY --from=builder ./hero-sekai-shop /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/hero-sekai-shop"]
