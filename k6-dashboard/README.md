```bash
docker-compose up -d influxdb grafana
```


```bash
docker-compose run --rm k6 run /scripts/script.js --out influxdb=http://influxdb:8086/k6
```
