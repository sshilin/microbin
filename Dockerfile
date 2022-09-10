FROM golang:1.19-alpine as builder

ENV CGO_ENABLED=0

WORKDIR /src

RUN apk add --no-cache --update git

COPY go.mod go.sum ./

RUN go mod download

COPY . /src

RUN git_rev=$(git rev-parse --abbrev-ref HEAD)-$(git log -1 --format=%h) && \
    cd app && go build -trimpath -ldflags="-w -s -X main.build=${git_rev}" -o /build/app

RUN echo "nobody:x:65534:65534:nobody:/:" > /tmp/passwd

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /tmp/passwd /etc/passwd
COPY --from=builder /build/app /app

USER nobody

ENTRYPOINT ["/app"]