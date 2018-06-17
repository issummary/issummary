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
RUN make CGO_ENABLED=0 install

FROM alpine
RUN apk add --no-cache ca-certificates
EXPOSE 8080
COPY --from=builder /go/bin/* /usr/local/bin/
ENTRYPOINT ["issummary"]
