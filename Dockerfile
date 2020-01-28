# Stage 1
FROM golang:1.11
ENV GOPATH /go

# WORKDIR /go/src/app
# COPY Gopkg.toml Gopkg.lock ./
# RUN dep ensure --vendor-only -v

ADD . /go/src/app
WORKDIR /go/src/app/cmd
RUN go build -o main

# Stage 2
FROM ubuntu:18.04

RUN apt update && apt install -y ca-certificates
COPY --from=0 /go/src/app/cmd/main /app/main

CMD /app/main