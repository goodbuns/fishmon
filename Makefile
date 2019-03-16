# Makefile variables.
SHELL=/bin/bash -o pipefail
BIN=$(shell go env GOPATH)/bin
FISHMON=github.com/goodbuns/fishmon/cmd/fishmon
FISHMON=github.com/goodbuns/fishmon/cmd/fmmon

# Building the CLI.
.PHONY: all deploy
all: $(BIN)/fishmon $(BIN)/fmmon

$(BIN)/fishmon: $(shell find . -name *.go)
	@echo "-------> Building fishmon"
	go build -o $@ $(FISHMON)

$(BIN)/fmmon: $(shell find . -name *.go)
	@echo "-------> Building fmmon"
	go build -o $@ $(FMMON)

deploy: $(BIN)/fishmon
	@echo "-------> Deploying fishmon"
	GOOS=linux GOARCH=arm go build -o /tmp/fishmon-arm $(CMD)
	test -n "$$RPI" || { echo "must set \$$RPI"; exit 1; }
	scp /tmp/fishmon-arm $$RPI:./fishmon
	scp fishmonconfig.json $$RPI:./fishmonconfig.json

deploy: $(BIN)/fmmon
	@echo "-------> Deploying fmmon"
	GOOS=linux GOARCH=arm go build -o /tmp/fmmon-arm $(CMD)
	test -n "$$RPI" || { echo "must set \$$RPI"; exit 1; }
	scp /tmp/fmmon-arm $$RPI:./fmmon
