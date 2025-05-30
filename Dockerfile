# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o quotes-api .

# Run stage
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/quotes-api .

EXPOSE 8080
CMD ["./quotes-api"]