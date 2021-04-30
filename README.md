intigriti-webhook
========================
Go bot that publishes (non-sensitive) new intigriti findings to a webhook.
This can publish submission notifications to e.g. [a Teams Channel](https://docs.microsoft.com/en-us/microsoftteams/platform/webhooks-and-connectors/how-to/add-incoming-webhook#add-an-incoming-webhook-to-a-teams-channel).

## Setup
1. Download [the latest iwh release](https://github.com/hazcod/intigriti-webhook/releases).
2. Retrieve your [intigriti API token](https://intigriti.com/) and pass your (external) IP address for whitelisting.
4. Create your configuration file:
```yaml
# skip findings in audit, archived and closed
include_non_ready: false

# how often to check in minutes
check_interval_minutes: 15

# desired output format
format: json

# your slack webhook
webhook_url: "https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX"
webhook_headers:
  -  "x-auth": "my-secret"

# your intigriti API credentials
intigriti_client_id: "XXXXXXXXXXX"
intigriti_client_secret: "XXXXXXXXXXX"
```
5. Run `iwh` (preferably as a service) with arguments:
```shell
./iwh -conf=my-conf.yaml
```
3. See new intigriti findings roll in on your Slack channel.
Any findings already sent to your Slack channel will be added to your YAML configuration file for portability.

## Building
This requires `make` and `go` to be installed.
Just run `make`.
