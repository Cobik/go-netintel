# go-netintel (skeleton)

Сбор сетевых сигналов → Kafka → ClickHouse. Метрики Prometheus, health-check.

## Быстрый старт
```bash
cd deploy
docker compose up -d --build
curl 'http://localhost:8080/healthz'
curl 'http://localhost:8080/v1/collect?domain=example.com'
