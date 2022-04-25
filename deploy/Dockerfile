FROM golang:1.18-alpine AS builder
WORKDIR ..
COPY . .
RUN go build -o got cmd/got/main.go
CMD ["./got"]
