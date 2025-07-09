```bash
docker-compose up -d influxdb grafana
```

```bash
docker-compose -f ./k6-dashboard/docker-compose.yml run --rm k6 run /scripts/script.js --out influxdb=http://influxdb:8086/k6
```

```bash
docker-compose run --rm k6 run /scripts/script.js --out influxdb=http://influxdb:8086/k6
```
