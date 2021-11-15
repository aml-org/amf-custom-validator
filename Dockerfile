FROM golang:1.15 AS ci-go
# Install make
RUN apt-get update && apt-get install make

# First copy dependencies to enable Docker caching them
COPY ./go.mod ./go.sum /go/src/
WORKDIR /go/src
RUN go mod tidy

COPY . .
RUN make ci-go

FROM ci-go AS go-coverage
RUN make go-coverage

FROM sonarsource/sonar-scanner-cli as coverage
COPY --from=go-coverage /go/src/coverage.out .

FROM openjdk:8u292-jre AS ci-java
#FROM openjdk/8u292-jdk AS ci-java
# Copy content
COPY . ./src
WORKDIR ./src

# Install make
RUN apt-get update && apt-get install make

# Install
RUN make ci-java

FROM node:12 AS ci-js

# First copy dependencies to enable Docker caching them
COPY . ./src
WORKDIR ./src/wrappers/js
RUN npm install

# Copy generated WASM
COPY --from=ci-go /go/src/wrappers/js/lib/main.wasm.gz ./lib

WORKDIR ../../
RUN make ci-js

FROM ci-js AS publish-snapshot

COPY . .
RUN chmod -R 777 ./
RUN npm install -g npm-snapshot
RUN npm install -g npm-cli-login
RUN make bundle-web-js