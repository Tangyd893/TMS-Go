FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/tms-api ./cmd/server

FROM alpine:3.20

RUN addgroup -S app && adduser -S app -G app
WORKDIR /app

COPY --from=builder /out/tms-api /app/tms-api

USER app
EXPOSE 8080

ENTRYPOINT ["/app/tms-api"]
