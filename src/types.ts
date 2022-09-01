import { DataQuery } from '@grafana/data';
import { AwsAuthDataSourceSecureJsonData, AwsAuthDataSourceJsonData } from '@grafana/aws-sdk';

export interface CloudWatchAlarmQuery extends DataQuery {
    region?: string;
    includeOk?: boolean;
    includeAlarm?: boolean;
    includeInsufficientData?: boolean;
    alarmNamePrefix?: string;
}

export const defaultQuery: Partial<CloudWatchAlarmQuery> = {
    region: 'default',
    includeOk: true,
    includeAlarm: true,
    includeInsufficientData: true,
    alarmNamePrefix: '',
};

export interface CloudWatchAlarmJsonData extends AwsAuthDataSourceJsonData {}

export interface CloudWatchAlarmSecureJsonData extends AwsAuthDataSourceSecureJsonData {}
