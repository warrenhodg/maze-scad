PRODUCT=maze-scad
FULL_PRODUCT=github.com/warrenhodg/${PRODUCT}
EXECUTABLE=maze-scad
EXECUTABLE_LINUX=${EXECUTABLE}.linux
EXECUTABLE_MAC=${EXECUTABLE}.mac
EXECUTABLE_WINDOWS=${EXECUTABLE}.exe
EXECUTABLES=${EXECUTABLE_LINUX} ${EXECUTABLE_MAC} ${EXECUTABLE_WINDOWS}
GOLANG_DOCKER_IMAGE="golang:1.10"
GO_SRC=/go/src
DEP=${GOPATH}/bin/dep
VENDOR=vendor
VERSION=$(shell cat index.go|grep version |grep -oe '[0-9]\+\.[0-9]\+\.[0-9]\+')
UID=$(shell id -u)
GID=$(shell id -g)
CUR_DIR=$(shell pwd)

all: ${EXECUTABLE}
${EXECUTABLE}: ${EXECUTABLES}

linux: ${EXECUTABLE_LINUX}
mac: ${EXECUTABLE_MAC}
windows: ${EXECUTABLE_WINDOWS}

help:
	@echo "Build normally"
	@echo "\tmake"
	@echo "Build using docker"
	@echo "\tUSE_DOCKER=1 make"

${EXECUTABLE_LINUX}: *.go ${VENDOR}
ifeq (${USE_DOCKER}, 1)
	docker run --user ${UID}:${GID} --rm -v ${CUR_DIR}:${GO_SRC}/${FULL_PRODUCT} -w ${GO_SRC}/${FULL_PRODUCT} golang:1.10 make ${EXECUTABLE_LINUX}
else
	GOOS=linux GOARCH=386 go build -o ${EXECUTABLE}.linux ${FULL_PRODUCT}
endif

${EXECUTABLE_MAC}: *.go ${VENDOR}
ifeq (${USE_DOCKER}, 1)
	docker run --user ${UID}:${GID} --rm -v ${CUR_DIR}:${GO_SRC}/${FULL_PRODUCT} -w ${GO_SRC}/${FULL_PRODUCT} golang:1.10 make ${EXECUTABLE_MAC}
else
	GOOS=darwin GOARCH=386 go build -o ${EXECUTABLE}.mac ${FULL_PRODUCT}
endif

${EXECUTABLE_WINDOWS}: *.go ${VENDOR}
ifeq (${USE_DOCKER}, 1)
	docker run --user ${UID}:${GID} --rm -v ${CUR_DIR}:${GO_SRC}/${FULL_PRODUCT} -w ${GO_SRC}/${FULL_PRODUCT} golang:1.10 make ${EXECUTABLE_WINDOWS}
else
	GOOS=windows GOARCH=386 go build -o ${EXECUTABLE}.exe ${FULL_PRODUCT}
endif

${VENDOR}: ${DEP}
ifneq (${USE_DOCKER}, 1)
	dep ensure
	chmod 777 vendor
endif

${DEP}:
ifneq (${USE_DOCKER}, 1)
	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
endif

clean:
	rm -f ${EXECUTABLES}
