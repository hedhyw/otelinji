ARG GOLANG_DOCKER_TAG=1.22.4-alpine3.17
ARG ALPINE_DOCKER_TAG=3.17

FROM golang:$GOLANG_DOCKER_TAG as builder

RUN apk update && apk upgrade && apk add --no-cache make

WORKDIR /build
COPY . .

RUN go build -o /build/otelinji cmd/otelinji/main.go

FROM alpine:$ALPINE_DOCKER_TAG

WORKDIR /app
COPY --from=builder /build/otelinji /app/otelinji

ENTRYPOINT [ "/app/otelinji" ]
