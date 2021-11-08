#!/usr/bin/env bash
nexus="https://repository-master.mulesoft.org/nexus/content/repositories/snapshots"
module="com/github/amlorg/amf-cli_2.12"
version="5.0.0-SUPER-SECRET-SNAPSHOT"
artifact="amf-cli_2.12-5.0.0-SUPER-SECRET-SNAPSHOT-assembly.jar"

curl $nexus/$module/$version/$artifact --output amf.jar
