FROM golang:1.19-alpine as build-env

COPY . /build

WORKDIR /build/app

RUN apk add --no-cache --update git

RUN revision=$(git rev-parse --abbrev-ref HEAD)-$(git log -1 --format=%h) && \
    go build -ldflags="-X main.build=${revision}" -o microbin

FROM alpine:3.16

RUN adduser -H -S microbin -G users

USER microbin

COPY --from=build-env /build/app/microbin /

CMD ["/microbin"]