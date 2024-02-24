# plan-dhcp-grafana

Use the PLAN DHCP lease API to fetch participants data, then update
Grafana dashboard with queries to gather data about bandwidth usage on
participant ports.

## Requirements

- Go
- Grafana token set in env var `GRAFANA_TOKEN`

## Usage

`go run .`

or

`go build . && ./plan-dhcp-grafana`