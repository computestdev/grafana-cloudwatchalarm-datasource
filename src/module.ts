import { DataSourcePlugin } from '@grafana/data';
import { DataSource } from './datasource';
import { ConfigEditor } from './ConfigEditor';
import { QueryEditor } from './QueryEditor';
import { CloudWatchAlarmQuery, CloudWatchAlarmJsonData } from './types';

export const plugin = new DataSourcePlugin<DataSource, CloudWatchAlarmQuery, CloudWatchAlarmJsonData>(DataSource)
    .setConfigEditor(ConfigEditor)
    .setQueryEditor(QueryEditor);
