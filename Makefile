PROJECT=yochu
ORGANIZATION=giantswarm

SOURCE := $(shell find . -name '*.go')
VERSION := $(shell cat VERSION)
GOPATH := $(shell pwd)/.gobuild
PROJECT_PATH := $(GOPATH)/src/github.com/$(ORGANIZATION)
TEMPLATES=$(shell find . -name '*.tmpl')
GOOS := "linux"
GOARCH := "amd64"

.PHONY=all clean test deps bin run-integration-tests

all: deps $(PROJECT)

clean:
	rm -rf $(GOPATH) $(PROJECT)

test: deps
	docker run \
	    --rm \
	    -v $(shell pwd):/usr/code \
	    -e GOPATH=/usr/code/.gobuild \
	    -e GOOS=$(GOOS) \
	    -e GOARCH=$(GOARCH) \
	    -w /usr/code \
		golang:1.5 go test ./... -cover

# deps
deps: .gobuild .gobuild/bin/go-bindata
	
.gobuild/bin/go-bindata:
	docker run \
	--rm \
	-v $(shell pwd):/usr/code \
	-e GOPATH=/usr/code/.gobuild \
	-e GOOS=linux \
	-e GOARCH=$(GOARCH) \
	-w /usr/code \
	golang:1.5 \
	 go get github.com/jteeuwen/go-bindata/...

.gobuild:
	mkdir -p $(PROJECT_PATH)
	rm -f $(PROJECT_PATH)/$(PROJECT) && cd "$(PROJECT_PATH)" && ln -s ../../../.. $(PROJECT)
	#
	# Fetch public dependencies via `go get`
	# All of the dependencies are listed here to make best use of caching in `builder go get`
	GOPATH=$(GOPATH) builder go get github.com/goamz/goamz/aws
	GOPATH=$(GOPATH) builder go get github.com/goamz/goamz/s3
	GOPATH=$(GOPATH) builder go get github.com/juju/errgo
	GOPATH=$(GOPATH) builder go get github.com/spf13/cobra
	GOPATH=$(GOPATH) builder go get github.com/coreos/go-systemd/dbus

# build
$(PROJECT): $(SOURCE) VERSION
		docker run \
	    --rm \
	    -v $(shell pwd):/usr/code \
	    -e GOPATH=/usr/code/.gobuild \
	    -e GOOS=$(GOOS) \
	    -e GOARCH=$(GOARCH) \
	    -w /usr/code \
		golang:1.5 \
		  /bin/bash -c ".gobuild/bin/go-bindata -pkg templates -o templates/templates_bindata.go templates/ && \
	    go build -a -ldflags \"-X main.projectVersion=$(VERSION) -X main.projectBuild=$(COMMIT)\" -o $(PROJECT)"

fmt:
	gofmt -l -w .

run-integration-tests:
	cd tests/integration/ && vagrant destroy -f || true && vagrant up
	
	cd tests/integration/ && vagrant ssh --command sh -c 'cat - > /home/core/yochu' < ../../yochu

	cd tests/integration/ &&  vagrant ssh --command \
  	  "sudo systemctl enable /etc/systemd/system/yochu.service && \
	  sudo systemctl start yochu.service && \
	  sleep 10 && \
	  systemctl is-active yochu.service"

godoc: 
	@echo Opening godoc server at http://localhost:6060/pkg/github.com/$(ORGANIZATION)/$(PROJECT)/
	@echo If using docker-machine, you will need to use the IP address of your machine
	docker run \
	--rm \
	-v $(shell pwd):/usr/code \
	-e GOPATH=/usr/code/.gobuild \
	-e GOROOT=/usr/code/.gobuild \
	-e GOOS=$(GOOS) \
	-e GOARCH=$(GOARCH) \
	-e GO15VENDOREXPERIMENT=1 \
	-w /usr/code \
	-p 6060:6060 \
	golang:1.5 \
	  godoc -http=:6060 -goroot=/usr/code/.gobuild
	  
bin-dist: all
	mkdir -p bin-dist/
	cp -f README.md bin-dist/
	cp -f LICENSE bin-dist/
	cp $(PROJECT) bin-dist/
	cd bin-dist/ && tar czf $(PROJECT).$(VERSION).tar.gz *
