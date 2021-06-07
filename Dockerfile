FROM golang:1.16.5-alpine as build-env

ARG VERSION=v.dev

WORKDIR /go/src/app

COPY . .

RUN go build -ldflags="-X main.Version=$VERSION" -o /go/bin/app

FROM alpine

COPY --from=build-env /go/bin/app /

CMD ["/app"]