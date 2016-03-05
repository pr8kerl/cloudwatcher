GOROOT := /usr/local/go
GOPATH := $(shell pwd)
GOBIN  := $(GOPATH)/bin
PATH   := $(GOROOT)/bin:$(PATH)
GO=$(GOROOT)/bin/go
DEPS   := github.com/aws/aws-sdk-go/aws github.com/aws/aws-sdk-go/service/cloudwatch

LDFLAGS := -ldflags "-X main.commit=`git rev-parse HEAD`"

all: cloudwatcher

deps: $(DEPS)
	GOPATH=$(GOPATH) go get -u $^

cloudwatcher: main.go
    # always format code
		GOPATH=$(GOPATH) go fmt $^
    # vet it
		GOPATH=$(GOPATH) go tool vet $^
		# binary
		GOPATH=$(GOPATH) go build $(LDFLAGS) -o $@ -v $^
		touch $@

win64: main.go config.go object.go
		# always format code
		GOPATH=$(GOPATH) $(GO) fmt $^
		# vet it
		GOPATH=$(GOPATH) $(GO) tool vet $^
		# binary
		GOOS=windows GOARCH=amd64 GOPATH=$(GOPATH) go build $(LDFLAGS) -o server-win-amd64.exe -v $^
		touch server-win-amd64.exe

.PHONY: $(DEPS) clean

clean:
	rm -f server server-win-amd64.exe
