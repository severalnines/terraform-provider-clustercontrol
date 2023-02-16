TEST?=$$(go list ./... | grep -v 'vendor')
HOSTNAME=severalnines.com
NAMESPACE=severalnines
#NAME=cc
NAME=clustercontrol
BINARY=terraform-provider-${NAME}
VERSION=0.0.1
OS_ARCH=darwin_amd64
TARGET_DIR=~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

default: install

build:
	go build -o ${BINARY}

release:
	goreleaser release --rm-dist --snapshot --skip-publish  --skip-sign

install: build
	mkdir -p ${TARGET_DIR}
	mv ${BINARY} ${TARGET_DIR}

test: 
	go test -i $(TEST) || exit 1                                                   
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4                    

testacc: 
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m   

clean:
	/bin/rm -rf ${TARGET_DIR}
