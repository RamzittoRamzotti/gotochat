FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /gotochat ./cmd/main.go


FROM alpine:3.21

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /gotochat .
COPY --from=builder /app/static ./static
COPY --from=builder /app/keys ./keys

EXPOSE 8080

CMD ["./gotochat"]
