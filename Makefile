############## COMMON VARS SECTION ##############
APP_NAME = sego

############## TARGETS SECTION ##############
.PHONY: all test clean

build: # @HELP builds for current GOOS/GOARCH
build:
	@CGO_ENABLED=0 go build -ldflags "-s -w" -o $(APP_NAME) main.go

snapshot: # @HELP generate a snapshot for all OS_ARCH combinations
snapshot:
	@goreleaser --snapshot --skip-publish

release: # @HELP releases a new version for all OS_ARCH combinations
release:
	@goreleaser release

dep-upgrade: # @HELP upgrades all dependencies
dep-upgrade:
	@go get -u ./...
	@go mod tidy

clean: # @HELP removes built binaries and temporary files
clean:
	@rm -rf dist

bench: # @HELP executes benchmarks
bench:
	@go test -v -bench=Bench -benchmem ./...

help: # @HELP prints this message
help:
	@echo "TARGETS:"
	@grep -E '^.*: *# *@HELP' Makefile            \
	    | awk '                                   \
	        BEGIN {FS = ": *# *@HELP"};           \
	        { printf "  %-30s %s\n", $$1, $$2 };  \
	    '
