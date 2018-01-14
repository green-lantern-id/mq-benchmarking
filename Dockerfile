FROM golang:latest AS build

RUN mkdir -p /go/src/github.com/pchaivong/mq-benchmarking
COPY . /go/src/github.com/pchaivong/mq-benchmarking

WORKDIR /go/src/github.com/pchaivong/mq-benchmarking

RUN wget -O /bin/dep https://github.com/golang/dep/releases/download/v0.3.2/dep-linux-amd64 \
 && chmod +x /bin/dep \
 && /bin/dep init -v \
 && /bin/dep ensure \
 && go install

 FROM nimmis/alpine-golang:1.9
 COPY --from=build /go/bin/mq-benchmarking /usr/local/bin/mq-benchmarking
 RUN chmod +x /usr/local/bin/mq-benchmarking
 WORKDIR /usr/local/bin
 ENTRYPOINT [ "/usr/local/bin/mq-benchmarking" ]