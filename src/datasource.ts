import { DataSourceInstanceSettings, ScopedVars } from '@grafana/data';
import { DataSourceWithBackend, getTemplateSrv } from '@grafana/runtime';

import { CloudWatchAlarmJsonData, CloudWatchAlarmQuery, TEMPLATED_QUERY_KEYS } from './types';


export class DataSource extends DataSourceWithBackend<CloudWatchAlarmQuery, CloudWatchAlarmJsonData> {
    constructor(instanceSettings: DataSourceInstanceSettings<CloudWatchAlarmJsonData>) {
        super(instanceSettings);
    }

    applyTemplateVariables(query: CloudWatchAlarmQuery, scopedVars: ScopedVars): CloudWatchAlarmQuery {
        const templateSrv = getTemplateSrv();
        const result = {...query};

        for (const key of TEMPLATED_QUERY_KEYS) {
            const value = result[key];
            if (typeof value === 'string') {
                (result as any)[key] = templateSrv.replace(value, scopedVars);
            }
        }

        return result;
    }
}
