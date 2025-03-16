- docker image

```
jilexandr/smart-maic-exporter:latest
```

- prometheus.yml

```yaml
scrape_configs:
  - job_name: "smart_maic"
    scrape_interval: 6s
    static_configs:
      - targets: ["192.168.10.98:8000"]
        labels:
          devid: "1828727481"
```

- grafana dashboard
import `grafana-dashboard.json`
