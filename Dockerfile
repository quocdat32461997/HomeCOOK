FROM golang:1.11 as builder

LABEL maintainer="HackDFW2019"
RUN apt-get update && \
    apt-get -y install git unzip build-essential autoconf libtool

RUN git clone https://github.com/google/protobuf.git && \
    cd protobuf && \
    ./autogen.sh && \
    ./configure && \
    make && \
    make install && \
    ldconfig && \
    make clean && \
    cd .. && \
    rm -r protobuf

RUN adduser -D -g '' homecook

WORKDIR $GOPATH/src/github.com/quocdat32461997/HomeCOOK/
COPY . .
WORKDIR $GOPATH/src/github.com/quocdat32461997/HomeCOOK/cmd/homecook

RUN go get -u google.golang.org/grpc
RUN go get -u github.com/golang/protobuf/protoc-gen-go
RUN go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
RUN go get -u github.com/globalsign/mgo
RUN go get -u github.com/joho/godotenv
RUN go get -u golang.org/x/net
RUN go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger

RUN CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -a -installsuffix cgo -ldflags='-w -s' -o ./HomeCOOK cmd/homecook/main.go
WORKDIR $GOPATH/src/github.com/quocdat32461997/HomeCOOK/

EXPOSE 8080
EXPOSE 8081
ENTRYPOINT ["./cmd/homecook/HomeCOOK"]

# FROM scratch
# COPY --from=builder /bin/HomeCOOK /
# EXPOSE 8080
# EXPOSE 8081
# ENTRYPOINT ["/HomeCOOK"]
