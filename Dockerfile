FROM golang:1.18-alpine AS builder
# Needed for text-to-speech to build
RUN apk update && apk add git pkgconfig alsa-lib-dev musl-dev gcc
WORKDIR /app
COPY . /app
RUN go generate ./...
RUN go build -o got /app/cmd/got/main.go

FROM alpine:3.15
# Needed for text-to-speech to play, note that even so it may not work.
# This isn't considered critical and it may or may not be fixed, since
# audio can easily be tested by building the binary (see: make)
RUN apk add alsa-lib-dev alsa-utils
COPY --from=builder /app/got /
CMD ["/got"]
