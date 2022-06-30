# Grafana CloudWatch Alarm data source
A data source plugin for [Grafana](https://grafana.com/) which queries for the current state of [Amazon Web Services CloudWatch](https://aws.amazon.com/cloudwatch/) alarms.

This Grafana plugin has been designed to help us to easily create a single dashboard that gives us a quick overview of all alerts/alarms from different monitoring solutions (ex: Grafana and CloudWatch).

Note that it only queries CloudWatch for the current state of alarms, it does not query for the history. As such, this data source is best used with the "table" visualization. 

If you would like to create graphs based on the alarm history, it would be better to configure AWS to store the alarm state as custom metrics (the regular alarm history only goes back 14 days). For example by following this guide: [Extending and exploring alarm history in Amazon CloudWatch](https://aws.amazon.com/blogs/mt/extending-and-exploring-alarm-history-in-amazon-cloudwatch-part-1/) and then using the regular CloudWatch data source that ships with Grafana.

## Building this plugin

### Frontend

1. Install dependencies

```sh
npm install
```

2. Build plugin in development mode

```bash
npm run dev
```

3. Build plugin in production mode

```bash
npm run build
```

### Backend

1. Update [Grafana plugin SDK for Go](https://grafana.com/docs/grafana/latest/developers/plugins/backend/grafana-plugin-sdk-for-go/) dependency to the latest minor version:

```bash
go get -u github.com/grafana/grafana-plugin-sdk-go
go mod tidy
```

2. Build backend plugin binaries for Linux, Windows and Darwin:

```bash
mage -v
```

3. List all available Mage targets for additional commands:

```bash
mage -l
```
