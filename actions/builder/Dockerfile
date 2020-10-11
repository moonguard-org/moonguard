FROM alpine:3.12

# TODO add specific versions of everything ?
ENV GO111MODULE=on
ENV GOPATH=/go

RUN apk add --update go git grpc libc-dev tzdata

RUN go get github.com/golang/protobuf/protoc-gen-go
RUN go get google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.0
# RUN go get github.com/bpicolo/moonguard/cmd/moonguard-gen

ENV PATH="$PATH:$(go env GOPATH)/bin"
