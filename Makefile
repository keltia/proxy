
# Main Makefile for proxy
#
# Copyright 2018 © by Ollivier Robert <roberto@keltia.net>
#

GO=		go
SRCS=	parse.go utils.go

USRCS=	config_unix.go
WSRCS=	config_windows.go

OPTS=	-ldflags="-s -w" -v

all: build

build: ${SRCS} ${USRCS}
	${GO} build ${OPTS}

test: build
	(unset http_proxy https_proxy && go test )

windows: ${SRCS} ${WSRCS}
	GOOS=windows ${GO} build ${OPTS} .

install:
	${GO} install ${OPTS} .

lint:
	gometalinter .

clean:
	${GO} clean .

push:
	git push --all
	git push --tags
