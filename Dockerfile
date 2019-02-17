FROM golang:1.11

LABEL maintainer="HackDFW2019"

COPY . /homecook
WORKDIR /homecook

EXPOSE 8080
EXPOSE 8081
RUN ls -la
ENTRYPOINT ["cmd/homecook/HomeCOOK"]
