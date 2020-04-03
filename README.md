Prometheus Exporter for Folding@home Donor Stats
================================================

This is an exporter that exposes information gathered from the
[Folding@home stats JSON API][statsapi] for use by the Prometheus monitoring
system.

Prometheus Configuration
------------------------

The fahstats exporter needs to be passed the user as a parameter. This can be
done with relabeling.

Example config:
```YAML
scrape_configs:
  - job_name: 'fahstats'
    scrape_interval: 30m  # Be kind to the Folding@home servers...
    scrape_timeout: 120s  # ...and be patient.
    static_configs:
      - targets:
        - corhere  # Folding@home donor user name
    relabel_configs:
      - source_labels: [__address__]
        target_label: __param_target
      - target_label: __address__
        replacement: 127.0.0.1:9702  # the fahstats exporter's real hostname:port
```

[statsapi]: https://stats.foldingathome.org/api
