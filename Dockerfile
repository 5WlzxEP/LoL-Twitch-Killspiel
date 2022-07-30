FROM golang:latest

WORKDIR /go/src/app

RUN mkdir "cmd"

COPY src/Killspiel/ ./

RUN go install cmd/main.go

RUN mkdir "/data"
VOLUME ["/data"]

WORKDIR /data

CMD ["main"]