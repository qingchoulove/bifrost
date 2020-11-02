BINARY="bifrost"
VERSION=1.0.0
BUILD=`date +%FT%T%z`
OS := $(if $(GOOS),$(GOOS),$(shell go env GOOS))
ARCH := $(if $(GOARCH),$(GOARCH),$(shell go env GOARCH))

PACKAGES=`go mod graph`
GOFILES=`find . -name "*.go" -type f -not -path "./vendor/*"`

default:
	@echo "build ${BINARY} version ${VERSION} time ${BUILD}"
	@rm -rf release/ && mkdir -p release/front
	@GOOS=${OS} GOARCH=${ARCH} go build -o ${BINARY} -tags=jsoniter
	@cd front && yarn build
	@mv ${BINARY} release
	@cp -r front/dist release/front/

list:
	@echo ${PACKAGES}
	@echo ${GOFILES}

fmt:
	@gofmt -s -w ${GOFILES}

install:
	@go mod download

test:
	@go test -cpu=1,2,4 -v -tags integration ./...

clean:
	@rm -rf release

