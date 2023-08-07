# Makefile.

PROJECT_PATH=$(shell cd "$(dirname "$0" )" &&pwd)
PROJECT_NAME=omnigram-server
VERSION=$(shell git describe --tags | sed 's/\(.*\)-.*/\1/')
BUILD_DATE=$(shell date -u '+%Y-%m-%d_%I:%M:%S%p')
BUILD_HASH=$(shell git rev-parse HEAD)
LDFLAGS="-X github.com/nexptr/omnigram-server.BuildStamp=${BUILD_DATE} -X github.com/nexptr/omnigram-server.GitHash=${BUILD_HASH} -X github.com/nexptr/omnigram-server.VERSION=${VERSION} -s -w"

DESTDIR=${PROJECT_PATH}/build
VERSION=v0.0.2

ifeq ($(BUILD_TYPE), "generic")
	GENERIC_PREFIX:=generic-
else
	GENERIC_PREFIX:=
endif


.PHONY: all


all : omnigram-server


omnigram-server: 
	@echo "create omnigram-server-${VERSION} "
	@mkdir -p ${DESTDIR}/omnigram-server-${VERSION}/{conf,bin,i18n}

	@echo "copy default configure file"
	@cp -f ${PROJECT_PATH}/conf/conf.yaml ${DESTDIR}/omnigram-server-${VERSION}/conf/conf.yaml

	@echo "build github.com/nexptr/omnigram-server"
	@env  go build -ldflags ${LDFLAGS} -o ${DESTDIR}/llmchain-${VERSION}/bin/app github.com/nexptr/omnigram-server/cmd/omni-server


clean:
	rm -rf ${DESTDIR}
	docker rmi llmchain:${VERSION}

