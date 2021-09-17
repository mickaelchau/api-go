FROM golang:1.17

WORKDIR /api-go

ADD . /api-go

RUN go build -o main . 

CMD ["/api-go/main"]
