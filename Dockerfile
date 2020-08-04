FROM golang:1.13-alpine AS builder

RUN apk add make upx
WORKDIR /app
COPY . /app
RUN make build
RUN upx dist/*

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/dist/* /app/
COPY db/ /app/db/

CMD /app/server
