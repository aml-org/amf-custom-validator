.PHONY: test

all: build test

## TEST ===================================================================================
test: test-go test-js

test-go:
	go test ./internal/...

go-coverage:
	go test -coverprofile=coverage.out ./internal/...

# must run build-js first
test-js:
	cd ./wrappers/js && npm install && npm test

## BUILD ==================================================================================

build: build-native build-js bundle-web-js

build-native:
	rm -f amf-opa-validator
	go build -o amf-opa-validator ./cmd/validate/validate.go

build-js:
	./scripts/gen_js_package.sh

build-js-web: build-js bundle-web-js

bundle-web-js:
	./scripts/bundle_js_web_package.sh

## CI =====================================================================================

ci-go: test-go build-native build-js

ci-js: test-js

ci-java:
	./scripts/download-cli.sh
	./scripts/download-dialects-and-vocabularies.sh
	./scripts/validate-profiles.sh
	#./scripts/validate-reports.sh