all:
	go run cmd/got/main.go

build:
	go build -o got cmd/got/main.go

clean:
	@if [ -f got ] && [ -x got ]; then \
		rm got; \
	fi

docker-build:
	docker build -t got deploy

docker-run:
	docker run -it -e "TERM=xterm-256color" got
