[![Build](https://github.com/sshilin/microbin/actions/workflows/build.yml/badge.svg)](https://github.com/sshilin/microbin/actions/workflows/build.yml)&nbsp;[![Go Report Card](https://goreportcard.com/badge/github.com/sshilin/microbin)](https://goreportcard.com/report/github.com/sshilin/microbin)&nbsp;[![Coverage Status](https://coveralls.io/repos/github/sshilin/microbin/badge.svg)](https://coveralls.io/github/sshilin/microbin)

Microbin replies with http request headers and some additional details on any URL. This can be used in a situation when you need to understand how the intermediate proxies modify the request until it reaches the destination.

## Features
- Output formatted json
- Upgrade protocol to http2 when requested (on ALPN or H2C)
- Expose Promethus metrics

> **Warning**
> Output may expose sensitive data contained in the request

## Example
```
$ curl http://localhost/microbin/headers

{
  "host": "184a0422e442",
  "remote": "172.18.0.6:47626",
  "proto": "HTTP/1.1",
  "method": "GET",
  "uri": "/headers",
  "headers": {
    "Accept": "*/*",
    "Accept-Encoding": "gzip",
    "User-Agent": "curl/7.83.0",
    "X-Forwarded-For": "172.18.0.1",
    "X-Forwarded-Host": "localhost",
    "X-Forwarded-Port": "80",
    "X-Forwarded-Prefix": "/microbin",
    "X-Forwarded-Proto": "http",
    "X-Forwarded-Server": "23569d46b42c",
    "X-Real-Ip": "172.18.0.1"
  }
}
```

## How to install
Docker:

    docker run --rm -d -p 8080:8080 --name microbin ghcr.io/sshilin/microbin:latest

Kubectl:

    kubectl apply -f microbin.yaml

Helm:

    helm repo add microbin https://sshilin.github.io/microbin

    helm install --generate-name  microbin/microbin

## Configuration
| Key                 |  Description                | Default         |
|:--------------------|:----------------------------|:----------------|
| `LISTEN`            | Listen on host:port         | 0.0.0.0:8080    |
| `TLS_ENABLED`       | Enable TLS                  | false           |
| `TLS_KEY_FILE`      | TLS key filepath            | ""              |
| `TLS_CERT_FILE`     | TLS cert filepath           | ""              |

## HTTPS
This example shows how to enable HTTPS with a self signed certificate.

Create cert:

    openssl req -x509 -nodes -newkey rsa:4096 -keyout key.pem -out cert.pem -subj "//CN=localhost" -days 365

**Docker**

Run a container with mounted key and cert files:
```bash
$ docker run --rm -d -p 8080:8080 --name microbin \
  -v "$(pwd)/certs:/var/tls:ro" \
  -e TLS_ENABLED=true \
  -e TLS_KEY_FILE=/var/tls/key.pem \
  -e TLS_CERT_FILE=/var/tls/cert.pem \
  ghcr.io/sshilin/microbin:latest
```

**Kubernetes**

First create `microbin-certs` secret with cert and key files, then create `microbin-https.yaml`

    kubectl create secret generic microbin-certs --from-file=./certs

    kubectl apply -f microbin-https.yaml
