export GO111MODULE=on

GOCMD=go
GOPHER='ʕ◔ϖ◔ʔ'

VERSION=v0.0.5
LDFLAGS=-ldflags "-X main.Version=$(VERSION)"

install: devs
	cd ./cmd/monitor; go install ${LDFLAGS}
	
devs:
	go get
