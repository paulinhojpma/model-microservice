FROM golang:latest
ENV GO111MODULE=on



#RUN echo "[url \"git@github.com:\"]\n\tinsteadOf = https://github.com/" >> /root/.gitconfig
#RUN mkdir /root/.ssh && echo "StrictHostKeyChecking no " > /root/.ssh/config
ADD . /go/src/github.com/sab.io/escola-service
#WORKDIR /go/src/github.com/sab.io/escola-service/
#RUN go get ./...

WORKDIR /go/src/github.com/sab.io/escola-service/cmd/escola-service

RUN go build
CMD ["./escola-service"]
