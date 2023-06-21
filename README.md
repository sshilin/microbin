[![Build](https://github.com/sshilin/microbin/actions/workflows/quality.yml/badge.svg)](https://github.com/sshilin/microbin/actions/workflows/quality.yml)&nbsp;[![Go Report Card](https://goreportcard.com/badge/github.com/sshilin/microbin)](https://goreportcard.com/report/github.com/sshilin/microbin)&nbsp;[![Coverage Status](https://coveralls.io/repos/github/sshilin/microbin/badge.svg)](https://coveralls.io/github/sshilin/microbin)

Microbin accepts any HTTP request and returns the request's details including headers, parameters, query strings in JSON format. It is designed to gain insights into how requests are processed and modified as they traverse a cluster environment.

**Running with Docker**
```
docker run --rm --name microbin -p 8080:8080 ghcr.io/sshilin/microbin:latest
```

**Running with kubectl**
```
git clone https://github.com/sshilin/microbin.git && cd microbin
kubectl apply -f k8s/microbin.yaml
```

**Running with Helm**
```
helm repo add microbin https://sshilin.github.io/microbin
helm install --generate-name  microbin/microbin
```

## Usage

> **Note:**
> `GET /metrics` is reserved for Prometheus metrics

HTTP/1.1:
```
$ curl http://localhost:8080/microbin/headers
```

HTTP/2:
```
$ curl --http2 --cacert ./cert/cert.pem https://localhost/microbin/headers
```

HTTP/2 Cleartext (H2C):
```
$ curl --http2-prior-knowledge http://localhost:8080/headers
```

Example response:
```
{
  "host": "microbin-6865c46f94-gpbsj",
  "remote": "10.1.0.1:51346",
  "proto": "HTTP/1.1",
  "method": "GET",
  "uri": "/headers",
  "headers": {
    "Accept": "*/*",
    "Accept-Encoding": "gzip",
    "User-Agent": "curl/7.87.0",
    "X-Forwarded-For": "127.0.0.1, 127.0.0.1",
    "X-Forwarded-Uri": "/api/v1/namespaces/default/services/microbin:8080/proxy/headers"
  }
}
```

## Configuration

| Key                 |  Description                | Default         |
|:--------------------|:----------------------------|:----------------|
| `LISTEN`            | Listen on host:port         | 0.0.0.0:8080    |
| `TLS_ENABLED`       | Enable TLS                  | false           |
| `TLS_KEY_FILE`      | TLS key filepath            | ""              |
| `TLS_CERT_FILE`     | TLS cert filepath           | ""              |
| `LOG_FORMAT_JSON`   | Enable structured logs      | false           |

## HTTPS Examples

```
docker run --name microbin -p 8080:8080 \
  -v $(pwd)/cert:/var/cert:ro \
  -e TLS_ENABLED=true \
  -e TLS_KEY_FILE=/var/cert/key.pem \
  -e TLS_CERT_FILE=/var/cert/cert.pem \
  ghcr.io/sshilin/microbin:latest
```

```
kubectl create secret generic microbin-cert --from-file=./cert
kubectl apply -f k8s/microbin-https.yaml
```
