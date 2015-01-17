FROM joernweissenborn/aursir4go:0.2.0

MAINTAINER Joern Weissenborn <jowen.weissenborn@gmail.com>



RUN go get github.com/joernweissenborn/propertygraph2go
RUN go get code.google.com/p/go.net/websocket
COPY . /var/local/gopath/src/aursirrt/
WORKDIR /var/local/gopath/src/aursirrt/

ENTRYPOINT ["go", "run", "src/main.go"]

