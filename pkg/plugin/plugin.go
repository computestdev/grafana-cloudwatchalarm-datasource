package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/grafana/grafana-aws-sdk/pkg/awsds"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"time"
)

// Make sure CloudWatchAlarmDatasource implements required interfaces.
var (
	_ backend.QueryDataHandler      = (*CloudWatchAlarmDatasource)(nil)
	_ backend.CheckHealthHandler    = (*CloudWatchAlarmDatasource)(nil)
	_ instancemgmt.InstanceDisposer = (*CloudWatchAlarmDatasource)(nil)
)

var plog = log.DefaultLogger

type dataSourceJsonModel struct {
	AuthType      string `json:"authType"`
	AssumeRoleARN string `json:"assumeRoleArn"`
	ExternalID    string `json:"externalId"`
	Profile       string `json:"profile"`
	Region        string `json:"defaultRegion"`
	Endpoint      string `json:"endpoint"`
}

type queryModel struct {
	Region                  string `json:"region"`
	IncludeTypeMetric       bool   `json:"includeTypeMetric"`
	IncludeTypeComposite    bool   `json:"includeTypeComposite"`
	IncludeOk               bool   `json:"includeOk"`
	IncludeAlarm            bool   `json:"includeAlarm"`
	IncludeInsufficientData bool   `json:"includeInsufficientData"`
	AlarmNamePrefix         string `json:"alarmNamePrefix"`
}

func defaultQuery() queryModel {
	return queryModel{
		Region:                  "default",
		IncludeTypeComposite:    true,
		IncludeTypeMetric:       true,
		IncludeOk:               false,
		IncludeAlarm:            true,
		IncludeInsufficientData: true,
		AlarmNamePrefix:         "",
	}
}

func (query *queryModel) stateValues() []string {
	stateValues := make([]string, 0)
	if query.IncludeAlarm {
		stateValues = append(stateValues, cloudwatch.StateValueAlarm)
	}
	if query.IncludeInsufficientData {
		stateValues = append(stateValues, cloudwatch.StateValueInsufficientData)
	}
	if query.IncludeOk {
		stateValues = append(stateValues, cloudwatch.StateValueOk)
	}
	return stateValues
}

func (query *queryModel) alarmTypes() []*string {
	typeValues := make([]*string, 0)
	if query.IncludeTypeMetric {
		value := cloudwatch.AlarmTypeMetricAlarm
		typeValues = append(typeValues, &value)
	}
	if query.IncludeTypeComposite {
		value := cloudwatch.AlarmTypeCompositeAlarm
		typeValues = append(typeValues, &value)
	}
	return typeValues
}

type CloudWatchAlarmDatasource struct {
	authType      awsds.AuthType
	assumeRoleARN string
	externalID    string
	endpoint      string
	profile       string
	region        string

	accessKey string
	secretKey string

	sessions *awsds.SessionCache
}

func parseAuthType(authType string) awsds.AuthType {
	switch authType {
	case "credentials":
		return awsds.AuthTypeSharedCreds
	case "keys":
		return awsds.AuthTypeKeys
	case "default":
		return awsds.AuthTypeDefault
	case "ec2_iam_role":
		return awsds.AuthTypeEC2IAMRole
	case "arn":
		plog.Warn("Authentication type \"arn\" is deprecated, falling back to default")
		return awsds.AuthTypeDefault
	default:
		plog.Warn("Unrecognized AWS authentication type", "type", authType)
		return awsds.AuthTypeDefault
	}
}

func NewCloudWatchAlarmDatasource(sessions *awsds.SessionCache) datasource.InstanceFactoryFunc {
	return func(ctx context.Context, settings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
		jsonData := dataSourceJsonModel{}
		err := json.Unmarshal(settings.JSONData, &jsonData)
		if err != nil {
			return nil, fmt.Errorf("error reading settings: %w", err)
		}

		d := CloudWatchAlarmDatasource{
			authType:      parseAuthType(jsonData.AuthType),
			assumeRoleARN: jsonData.AssumeRoleARN,
			externalID:    jsonData.ExternalID,
			endpoint:      jsonData.Endpoint,
			profile:       jsonData.Profile,
			region:        jsonData.Region,

			accessKey: settings.DecryptedSecureJSONData["accessKey"],
			secretKey: settings.DecryptedSecureJSONData["secretKey"],

			sessions: sessions,
		}

		return &d, nil
	}
}

