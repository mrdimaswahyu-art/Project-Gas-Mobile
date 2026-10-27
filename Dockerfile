# syntax=docker/dockerfile:1
FROM golang:1.27-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server main.go

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/server .
EXPOSE 8080
CMD ["./server"]
