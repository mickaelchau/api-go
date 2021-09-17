FROM golang:1.17

RUN mkdir /api-go

ADD . /api-go

WORKDIR /api-go

##RUN go mod download 
## Seems to be a good practice to download before the build

RUN go build -o main .

EXPOSE 8000:8000

CMD ["/api-go/main"]

