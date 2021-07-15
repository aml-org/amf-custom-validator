ARG OPENJDK_TAG=8u292
ARG SBT_VERSION=1.5.4

#FROM mozilla/sbt AS CI-JAVA
## Copy content
#COPY . ./src
#WORKDIR ./src
#
## Install make
#RUN apt-get update && apt-get install make
#
## Generate amf.jar
#RUN git clone https://github.com/aml-org/amf.git
#WORKDIR ./amf
#RUN sbt "clientJVM/assembly"
#RUN mv *.jar amf.jar
#ENV AMF_JAR=$PWD/amf.jar
#WORKDIR ../
#
#RUN make ci-java


FROM golang:1.15 AS CI-GO
# Install make
RUN apt-get update && apt-get install make

COPY . /go/src
WORKDIR /go/src

RUN go mod tidy
RUN make ci-go

FROM node:12 AS CI-JS
COPY --from=CI-GO /go/src ./src
WORKDIR ./src

RUN make ci-js

