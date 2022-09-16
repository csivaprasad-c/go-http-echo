BUILD_TARGET_ARCH=${BUILD_ARCH}
BIN_NAME=httpecho

ifeq ($(BUILD_TARGET_ARCH),)
        BUILD_TARGET_ARCH=amd64
endif

build:
	env GOOS=linux CGO_ENABLED=0 GOARCH=${BUILD_TARGET_ARCH} go build -o ${BIN_NAME} .

clean:
	rm -rf ${BIN_NAME}