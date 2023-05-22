# Grafana CloudWatch Alarm data source
A data source plugin for [Grafana](https://grafana.com/) which queries for the current state of [Amazon Web Services CloudWatch](https://aws.amazon.com/cloudwatch/) alarms.

This Grafana plugin has been designed to help us to easily create a single dashboard that gives us a quick overview of all alerts/alarms from different monitoring solutions (ex: Grafana and CloudWatch).

Note that it only queries CloudWatch for the current state of alarms, it does not query for the history. As such, this data source is best used with the "table" visualization. 

If you would like to create graphs based on the alarm history, it would be better to configure AWS to store the alarm state as custom metrics (the regular alarm history only goes back 14 days). For example by following this guide: [Extending and exploring alarm history in Amazon CloudWatch](https://aws.amazon.com/blogs/mt/extending-and-exploring-alarm-history-in-amazon-cloudwatch-part-1/) and then using the regular CloudWatch data source that ships with Grafana.
