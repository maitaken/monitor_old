SRCS := $(shell find . -type f -name '*.go')

export GO111MODULE=on

GOCMD=go
GOPHER='ʕ◔ϖ◔ʔ'

VERSION=v0.2.0
LDFLAGS=-ldflags "-X main.Version=$(VERSION)"

install:
	go install ${LDFLAGS} ./cmd/monitor

fmt:
	goimports -w $(SRCS)
