#!/usr/bin/env bash
url=https://repository-master.mulesoft.org/nexus/content/repositories/snapshots/com/github/amlorg/adhoccli_2.12/0.1-SNAPSHOT/adhoccli_2.12-0.1-SNAPSHOT-assembly.jar
curl $url --output amf.jar
