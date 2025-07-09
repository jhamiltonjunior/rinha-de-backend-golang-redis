```bash
docker-compose up -d influxdb grafana
```

```bash
docker-compose -f ./k6-dashboard/docker-compose.yml run --publish 5665:5665 --rm k6 run /scripts/script.js --out influxdb=http://influxdb:8086/k6
```

```bash
docker-compose run --rm k6 run --publish 5665:5665 /scripts/script.js --out influxdb=http://influxdb:8086/k6
```


*open dashboard in http://localhost:5665*

