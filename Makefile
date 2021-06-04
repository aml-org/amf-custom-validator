.PHONY: test

all: test build

test-profiles:
	./scripts/check_profile_syntax.sh
test-go:
	go test ./internal/...

test-js: build-js
	cd ./wrappers/js && npm install && ./node_modules/.bin/mocha

test: test-profiles test-go test-js

build-js:
	./scripts/gen_js_package.sh
build-native:
	rm -f amf-opa-validator
	${GO} build -o amf-opa-validator ./cmd/validator.go

build: build-native build-js
