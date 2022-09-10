[![Build](https://github.com/sshilin/microbin/actions/workflows/build.yml/badge.svg)](https://github.com/sshilin/microbin/actions/workflows/build.yml)&nbsp;[![Go Report Card](https://goreportcard.com/badge/github.com/sshilin/microbin)](https://goreportcard.com/report/github.com/sshilin/microbin)&nbsp;[![Coverage Status](https://coveralls.io/repos/github/sshilin/microbin/badge.svg)](https://coveralls.io/github/sshilin/microbin)

Microbin is an http(s) service that inspects any request sent to it. This is can useful for understanding how proxies modified the request.

**Features**
---
- Outputs formatted json
- Upgrades protocol to http2 (ALPN and H2C)
- Exposes Promethus metrics

> **Warning**
> Output may expose sensitive data contained in the request

**Example**
---
```
$ curl http://localhost:8080/foo/bar?p1=1

{
  "host": "ef1ddf2d8af6",
  "remote": "127.0.0.1:57625",
  "proto": "HTTP/1.1",
  "method": "GET",
  "uri": "/foo/bar?p1=1",
  "headers": {
    "Accept": "*/*",
    "User-Agent": "curl/7.83.0"
  }
}
```

**How to install**
---
Docker:

    docker run --rm -d -p 8080:8080 --name microbin ghcr.io/sshilin/microbin:latest

Kubectl:

    kubectl apply -f microbin.yaml

Helm:

    helm repo add microbin https://sshilin.github.io/microbin

    helm install --generate-name  microbin/microbin

**Configuration**
---

| Key                 |  Description                | Default         |
|:--------------------|:----------------------------|:----------------|
| `LISTEN`            | Listen on host:port         | 0.0.0.0:8080    |
| `TLS_ENABLED`       | Enable TLS                  | false           |
| `TLS_KEY_FILE`      | TLS key filepath            | ""              |
| `TLS_CERT_FILE`     | TLS cert filepath           | ""              |

**HTTPS**
---

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