func (d *CloudWatchAlarmDatasource) newSession(region string) (*session.Session, error) {
	if region == "default" {
		region = d.region
	}

	sessionConfig := awsds.SessionConfig{
		Settings: awsds.AWSDatasourceSettings{
			Profile:       d.profile,
			Region:        region,
			AuthType:      d.authType,
			AssumeRoleARN: d.assumeRoleARN,
			ExternalID:    d.externalID,
			Endpoint:      d.endpoint,
			DefaultRegion: d.region,
			AccessKey:     d.accessKey,
			SecretKey:     d.secretKey,
		},
		UserAgentName: aws.String("Cloudwatch"),
	}

	return d.sessions.GetSession(sessionConfig)
}

func (d *CloudWatchAlarmDatasource) Dispose() {
	// Clean up datasource instance resources.
}

func (d *CloudWatchAlarmDatasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	response := backend.NewQueryDataResponse()

	for _, q := range req.Queries {
		res := d.query(ctx, req.PluginContext, q)
		response.Responses[q.RefID] = res
	}

	return response, nil
}

func (d *CloudWatchAlarmDatasource) query(_ context.Context, _ backend.PluginContext, dataQuery backend.DataQuery) backend.DataResponse {
	response := backend.DataResponse{}

	query := defaultQuery()
	err := json.Unmarshal(dataQuery.JSON, &query)
	if err != nil {
		response.Error = fmt.Errorf("%v: %w", "failed to unmarshal DataQuery", err)
		return response
	}

	awsSession, err := d.newSession(query.Region)
	if err != nil {
		response.Error = fmt.Errorf("%v: %w", "failed to create AWS SDK session", err)
		return response
	}

	cloudwatchApi := cloudwatch.New(awsSession)

	alarmName := make([]*string, 0)
	alarmDescription := make([]*string, 0)
	stateUpdatedTimestamp := make([]*time.Time, 0)
	stateValue := make([]*string, 0)
	stateReason := make([]*string, 0)

	for _, queryStateValue := range query.stateValues() {
		alarmsParams := cloudwatch.DescribeAlarmsInput{
			StateValue: &queryStateValue,
			AlarmTypes: query.alarmTypes(),
		}

		if query.AlarmNamePrefix != "" {
			alarmsParams.AlarmNamePrefix = &query.AlarmNamePrefix
		}

		for {
			alarmsResponse, err := cloudwatchApi.DescribeAlarms(&alarmsParams)
			if err != nil {
				response.Error = fmt.Errorf("%v: %w", "failed to call cloudwatch:DescribeAlarms", err)
				return response
			}

			for _, metricAlarm := range alarmsResponse.CompositeAlarms {
				alarmName = append(alarmName, metricAlarm.AlarmName)
				alarmDescription = append(alarmDescription, metricAlarm.AlarmDescription)
				stateUpdatedTimestamp = append(stateUpdatedTimestamp, metricAlarm.StateUpdatedTimestamp)
				stateValue = append(stateValue, metricAlarm.StateValue)
				stateReason = append(stateReason, metricAlarm.StateReason)
			}

			for _, metricAlarm := range alarmsResponse.MetricAlarms {
				alarmName = append(alarmName, metricAlarm.AlarmName)
				alarmDescription = append(alarmDescription, metricAlarm.AlarmDescription)
				stateUpdatedTimestamp = append(stateUpdatedTimestamp, metricAlarm.StateUpdatedTimestamp)
				stateValue = append(stateValue, metricAlarm.StateValue)
				stateReason = append(stateReason, metricAlarm.StateReason)
			}

			if alarmsResponse.NextToken != nil {
				alarmsParams.NextToken = alarmsResponse.NextToken
			} else {
				break
			}
		}
	}

	frame := data.NewFrame("response")
	frame.Meta = &data.FrameMeta{
		PreferredVisualization: data.VisTypeTable,
	}

	frame.Fields = append(frame.Fields,
		data.NewField("AlarmName", nil, alarmName),
		data.NewField("AlarmDescription", nil, alarmDescription),
		data.NewField("StateUpdatedTimestamp", nil, stateUpdatedTimestamp),
		data.NewField("StateValue", nil, stateValue),
		data.NewField("StateReason", nil, stateReason),
	)

	response.Frames = append(response.Frames, frame)
	return response
}

func (d *CloudWatchAlarmDatasource) CheckHealth(_ context.Context, _ *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	var status = backend.HealthStatusOk
	var message = "Data source is working"

	awsSession, err := d.newSession("default")
	if err != nil {
		status = backend.HealthStatusError
		message = err.Error()

	} else {
		cloudwatchApi := cloudwatch.New(awsSession)

		maxRecords := int64(1)
		_, err := cloudwatchApi.DescribeAlarms(&cloudwatch.DescribeAlarmsInput{
			MaxRecords: &maxRecords,
		})

		if err != nil {
			status = backend.HealthStatusError
			message = err.Error()
		}
	}

	return &backend.CheckHealthResult{
		Status:  status,
		Message: message,
	}, nil
}
