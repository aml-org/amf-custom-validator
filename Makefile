.PHONY: test performance

all: build test

## TEST ===================================================================================
test: test-go test-js

test-go:
	go test ./internal/...

go-coverage:
	go test -coverprofile=coverage.out -coverpkg=./internal/... ./internal/...

# must run build-js first
test-js:
	cd ./wrappers/js && npm install && npm run build && npm test

# must have a downloaded AMF CLI
test-profiles:
	./scripts/validate-profiles.sh

test-reports:
	./scripts/validate-reports.sh

## PERFORMANCE ============================================================================
performance:
	time (go run ./performance/main/performance.go)
	time (go run ./performance/main/performance.go "--pre-compiled")

## BUILD ==================================================================================

build: build-native build-js bundle-web-js

build-native:
	rm -f acv
	go build -o acv ./cmd/main.go

build-js:
	GOOS=js GOARCH=wasm go build -o wrappers/js/lib/main.wasm -ldflags "-s -w" js/validator.go &&\
	gzip -9 -v -c wrappers/js/lib/main.wasm > wrappers/js/lib/main.wasm.gz &&\
    rm wrappers/js/lib/main.wasm

build-js-web: build-js bundle-web-js

## CI =====================================================================================

ci-go: test-go build-native build-js

ci-js: test-js

ci-java:
	./scripts/download-amf-cli.sh
	./scripts/validate-profiles.sh
	./scripts/validate-reports.sh

## Helpers ==========================================================================
validate-profiles:
	./scripts/validate-profiles.sh

generate-list-file:
	go list -m all > go.list

install:
	go build -o ${GOPATH}/bin/acv ./cmd/main.go
