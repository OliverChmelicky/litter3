#Dockerfile using reflex watcher
#Another option is using https://mikemadisonweb.github.io/2018/03/06/go-autoreload/ REALIZE watcher

FROM golang:1.13

#init
COPY . /go/src/github.com/olo/litter3
WORKDIR /go/src/github.com/olo/litter3
#RUN apk add --update --no-cache ca-certificates git

#init wathcer go-watcher
ENV PATH /usr/local/go/bin:$PATH
RUN go get github.com/canthefason/go-watcher
#ENTRYPOINT ["watcher", "-run", "be/cmd/main.go"]

#init watcher reflex
#RUN ["go", "get", "github.com/cespare/reflex"]
#ENTRYPOINT [ "reflex" ,"-c" ,"reflex.conf", "--decoration=fancy" ]