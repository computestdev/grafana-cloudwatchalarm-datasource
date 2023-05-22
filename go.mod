module github.com/computestdev/grafana-cloudwatchalarm-datasource

go 1.16

require (
	github.com/aws/aws-sdk-go v1.44.266
	github.com/grafana/grafana-aws-sdk v0.15.1
	github.com/grafana/grafana-plugin-sdk-go v0.161.0
)

// fix for CVE-2022-35737
replace github.com/mattn/go-sqlite3 v1.14.14 => github.com/mattn/go-sqlite3 v1.14.15
