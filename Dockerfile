# builder image
FROM golang:1.12-alpine3.9 as builder

ENV CGO_ENABLED 0
RUN apk --no-cache add git
RUN go get github.com/golang/dep/cmd/dep

WORKDIR /go/src/github.com/bouweceunen/clear-impaired-volumes-taint

COPY src/clear-impaired-volumes-taint/ .

RUN dep ensure -vendor-only

ENV GOARCH amd64

RUN go build -o /bin/clear-impaired-volumes-taint

RUN mkdir /tmp/result/ && \
  cp /bin/clear-impaired-volumes-taint /tmp/result/clear-impaired-volumes-taint

# final image
FROM gcr.io/distroless/base
MAINTAINER Bouwe Ceunen <bouwe.ceunen@gmail.com>
COPY --from=builder /tmp/result/ /
ENTRYPOINT ["/clear-impaired-volumes-taint"]
