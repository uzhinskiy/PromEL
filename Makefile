GOPATH=$(shell pwd)/vendor:$(shell pwd)
GOBIN=$(shell pwd)/build/
GOFILES=./cmd/$(wildcard *.go)
GONAME=$(shell basename "$(PWD)")
PID=/tmp/go-$(GONAME).pid
GOOS=linux
#GOARCH=386
BUILD=`date +%FT%T%z`
MKDIR_P = mkdir -p

all :: build

get:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go get -d $(GOFILES)

build:
	@echo "Building $(GOFILES) to ./build"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags "-s -w -X main.vBuild=${BUILD}" -o build/$(GONAME) $(GOFILES)
	strip ./build/$(GONAME)

install:
	install -m 755 build/$(GONAME) /usr/local/sbin/
	install -m 644 promel.service /etc/systemd/system/
	install -m 644 promel.yml /etc/promel/

run:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go run $(GOFILES)

clear:
	@clear

clean:
	@echo "Cleaning"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go clean

dirs:
	@$(MKDIR_P) /etc/promel
	@$(MKDIR_P) /var/log/promel

.PHONY:	build get install run clean dirs

