FROM golang:1.15-alpine AS builder
RUN apk add --no-cache translate-shell
WORKDIR /build
COPY . .
RUN go build -o got
# enable tty
CMD ["./got"]
