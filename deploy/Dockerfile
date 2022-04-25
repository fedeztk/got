FROM golang:1.15-alpine AS builder
RUN apk add --no-cache translate-shell
WORKDIR ..
COPY . .
RUN go build -o got cmd/got/main.go
CMD ["./got"]
