import React, { FC } from 'react';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { ConnectionConfig } from '@grafana/aws-sdk';

import { CloudWatchAlarmJsonData, CloudWatchAlarmSecureJsonData } from './types';

interface Props extends DataSourcePluginOptionsEditorProps<CloudWatchAlarmJsonData, CloudWatchAlarmSecureJsonData> {}

export const ConfigEditor: FC<Props> = (props) => {
    const { options, onOptionsChange } = props;

    // todo: add warnings about invalid `authType`, just like:
    //       https://github.com/grafana/grafana/blob/23956557d8c6a119b7de5be5c42024e29634d002/public/app/plugins/datasource/cloudwatch/components/ConfigEditor.tsx

    return <ConnectionConfig options={options} onOptionsChange={onOptionsChange} />;
};
