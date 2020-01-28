FROM golang:latest
ENV GO111MODULE=on

ADD . /sab.io/escola-service
WORKDIR /sab.io/escola-service/
#RUN ls
COPY go.mod .
COPY go.sum .
RUN go mod download
#ENV GO111MODULE=off
#WORKDIR /go/src/github.com/sab.io/escola-service/cmd/escola-service
RUN go build sab.io/escola-service/cmd/escola-service
WORKDIR /sab.io/escola-service/cmd/escola-service
CMD ["./escola-service"]
