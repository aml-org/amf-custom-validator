.PHONY: test

all: build test

## TEST ===================================================================================
test: test-go test-js

test-go:
	go test ./internal/...

# must run build-js first
test-js:
	cd ./wrappers/js && npm install && npm test

## BUILD ==================================================================================

build: build-native build-js

build-native:
	rm -f amf-opa-validator
	go build -o amf-opa-validator ./cmd/validator.go

build-js:
	./scripts/gen_js_package.sh

## CI =====================================================================================

ci-go: test-go build-native build-js

ci-js: test-js
