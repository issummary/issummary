#FROM node:9.4.0-stretch AS front-builder
#COPY . /app
##RUN ls -la
#WORKDIR /app
#RUN make build-front

FROM golang:1 AS builder
RUN apt update && apt -y upgrade

# Install nodejs
RUN curl -sL https://deb.nodesource.com/setup_8.x | bash -
RUN apt install -y nodejs

WORKDIR /go/src/github.com/issummary/issummary
COPY Makefile /go/src/github.com/issummary/issummary/Makefile

# Install build tools
RUN make setup

COPY ./static /go/src/github.com/issummary/issummary/static
RUN make install-front

COPY . /go/src/github.com/issummary/issummary
RUN make install

#FROM golang:1-alpine3.7
#COPY --from=builder /go/bin/* /go/bin/
#ENTRYPOINT ["issummary"]
