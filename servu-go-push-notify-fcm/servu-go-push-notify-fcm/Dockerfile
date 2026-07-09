# Build the Go application
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o servu .

# Final runtime image
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/servu ./servu
COPY --from=builder /app/firebase-service-account.json ./firebase-service-account.json
EXPOSE 8001 8002
ENTRYPOINT ["./servu"]
