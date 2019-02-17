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

RUN go get -v ./cmd/homecook/...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags='-w -s' -o ./cmd/homecook/HomeCOOK cmd/homecook/main.go

EXPOSE 8080
EXPOSE 8081
ENTRYPOINT ["./cmd/homecook/HomeCOOK"]

# FROM scratch
# COPY --from=builder /bin/HomeCOOK /
# EXPOSE 8080
# EXPOSE 8081
# ENTRYPOINT ["/HomeCOOK"]
