FROM golang:alpine as base

FROM base as builder
WORKDIR /opt/bot

ADD go.mod go.sum main.go ./

RUN go build

FROM base
WORKDIR /opt

COPY --from=builder /opt/bot/bot bot

RUN chmod +x /opt/bot

ENTRYPOINT [ "/opt/bot" ]