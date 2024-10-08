SHELL = /bin/bash
MODULE = $(shell go list -m)
BIN = $(CURDIR)/bin
BUILD = $(CURDIR)/build
BUILDRQ = $(CURDIR)/src/github.com/rqlite
GO = go

export PATH := $(CURDIR)/bin:$(PATH)

#DEPS AND TOOLS
$(BIN):
	@mkdir -p $@
$(BIN)/%: | $(BIN); $(info $(M) installing $(PACKAGE)...)
	GOBIN=$(BIN) $(GO) install -mod=vendor $(PACKAGE)
# 


.PHONY: build
build:
	go build -o pthread main.go

clean:
	rm -rf src build pkg bin

rqlite:
	@mkdir -p $(BUILD)
	@mkdir -p $(BUILDRQ)
	cd $(BUILDRQ) && git clone https://github.com/rqlite/rqlite.git
	export GOPATH=$(CURDIR) && cd $(BUILDRQ)/rqlite && CGO_ENABLED=1 go install ./...