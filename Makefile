MAKEFLAGS += --silent --always-make

sockets:
	socat -d -d pty,raw,echo=0 pty,raw,echo=0

build:
	go build -ldflags "-s -w" .
