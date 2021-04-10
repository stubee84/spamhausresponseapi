FROM golang:1.16

WORKDIR /src
COPY ./ ./

RUN apt-get update && apt-get -y install dnsutils
RUN go mod vendor && go build