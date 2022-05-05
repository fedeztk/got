all:
	go build -ldflags "-X main.gotVersion=0.1" -o got cmd/got/main.go

run:
	go run cmd/got/main.go

clean:
	@if [ -f got ] && [ -x got ]; then \
		rm got; \
	fi

docker-build:
	docker build -t got deploy

docker-run:
	docker run -it -e "TERM=xterm-256color" got

install:
	go build -ldflags "-X main.gotVersion=0.1" -o got cmd/got/main.go
	mv -f got `go env GOPATH`/bin/

uninstall:
	rm `go env GOPATH`/bin/got
