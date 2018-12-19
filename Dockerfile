FROM golang:alpine as builder

RUN apk --no-cache add git

RUN go get github.com/Masterminds/glide

RUN mkdir -p /go/src/github.com/uber/prototool
ADD . /go/src/github.com/uber/prototool

WORKDIR /go/src/github.com/uber/prototool
RUN glide install

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /build/prototool internal/cmd/prototool/main.go

RUN cd /tmp && \
    wget https://github.com/protocolbuffers/protobuf/releases/download/v3.6.1/protoc-3.6.1-linux-x86_64.zip && \
    mkdir -p /assets/protoc && \
    unzip protoc-3.6.1-linux-x86_64.zip -d /assets/protoc

FROM alpine

RUN apk --no-cache add libc6-compat

COPY --from=builder /build/prototool /usr/bin
COPY --from=builder /assets/protoc /root/.cache/prototool/Linux/x86_64/protobuf/3.6.1

WORKDIR /protobufs

ENTRYPOINT ["/usr/bin/prototool"]
