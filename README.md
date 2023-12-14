# dp-population-types-api

Exposes representation of the population types resources which support selection of datasets exposed elsewhere

## Getting started

* Run `make debug`

## Dependencies

* No further dependencies other than those defined in `go.mod`

## Configuration

| Environment variable         | Default                 | Description
| ---------------------------- | ----------------------- | -----------
| BIND_ADDR                    | :27300                  | The host and port to bind to
| GRACEFUL_SHUTDOWN_TIMEOUT    | 5s                      | The graceful shutdown timeout in seconds (`time.Duration` format)
| HEALTHCHECK_INTERVAL         | 30s                     | Time between self-healthchecks (`time.Duration` format)
| HEALTHCHECK_CRITICAL_TIMEOUT | 90s                     | Time to wait until an unhealthy dependent propagates its state to make this app unhealthy (`time.Duration` format)
| OTEL_EXPORTER_OTLP_ENDPOINT  | localhost:4317          | Endpoint for OpenTelemetry service
| OTEL_SERVICE_NAME            | dp-population-types-api | Label of service for OpenTelemetry service

## Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

## License

Copyright © 2022, Office for National Statistics (<https://www.ons.gov.uk>)

Released under MIT license, see [LICENSE](LICENSE.md) for details.
