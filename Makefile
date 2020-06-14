install:
	go install -v

build: 
	mkdir -p dist
	go build -o dist/terraform-provider-sym

test:
	go test ./... -v

fmt:
	go fmt ./... -v

.PHONY: install build test fmt
