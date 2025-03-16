## Smart Maic Prometheus Exporter

### Exporter is Python app available as docker image:

```
jilexandr/smart-maic-exporter:latest
```

configure environment variables:

```
# also change IP your own (required)
DATA_SOURCE_URL=http://192.168.10.55/?page=getdata&devid={DEVICE_ID}&devpass={DEVICE_PASS}

EXPORTER_PORT=8000 # port by default, optional to change
```

### prometheus.yml

```yaml
scrape_configs:
  - job_name: "smart_maic"
    scrape_interval: 6s
    static_configs:
      - targets: ["192.168.10.98:8000"] # ip of exporter
        labels:
          devid: "1828727481" # device id
```

### grafana:

- configure Prometheus data source
- import dashboard from JSON file [grafana-dashboard.json](grafana-dashboard.json)
