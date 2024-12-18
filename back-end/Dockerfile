FROM golang:1.21-alpine as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o server main.go

# Use a minimal base image for running the binary
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 3000
CMD ["./server"]