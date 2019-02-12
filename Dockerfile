FROM golang:alpine AS buildocker
ADD . /go/src/wlx
WORKDIR /go/src/wlx
RUN apk update && apk add git
RUN go get -u github.com/golang/dep/cmd/dep && dep ensure && go build -o wlx main.go

FROM alpine
COPY --from=buildocker /go/src/wlx/wlx /usr/local/bin/wlx
ENTRYPOINT ["/usr/local/bin/wlx"]
