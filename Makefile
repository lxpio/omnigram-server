# Makefile.

PROJECT_PATH=$(shell cd "$(dirname "$0" )" &&pwd)
PROJECT_NAME=omnigram-server
VERSION=$(shell git describe --tags | sed 's/\(.*\)-.*/\1/')
BUILD_DATE=$(shell date -u '+%Y-%m-%d_%I_%M_%S%p')
BUILD_HASH=$(shell git rev-parse HEAD)
LDFLAGS="-X main.BUILDSTAMP=${BUILD_DATE} -X main.GITHASH=${BUILD_HASH} -X main.VERSION=${VERSION} -s -w"
# SHELL := /bin/bash
VERSION=v0.0.2
DESTDIR=${PROJECT_PATH}/build/omnigram-server-${VERSION}




ifeq ($(BUILD_TYPE), "generic")
	GENERIC_PREFIX:=generic-
else
	GENERIC_PREFIX:=
endif


.PHONY: all


all : omnigram-server


omnigram-server: 
	@echo "create omnigram-server-${VERSION} "
	@#debian上直接使用mkdir不会创建，需要额外调用 bash-c 
	@bash -c "mkdir -p ${DESTDIR}/{conf,bin,data}"

	@echo "copy default configure file"
	@cp -f ${PROJECT_PATH}/conf/conf.yaml ${DESTDIR}/conf/conf.yaml

	@echo "build github.com/nexptr/omnigram-server"
	@env  go build -ldflags ${LDFLAGS} -o ${DESTDIR}/bin/omni-server github.com/nexptr/omnigram-server/cmd/omni-server



docker:
	@docker build --build-arg BUILD_DATE=${BUILD_DATE} --build-arg BUILD_HASH=${BUILD_HASH} --build-arg BUILD_HASH=${VERSION} -t omnigram-server:${VERSION} ./

docker_cn:
	@docker buildx build --build-arg BUILD_COUNTRY="CN" --build-arg BUILD_DATE=${BUILD_DATE} --build-arg BUILD_HASH=${BUILD_HASH} --build-arg BUILD_VERSION=${VERSION} -t omnigram-server:${VERSION} ./

clean:
	rm -rf ${DESTDIR}
	docker rmi omnigram-server:${VERSION}

run_docker: docker
	docker run --rm -v ${LOCAL_EPUB_DIR}:/epub ${LOCAL_MATA_DIR}:/metadata omnigram-server:${VERSION}
