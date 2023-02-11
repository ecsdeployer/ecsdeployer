# CronJobs

You can easily create [Scheduled Tasks](https://docs.aws.amazon.com/scheduler/latest/UserGuide/what-is-scheduler.html) for your application by specifying a cronjobs block.

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

:   Unique name for your cron job. This will be used to create the schedule on EventBridge Scheduler. This should be a short identifier with only letters, numbers, dash or underscore.

[`schedule`](#cronjob.schedule){ #cronjob.schedule } - **(required)**

:   The schedule expression that is used to determine when your task runs. This can be either a [cron expression](https://docs.aws.amazon.com/scheduler/latest/UserGuide/schedule-types.html#cron-based), a [rate expression](https://docs.aws.amazon.com/scheduler/latest/UserGuide/schedule-types.html#rate-based), or a [one-time expression](https://docs.aws.amazon.com/scheduler/latest/UserGuide/schedule-types.html#one-time).

    Expression format: [EventBridge Cron/Rate Expressions](https://docs.aws.amazon.com/scheduler/latest/UserGuide/schedule-types.html)

[`command`](#cronjob.command){ #cronjob.command } - _(recommended)_

:   The command you want to run for this cronjob.

    This field is not required, but it doesn't make much sense to exclude it.
    For details on how to specify a command, see [Specifying Command/EntryPoints](common.md#commandentry-point)

[`description`](#cronjob.description){ #cronjob.description }

:   An optional description that will be added to the schedule. Any string is valid here.

    _Default_: _none_

[`disabled`](#cronjob.disabled){ #cronjob.disabled }

:   Create this cronjob, but do not enable it. It will not run.

    _Default_: `{{schema:default:CronJob.disabled}}`


[`timezone`](#cronjob.timezone){ #cronjob.timezone }

:   Sets the time zone used when interpreting the schedule expression.

    This should be an [IANA Timezone Identifier](https://www.iana.org/time-zones) like `UTC` or `America/Los_Angeles`.

    _Default_: _none_ (AWS will default to UTC)

[`start_date`](#cronjob.start_date){ #cronjob.start_date }

:   You can set a cron to start evaluating after a specific date if needed.

    If provided, must be in [RFC 3339](https://www.rfc-editor.org/rfc/rfc3339#section-5.8) format.

    _Default_: _none_

[`end_date`](#cronjob.end_date){ #cronjob.end_date }

:   Stop executing the cron job after the given date.

    If provided, must be in [RFC 3339](https://www.rfc-editor.org/rfc/rfc3339#section-5.8) format.

    _Default_: _none_

[`<anything from common>`](#cronjob.common){ #cronjob.common }

:   See [Common Task Options](common.md).

## See Also

* [Customizing Resource Names](naming.md)
* [Templating](../templating.md)
* [AWS EventBridge Scheduler](https://docs.aws.amazon.com/scheduler/latest/UserGuide/what-is-scheduler.html)