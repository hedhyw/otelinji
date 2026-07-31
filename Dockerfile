FROM golang:1.26.4-alpine3.22 AS builder

RUN apk update && apk upgrade && apk add --no-cache make

WORKDIR /build
COPY . .

RUN go build -o /build/otelinji cmd/otelinji/main.go

FROM alpine:3.22

WORKDIR /app
COPY --from=builder /build/otelinji /app/otelinji

ENTRYPOINT [ "/app/otelinji" ]
