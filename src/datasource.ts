import { DataSourceInstanceSettings } from '@grafana/data';
import { DataSourceWithBackend } from '@grafana/runtime';
import { CloudWatchAlarmJsonData, CloudWatchAlarmQuery } from './types';

export class DataSource extends DataSourceWithBackend<CloudWatchAlarmQuery, CloudWatchAlarmJsonData> {
  constructor(instanceSettings: DataSourceInstanceSettings<CloudWatchAlarmJsonData>) {
    super(instanceSettings);
  }
}
