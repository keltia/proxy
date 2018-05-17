
# Main Makefile for proxy
#
# Copyright 2018 © by Ollivier Robert <roberto@keltia.net>
#

SRCS=	proxy.go utils.go

USRCS=	config_unix.go
WSRCS=	config_windows.go

OPTS=	-ldflags="-s -w" -v

all: build

build: ${SRCS} ${USRCS}
	go build ${OPTS}

test: build
	(unset http_proxy https_proxy && go test )

windows: ${SRCS} ${WSRCS}
	GOOS=windows go build ${OPTS} .

install:
	go install ${OPTS} .

lint:
	gometalinter .

clean:
	go clean .

push:
	git push --all
	git push --tags
