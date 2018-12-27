FROM golang:alpine as builder

RUN apk --no-cache add git

ADD . /source
WORKDIR /source

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /build/prototool cmd/prototool/main.go

ENV PROTOC_VER=3.6.1

RUN mkdir -p /assets/protoc && wget -qO- https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VER}/protoc-${PROTOC_VER}-linux-x86_64.zip | \
    unzip -d /assets/protoc - && \
    chmod a+x /assets/protoc/bin/protoc

FROM alpine

RUN apk --no-cache add libc6-compat

ENV PROTOC_VER=3.6.1

COPY --from=builder /build/prototool /usr/bin
COPY --from=builder /assets/protoc /root/.cache/prototool/Linux/x86_64/protobuf/${PROTOC_VER}

WORKDIR /protobufs

ENTRYPOINT ["/usr/bin/prototool"]
