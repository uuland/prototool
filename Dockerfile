FROM golang:alpine as builder

RUN apk --no-cache add git

RUN go get github.com/Masterminds/glide

RUN mkdir -p /go/src/github.com/uber/prototool
ADD . /go/src/github.com/uber/prototool

WORKDIR /go/src/github.com/uber/prototool
RUN glide install

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /build/prototool internal/cmd/prototool/main.go

FROM alpine

COPY --from=builder /build/prototool /usr/bin

ENTRYPOINT ["/usr/bin/prototool"]
