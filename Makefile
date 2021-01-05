TEST?=$$(go list ./... | grep -v 'vendor')
HOSTNAME=symops.io
NAMESPACE=com
NAME=sym
BINARY=terraform-provider-${NAME}
VERSION=0.0.1
OS_ARCH=darwin_amd64
PLUGIN_DIR=~/.terraform.d/plugins/terraform.${HOSTNAME}/symopsio/${NAME}/${VERSION}/${OS_ARCH}

default: build

install:
	go install -v

clean:
	rm -rf dist/*

build:
	mkdir -p dist
	go build -o dist/${BINARY}

# Copy to plugin direction in v12 and v13 formats
local: build
	mkdir -p ${PLUGIN_DIR}
	cp dist/${BINARY} ${PLUGIN_DIR}/${BINARY}_v${VERSION}

test:
	go test -i $(TEST) || exit 1
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc:
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

.PHONY: install clean build local test testacc
