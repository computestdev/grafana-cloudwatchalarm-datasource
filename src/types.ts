import { DataQuery } from '@grafana/data';
import { AwsAuthDataSourceSecureJsonData, AwsAuthDataSourceJsonData } from '@grafana/aws-sdk';

export interface CloudWatchAlarmQuery extends DataQuery {
    region?: string;
    includeTypeMetric?: boolean;
    includeTypeComposite?: boolean;
    includeOk?: boolean;
    includeAlarm?: boolean;
    includeInsufficientData?: boolean;
    alarmNamePrefix?: string;
}

export const TEMPLATED_QUERY_KEYS: Array<keyof CloudWatchAlarmQuery> = ['region', 'alarmNamePrefix'];

export const defaultQuery: Partial<CloudWatchAlarmQuery> = {
    region: 'default',
    includeTypeMetric: true,
    includeTypeComposite: true,
    includeOk: true,
    includeAlarm: true,
    includeInsufficientData: true,
    alarmNamePrefix: '',
};

export interface CloudWatchAlarmJsonData extends AwsAuthDataSourceJsonData {}

export interface CloudWatchAlarmSecureJsonData extends AwsAuthDataSourceSecureJsonData {}
