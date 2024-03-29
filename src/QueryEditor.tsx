import { defaults } from 'lodash';
import React, { ChangeEvent, PureComponent, SyntheticEvent } from 'react';
import { InlineField, Input, Switch, Select } from '@grafana/ui';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { standardRegions } from '@grafana/aws-sdk';

import { DataSource } from './datasource';
import { defaultQuery, CloudWatchAlarmJsonData, CloudWatchAlarmQuery } from './types';

type Props = QueryEditorProps<DataSource, CloudWatchAlarmQuery, CloudWatchAlarmJsonData>;

const regionOptions = [
    {
        value: 'default',
        label: '',
    },
    ...standardRegions.map((region) => ({
        value: region,
        label: region,
    })),
];

const alarmTypeOptions = [
    {
        value: 'all',
        label: 'All',
    },
    {
        value: 'metric',
        label: 'Metric',
    },
    {
        value: 'composite',
        label: 'Composite',
    },
];

export class QueryEditor extends PureComponent<Props> {
    onRegionChange = (value: SelectableValue<string>) => {
        const { onChange, query, onRunQuery } = this.props;
        onChange({ ...query, region: value.value });
        onRunQuery();
    };

    onAlarmTypeChange = ({value}: SelectableValue<string>) => {
        const { onChange, query, onRunQuery } = this.props;

        onChange({
            ...query,
            includeTypeMetric: value === 'all' || value === 'metric',
            includeTypeComposite: value === 'all' || value === 'composite',
        });
        onRunQuery();
    };

    onAlarmNamePrefixChange = (event: ChangeEvent<HTMLInputElement>) => {
        const { onChange, query, onRunQuery } = this.props;
        onChange({ ...query, alarmNamePrefix: event.target.value });
        onRunQuery();
    };

    onIncludeOkChange = (event: SyntheticEvent<HTMLInputElement>) => {
        const { onChange, query, onRunQuery } = this.props;
        onChange({ ...query, includeOk: event.currentTarget.checked });
        onRunQuery();
    };

    onIncludeAlarmChange = (event: SyntheticEvent<HTMLInputElement>) => {
        const { onChange, query, onRunQuery } = this.props;
        onChange({ ...query, includeAlarm: event.currentTarget.checked });
        onRunQuery();
    };

    onIncludeInsufficientDataChange = (event: SyntheticEvent<HTMLInputElement>) => {
        const { onChange, query, onRunQuery } = this.props;
        onChange({ ...query, includeInsufficientData: event.currentTarget.checked });
        onRunQuery();
    };

    render() {
        const query = defaults(this.props.query, defaultQuery);
        const { region, includeTypeMetric, includeTypeComposite, includeOk, includeAlarm, includeInsufficientData, alarmNamePrefix } = query;
        let alarmType = 'all';
        if (includeTypeMetric && !includeTypeComposite) {
            alarmType = 'metric';
        }
        else if (!includeTypeMetric && includeTypeComposite) {
            alarmType = 'composite';
        }

        return (
            <div className="gf-form">
                <div>
                    <InlineField
                        label="Region"
                        labelWidth={20}
                        tooltip="Specify the region, such as for US West (Oregon) use ` us-west-2 ` as the region."
                    >
                        <Select
                            aria-label="Region"
                            className="width-20"
                            value={region}
                            options={regionOptions}
                            defaultValue="default"
                            allowCustomValue={true}
                            onChange={this.onRegionChange}
                            formatCreateLabel={(r) => `Use region: ${r}`}
                            menuShouldPortal={true}
                        />
                    </InlineField>

                    <InlineField label="Alarm name prefix" labelWidth={20}>
                        <Input
                            aria-label="Alarm name prefix"
                            className="width-20"
                            value={alarmNamePrefix || ''}
                            onChange={this.onAlarmNamePrefixChange}
                        />
                    </InlineField>

                    <InlineField
                        label="Alarm Type"
                        labelWidth={20}
                        tooltip="Specify which type of alarms to display (metric and/or composite)."
                    >
                        <Select
                            aria-label="Alarm Type"
                            className="width-20"
                            value={alarmType}
                            options={alarmTypeOptions}
                            defaultValue="all"
                            allowCustomValue={false}
                            onChange={this.onAlarmTypeChange}
                            menuShouldPortal={true}
                        />
                    </InlineField>
                </div>
                <div>
                    <InlineField label="Include OK" labelWidth={24}>
                        <Switch value={includeOk || false} aria-label="Include OK" onChange={this.onIncludeOkChange} />
                    </InlineField>

                    <InlineField label="Include ALARM" labelWidth={24}>
                        <Switch
                            value={includeAlarm || false}
                            aria-label="Include ALARM"
                            onChange={this.onIncludeAlarmChange}
                        />
                    </InlineField>

                    <InlineField label="Include INSUFFICIENT_DATA" labelWidth={24}>
                        <Switch
                            value={includeInsufficientData || false}
                            aria-label="Include INSUFFICIENT_DATA"
                            onChange={this.onIncludeInsufficientDataChange}
                        />
                    </InlineField>
                </div>
            </div>
        );
    }
}
