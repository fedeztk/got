FROM golang:1.18-alpine AS builder
WORKDIR ..
COPY . .
RUN go build -ldflags "-X main.gotVersion=0.1" -o got cmd/got/main.go
CMD ["./got"]
