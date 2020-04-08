ARG GO_VERSION=1.17
FROM golang:${GO_VERSION}-bullseye
RUN apt-get update && apt-get -y install --no-install-recommends libprotobuf-dev protobuf-compiler 

COPY install-tensorflowlib.sh /opt/bin/install-tensorflowlib.sh
R