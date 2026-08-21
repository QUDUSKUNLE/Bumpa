FROM golang:1.23-alpine AS builder

WORKDIR /src
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bumpa .

FROM alpine:3.20

RUN adduser -D -H appuser
WORKDIR /app
COPY --from=builder /app/bumpa ./bumpa
USER appuser

EXPOSE 8080
ENTRYPOINT ["/app/bumpa"]
