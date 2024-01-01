SHELL=/bin/bash -o pipefail

.PHONY: build
build:
	go build -o webex-breakouts main.go
