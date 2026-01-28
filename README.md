# Stress Test Full Cycle

CLI em Go para testes de carga HTTP.

## Uso local

```bash
go run ./cmd --url=http://localhost:8080 --requests=100 --concurrency=10
```

## Docker

Build:

```bash
docker build -t stress-test .
```

Run:

```bash
docker run --rm stress-test --url=http://google.com --requests=1000 --concurrency=10
```

## Testes

```bash
go test ./...
```