SOURCE_PATH :=TEST_SOURCE_PATH=$(PWD)
CURDIR := $(shell pwd)
GOBIN := $(CURDIR)/bin/
ENV:=GOBIN=$(GOBIN)
DIR:=FILE_DIR=$(CURDIR)/testfiles TEST_SOURCE_PATH=$(CURDIR)
GODEBUG:=GODEBUG=gocacheverify=1

install: clean mod
	@echo "Environment installed"

test-cover:
	rm -f coverage.out
	$(DIR) $(GODEBUG) go test -coverprofile=coverage.out -cover -race ./
	go tool cover -html=coverage.out -o coverage.html
	rm coverage.out

test-tc:
	$(DIR) $(GODEBUG) go test --check.format=teamcity -race ./

clean:
	rm -f coverage.out
	rm -fr ./vendor

# Download all dependencies
mod:
	@echo "======================================================================"
	@echo "Run MOD"
	GO111MODULE=on GONOSUMDB="*" GOPROXY=direct go mod verify
	GO111MODULE=on GONOSUMDB="*" GOPROXY=direct go mod tidy -v
	GO111MODULE=on GONOSUMDB="*" GOPROXY=direct go mod vendor -v
	GO111MODULE=on GONOSUMDB="*" GOPROXY=direct go mod download
	GO111MODULE=on GONOSUMDB="*" GOPROXY=direct go mod verify
	@echo "======================================================================"
