module github.com/computestdev/grafana-cloudwatchalarm-datasource

go 1.16

require (
	github.com/aws/aws-sdk-go v1.44.266
	github.com/chromedp/cdproto v0.0.0-20230914224007-a15a36ccbc2e // indirect
	github.com/elazarl/goproxy v0.0.0-20230808193330-2592e75ae04a // indirect
	github.com/getkin/kin-openapi v0.120.0 // indirect
	github.com/go-openapi/jsonpointer v0.20.0 // indirect
	github.com/google/uuid v1.3.1 // indirect
	github.com/grafana/grafana-aws-sdk v0.15.1
	github.com/grafana/grafana-plugin-sdk-go v0.178.0
	github.com/magefile/mage v1.15.0 // indirect
	github.com/unknwon/log v0.0.0-20200308114134-929b1006e34a // indirect
	github.com/urfave/cli v1.22.14 // indirect
	golang.org/x/sys v0.12.0 // indirect
)

// fix for CVE-2022-35737
replace github.com/mattn/go-sqlite3 v1.14.14 => github.com/mattn/go-sqlite3 v1.14.15
