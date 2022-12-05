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
	cd ./wrappers/js && npm install && npm test

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
	./scripts/gen_js_package.sh

build-js-web: build-js bundle-web-js

bundle-web-js:
	./scripts/bundle_js_web_package.sh

## CI =====================================================================================

ci-go: test-go build-native build-js

ci-js: test-js

ci-browser:
	./scripts/ci-browser.sh

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
