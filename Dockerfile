FROM golang:latest AS build

RUN mkdir -p /go/src/github.com/green-lantern-id/mq-benchmarking
COPY . /go/src/github.com/green-lantern-id/mq-benchmarking

WORKDIR /go/src/github.com/green-lantern-id/mq-benchmarking

RUN wget -O /bin/dep https://github.com/golang/dep/releases/download/v0.3.2/dep-linux-amd64 \
 && chmod +x /bin/dep \
 && /bin/dep ensure

RUN CGO_ENABLED=0 go build -a -installsuffix cgo
# && CGO_ENABLED=0 make DESTDIR=/opt PREFIX=/mq-benchmarking GOFLAGS='-ldflags="-s -w"' install

 FROM alpine:3.6
 COPY --from=build /go/src/github.com/green-lantern-id/mq-benchmarking/mq-benchmarking /usr/local/bin/mq-benchmarking
 RUN chmod +x /usr/local/bin/mq-benchmarking
 WORKDIR /usr/local/bin
 ENTRYPOINT [ "/usr/local/bin/mq-benchmarking" ]
