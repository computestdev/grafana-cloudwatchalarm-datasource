package main

import (
	"os"

	"github.com/computestdev/grafana-cloudwatch-alarm-datasource/pkg/plugin"
	"github.com/grafana/grafana-aws-sdk/pkg/awsds"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

func main() {
	sessionCache := awsds.NewSessionCache()
	instanceFactory := plugin.NewCloudWatchAlarmDatasource(sessionCache)
	opts := datasource.ManageOpts{}

	if err := datasource.Manage("computest-cloudwatch-alarm-datasource", instanceFactory, opts); err != nil {
		log.DefaultLogger.Error(err.Error())
		os.Exit(1)
	}
}
