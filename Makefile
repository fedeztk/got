all:
	go generate ./...
	go build -o got cmd/got/main.go

run:
	go run cmd/got/main.go

clean:
	@if [ -f got ] && [ -x got ]; then \
		rm got; \
	fi

docker-build:
	docker build -t got .

docker-run:
	docker run -it -e "TERM=xterm-256color" got

install:
	go generate ./...
	go build -o got cmd/got/main.go
	mv -f got `go env GOPATH`/bin/

uninstall:
	rm `go env GOPATH`/bin/got
