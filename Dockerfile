FROM golang:1.15 AS CI-GO
# Install make
RUN apt-get update && apt-get install make

# First copy dependencies to enable Docker caching them
COPY ./go.mod ./go.sum /go/src/
WORKDIR /go/src
RUN go mod tidy

COPY . .
RUN make ci-go

FROM node:12 AS CI-JS

# First copy dependencies to enable Docker caching them
COPY . ./src
WORKDIR ./src
RUN npm install

# Copy generated WASM
COPY --from=CI-GO /go/src/wrappers/js/lib/main.wasm ./wrappers/js/lib

RUN make ci-js

