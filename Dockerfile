FROM golang:tip-alpine3.22 AS builder

RUN mkdir /src

WORKDIR /src

COPY ./pkg pkg
COPY ./go.mod ./go.sum .
COPY ./main.go .

RUN go mod download
RUN go build -o bot_hsgq main.go

FROM alpine:3.22

COPY --from=builder /src/bot_hsgq /usr/local/bin/bot_hsgq

CMD ["bot_hsgq"]