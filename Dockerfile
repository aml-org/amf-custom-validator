FROM golang:1.19 AS ci-go
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

FROM cypress/included:11.2.0 as ci-js
COPY . ./src
WORKDIR ./src/wrappers/js
COPY --from=ci-go /go/src/wrappers/js/lib/main.wasm.gz ./lib

RUN npm install
RUN npm run build
RUN npm test

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
