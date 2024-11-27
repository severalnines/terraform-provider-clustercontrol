HOSTNAME=severalnines.com
NAMESPACE=severalnines
NAME=clustercontrol
BINARY=terraform-provider-${NAME}
TARGET=./bin/${BINARY}
VERSION=0.2.21
OS_ARCH=darwin_amd64
TARGET_DIR=~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

.PHONY: docs

default: install

build:
	CGO_ENABLED=0 go build -o ${TARGET}

all: ${TARGET}

clean:
	/bin/rm -rf ${TARGET_DIR}

install: build
	mkdir -p ${TARGET_DIR}
	mv ${TARGET} ${TARGET_DIR}

release:
	goreleaser release --rm-dist --snapshot --skip-publish  --skip-sign

docs:
	go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

