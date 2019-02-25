# Makefile variables.
SHELL=/bin/bash -o pipefail
BIN=$(shell go env GOPATH)/bin
CMD=github.com/liftM/fishmon/cmd/fishmon

# Building the CLI.
.PHONY: all deploy
all: $(BIN)/fishmon

$(BIN)/fishmon: $(shell find . -name *.go)
	@echo "-------> Building"
	go build -o $@ $(CMD)

deploy: $(BIN)/fishmon
	@echo "-------> Deploying"
	GOOS=linux GOARCH=arm go build -o /tmp/fishmon-arm $(CMD)
	test -n "$$RPI" || { echo "must set \$$RPI"; exit 1; }
	scp /tmp/fishmon-arm $$RPI:./fishmon
	scp fishmonconfig.json $$RPI:./fishmonconfig.json
