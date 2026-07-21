# ---------- Build Stage ----------
FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o promail .

# ---------- Runtime Stage ----------
FROM alpine:3.22

WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /app/promail .

EXPOSE 8080

CMD ["./promail"]