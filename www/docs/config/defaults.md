# Task Defaults

The task defaults section lets you define attributes that all [PreDeploy Tasks](predeploy.md), [CronJobs](cronjobs.md), and [Services](services.md) will inherit.

Anything you specify here can be overridden on the individual task as well.

This is where you would set things like CPU/Memory resources, storage, architecture, etc.

## Common Options
You can override **every** option in the [Common Task Options](common.md) section.

## Spot Configuration
This is also where you can specify the [Spot Configuration](spot.md) if you want to utilize Fargate Spot.


## Example

```yaml
task_defaults:
  cpu: 2048
  memory: 4x

  spot:
    enabled: true
```

## Fields

[`<anything from common>`](#defaults.common){ #defaults.common }

:   See [Common Task Options](common.md).

    Any values you specify will be used as the default for all tasks, cronjobs, services, etc. created by the deployer.

[`spot`](#defaults.spot){ #defaults.spot }

:   Allows you to specify the default [Spot Configuration](spot.md) for _services_ (spot does not apply to single run tasks)


### Default Values

To make onboarding easier, ECS Deployer provides a few default values for you out-of-the-box. You are welcome to override these if you want.

<div class="tbl-nowrap-key tbl-normal-font" markdown>

Field | Default Value
----|-----------
[`arch`](common.md#common.arch) | `{{schema:default:FargateDefaults.arch}}`
[`cpu`](common.md#common.cpu) | `{{schema:default:FargateDefaults.cpu}}`
[`memory`](common.md#common.memory) | `{{schema:default:FargateDefaults.memory}}`
[`platform_version`](common.md#common.platform_version) | `{{schema:default:FargateDefaults.platform_version}}`

</div>


