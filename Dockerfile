FROM ubuntu:latest
LABEL authors="87476"

FROM golang:1.23
WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod downlaod && go mod verify

COPY utils .

RUN go build -v -o /usr/src/app ./...

CMD ["app"]