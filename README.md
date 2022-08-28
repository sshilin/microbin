[![Build](https://github.com/sshilin/microbin/actions/workflows/build.yml/badge.svg)](https://github.com/sshilin/microbin/actions/workflows/build.yml)&nbsp;[![Go Report Card](https://goreportcard.com/badge/github.com/sshilin/microbin)](https://goreportcard.com/report/github.com/sshilin/microbin)&nbsp;[![Coverage Status](https://coveralls.io/repos/github/sshilin/microbin/badge.svg)](https://coveralls.io/github/sshilin/microbin)

Microbin sends back request headers and k8s pod metainfo in json format. This can be useful for quick checking the ingress setup and service mesh rules.

**Features**
- Include pod metainfo
- Decode JWT from Authorization header
- Expose Prometheus metrics
- Write structured logs

> **Warning**
> This may expose sensitive data contained in the headers

### How to use
---
Docker:

    docker run --rm -d -p 8080:8080 --name microbin ghcr.io/sshilin/microbin:latest

Kubectl:

    kubectl apply -f microbin.yaml

Helm:

    helm repo add microbin https://sshilin.github.io/microbin

    helm install --generate-name  microbin/microbin

### Endpoints
---
- `GET /headers` - returns request's headers in Json format
```
$ curl http://localhost:8080/headers
{
  "proto": "HTTP/1.1",
  "headers": {
    "Accept": "*/*",
    "User-Agent": "curl/7.77.0"
  },
  "pod": {
    "name": "microbin-58cf7dff8d-7ts8l",
    "namespace": "default",
    "node": "docker-desktop"
  }
}
```
### Environment Variables
---
The table below provides an overview of optional environment variables that can be used to configure microbin.

| Key                 |  Description                | Default         |
|:--------------------|:----------------------------|:----------------|
| `LISTEN`            | Listen on host:port         | 0.0.0.0:8080    |
| `TLS_ENABLED`       | Enable TLS                  | false           |
| `TLS_KEY_FILE`      | TLS key filepath            | ""              |
| `TLS_CERT_FILE`     | TLS cert filepath           | ""              |

### Enable HTTPS
---
This example shows how to enable HTTPS with a self signed certificate.

Create cert:

    openssl req -x509 -nodes -newkey rsa:4096 -keyout ./certs/key.pem -out ./certs/cert.pem -subj "//CN=localhost" -days 365

**Docker**

Run a container with mounted key and cert files:
```bash
$ docker run --rm -d -p 8080:8080 --name microbin \
  -v "$(pwd)/certs:/var/tls:ro" \
  -e TLS_ENABLED=true \
  -e TLS_KEY_FILE=/var/tls/key.pem \
  -e TLS_CERT_FILE=/var/tls/cert.pem \
  microbin:latest
```

**Kubernetes**

First create `microbin-certs` secret with cert and key files, then create `microbin-https.yaml`

    kubectl create secret generic microbin-certs --from-file=./certs

    kubectl apply -f microbin-https.yaml
