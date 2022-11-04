# CronJobs

You can easily create [Scheduled Tasks](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/scheduled_tasks.html) for your application by specifying a cronjobs block.

## Example
```yaml
cronjobs:
  - name: reaper
    schedule: "cron(0 9 * * ? *)"
    command: ["bundle", "exec", "rake", "cron:whatever"]
```


## Fields

In addition to the fields listed below, you can also specify anything in [Common Task Options](common.md).

[`name`](#cronjob.name){ #cronjob.name } - **(required)**

:   Unique name for your cron job. This will be used to create the Rule and Target on EventBridge. This should be a short identifier with only letters, numbers, dash or underscore.

[`schedule`](#cronjob.schedule){ #cronjob.schedule } - **(required)**

:   The schedule expression that is used to determine when your task runs. This can be either a [cron expression](https://docs.aws.amazon.com/eventbridge/latest/userguide/eb-create-rule-schedule.html#eb-cron-expressions), or a [rate expression](https://docs.aws.amazon.com/eventbridge/latest/userguide/eb-create-rule-schedule.html#eb-rate-expressions).

    Expression format: [EventBridge Cron/Rate Expressions](https://docs.aws.amazon.com/eventbridge/latest/userguide/eb-create-rule-schedule.html)

[`command`](#cronjob.command){ #cronjob.command } - _(recommended)_

:   The command you want to run for this cronjob.

    This field is not required, but it doesn't make much sense to exclude it.
    For details on how to specify a command, see [Specifying Command/EntryPoints](common.md#commandentry-point)

[`description`](#cronjob.description){ #cronjob.description }

:   An optional description that will be added to the Eventbridge Rule. Any string is valid here.

    _Default_: _none_

[`disabled`](#cronjob.disabled){ #cronjob.disabled }

:   Create this cronjob, but do not enable it. It will not run.

    _Default_: `{{schema:default:CronJob.disabled}}`

[`<anything from common>`](#cronjob.common){ #cronjob.common }

:   See [Common Task Options](common.md).

## See Also

* [Customizing Resource Names](naming.md)
* [Templating](../templating.md)
* [AWS Scheduled Tasks](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/scheduled_tasks.html)