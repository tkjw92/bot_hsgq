FROM tip-alpine3.22 AS builder

RUN mkdir /src

WORKDIR /src

COPY ./pkg /src/pkg
COPY ./go.mod /src/go.mod
COPY ./go.sum /src/go.mod
COPY ./main.go /src/main.go

RUN go mod download
RUN go build -o bot_hsgq

FROM alpine:3.22

COPY --from=builder /src/bot_hsgq /usr/local/bin/bot_hsgq

CMD ["bot_hsgq"]