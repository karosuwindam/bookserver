TAG = ${shell cat version}
APPNAME = "bookserver"
GOVERSION = "1.18"
TEMPLATE = ./Dockerfile_tmp
BASE_CONTANER = "debian:11"
TARGET = Dockerfile
NAME = bookserver2:31000/tool/${APPNAME}
TARGET_FILE = ./
MAKEFILE_DIR := $(dir $(lastword $(MAKEFILE_LIST)))

DOCKER = docker

BUILD = buildctl
BUILD_ADDR = tcp://buildkit.bookserver.home:1234 #arm64
BUILD_ADDR_ARM = tcp://buildkit-arm.bookserver.home:1235 #arm
BUILD_OPTION = "type=image,push=true,registry.insecure=true"



ARCH = ${shell arch}
ifeq (${ARCH},x86_64)
ARCH = amd64
BASECONTANA = ubuntu
else
ARCH = armv6l
BASECONTANA = ubuntu
endif

OPT = "--privileged"


test:
	@echo $(MAKEFILE_DIR)
test-run:

create:
	@echo "create dockerfile"
	@echo "for ${NAME}:${TAG}"
	cat ${TEMPLATE} | sed -e "s|TAG|${TAG}|g" | sed -e "s|GOVERSION|${GOVERSION}|g" | sed -e "s|APPNAME|${APPNAME}|g" | sed -e "s|BASE_CONTANER|${BASE_CONTANER}|g" > ${TARGET}
build: create
	@echo build
	${DOCKER} build -t ${NAME}:${TAG} ${TARGET_FILE}
rmi:
	${DOCKER} rmi ${NAME}:${TAG}
	${DOCKER} image prune -f
buildkit: create
        @echo "--- buildkit build --"
        ${BUILD} --addr ${BUILD_ADDR_ARM} build --output name=${NAME}:${TAG},${BUILD_OPTION} --frontend=dockerfile.v0 --local context=${TARGET_FILE}   --local dockerfile=${TARGET_FILE} --opt source=${TARGET_FILE}${TARGET}
help:
	@echo ""