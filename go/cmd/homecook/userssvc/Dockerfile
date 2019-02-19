FROM golang:1.11

LABEL maintainer="HackDFW2019"

COPY . /homecook
WORKDIR /homecook

RUN ls -la
EXPOSE 8080
ENTRYPOINT ["./users"]
