FROM golang:alpine AS builder

# Install gcc and standard C library for SQLite CGO
RUN apk add --no-cache gcc musl-dev

WORKDIR /app
COPY go.mod go.sum ./
# Downgrade go version requirement in go.mod to ensure compatibility with the container's Go version
RUN go mod edit -go=1.21
RUN go mod download

COPY *.go ./
# CGO_ENABLED=1 is required for go-sqlite3
RUN CGO_ENABLED=1 GOOS=linux go build -o logmanager .

FROM alpine:latest
RUN apk add --no-cache sqlite-libs tzdata

WORKDIR /app
COPY --from=builder /app/logmanager /app/logmanager

ENV PORT=8080
ENV DB_PATH=/app/data/logs.db

EXPOSE 8080
EXPOSE 514/udp

CMD ["/app/logmanager"]
