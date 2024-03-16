HOSTNAME=severalnines.com
NAMESPACE=severalnines
NAME=clustercontrol
BINARY=terraform-provider-${NAME}
TARGET=./bin/${BINARY}
VERSION=0.1.0
OS_ARCH=darwin_amd64
TARGET_DIR=~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

.PHONY: docs

default: install

build:
	go build -o ${TARGET}

all: ${TARGET}

${TARGET}:
	 go build -o ${TARGET}

clean:
	rm -rf ${TARGET}

install: ${TARGET}
	mkdir -p ${TARGET_DIR}
	mv ${TARGET} ${TARGET_DIR}
#	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/severalnines/ccx/0.2.0/linux_amd64
#	cp ./bin/terraform-provider-ccx ~/.terraform.d/plugins/registry.terraform.io/severalnines/ccx/0.2.0/linux_amd64/

docs:
	go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

