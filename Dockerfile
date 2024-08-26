FROM golang:1.22

LABEL maintainer="Ahuigo <github.com/ahuigo>"

ENV GOPATH /go
RUN apt-get update && apt-get install -y openssl
COPY . /app/2wayssl
WORKDIR /app/2wayssl
RUN make install

