FROM golang:1.11

ENV GO111MODULE=on

EXPOSE 8080

LABEL maintainer="Artelhin <artelhin@gmail.com>"

WORKDIR $GOPATH/src/test_telegram_bot

COPY . .

RUN go build -v

CMD ["./test_telegram_bot"]