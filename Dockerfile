FROM golang:1.18-alpine AS builder
WORKDIR /app
COPY . /app
RUN go build -ldflags "-X main.gotVersion=0.1" -o got /app/cmd/got/main.go

FROM alpine:3.15
COPY --from=builder /app/got /
CMD ["/got"]
