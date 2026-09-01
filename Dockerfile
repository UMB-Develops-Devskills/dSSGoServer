# Build stage
FROM golang:1.26.2-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/devskills

# Runtime stage
FROM gcr.io/distroless/static-debian12
WORKDIR /
COPY --from=builder /app/server /server
EXPOSE 8000
CMD ["/server"]