FROM golang:1.15 AS ci-go
# Install make
RUN apt-get update && apt-get install make

# First copy dependencies to enable Docker caching them
COPY ./go.mod ./go.sum /go/src/
WORKDIR /go/src
RUN go mod tidy

COPY . .
RUN make ci-go

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
RUN make build-js-web