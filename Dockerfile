FROM golang:1.21 AS ci-go
# Install make
RUN apt-get update && apt-get install make

# First copy dependencies to enable Docker caching them
COPY ./go.mod ./go.sum /go/src/
WORKDIR /go/src
RUN go mod tidy

COPY . .
RUN make ci-go

FROM eclipse-temurin:17-focal AS ci-java

# Copy content
COPY . ./src
WORKDIR ./src

# Install make
RUN apt-get update && apt-get install make

# Install
RUN make ci-java

FROM node:16 AS ci-js

# First copy dependencies to enable Docker caching them
COPY . ./src
WORKDIR ./src/wrappers/js
RUN npm install

# JS-WEB
WORKDIR ../js-web
RUN npm install

WORKDIR ../js

# Copy generated WASM
COPY --from=ci-go /go/src/wrappers/js/lib/main.wasm.gz ./lib

WORKDIR ../../
RUN make ci-js

FROM cypress/included:11.2.0 as ci-browser
COPY --from=ci-js /src ./src
WORKDIR ./src
RUN ./scripts/ci-browser.sh

FROM ci-go AS go-coverage
RUN make go-coverage

FROM sonarsource/sonar-scanner-cli as coverage
COPY --from=go-coverage /go/src/coverage.out .

FROM ci-go AS go-nexus-scan
RUN make generate-list-file

# image used by valkyr -> https://github.com/mulesoft/kilonova-nexusiq-cli
FROM artifacts.msap.io/mulesoft/kilonova-nexusiq-cli:v107.0.2 as nexus-scan
COPY --from=go-nexus-scan /go/src/go.list .

FROM ci-js AS publish-snapshot

COPY . .
RUN chmod -R 777 ./
RUN npm install -g npm-snapshot
RUN make bundle-web-js
