.PHONY: test

all: test build

## TEST ===================================================================================
test: test-profiles test-go test-js

test-profiles:
	./scripts/check_profile_syntax.sh

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

ci-java: test-profiles

ci-go: test-go build-native build-js

ci-js: test-js
