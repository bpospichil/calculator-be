# ---- Build stage ----
FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git

WORKDIR /src

# Cache module downloads
COPY go.mod ./
RUN go mod download

COPY . .

# Build a statically-linked binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /app ./cmd/server

# ---- Final stage ----
FROM scratch

COPY --from=builder /app /app

EXPOSE 8080

ENTRYPOINT ["/app"]
