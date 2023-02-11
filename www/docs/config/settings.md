# Settings

This allows you to override specific settings used in the deployer. For the most part, you should never need to modify these.


## Example

```yaml
settings:
  predeploy_timeout: 90m
  skip_deployment_env_vars: false

  ssm_import: /path/to/project/secrets
```

## Fields

[`ssm_import`](#settings.ssm_import){ #settings.ssm_import }

:   See [SSM Parameter Store Secrets Import](#ssm-import)

    **Highly recommended you enable this!**

    _Default_: disabled (see [SSM Parameter Store Secrets Import](#ssm-import) for defaults)

<!--
[`predeploy_parallel`](#settings.predeploy_parallel){ #settings.predeploy_parallel }

:   If true then [predeploy tasks](predeploy.md) will be run in parallel.
    By default, they are run sequentially. 
    If none of your tasks deploy on each other, you can speed up deployments by enabling this.

    _Default_: `{{schema:default:Settings.predeploy_parallel}}`
-->

[`predeploy_timeout`](#settings.predeploy_timeout){ #settings.predeploy_timeout }

:   The maximum time that a [predeploy task](predeploy.md) may take to run.
    This is the default time, and can be overridden on a per-task basis.

    _Default_: `{{schema:default:Settings.predeploy_timeout}}` (seconds)

[`skip_cron_env_vars`](#settings.skip_cron_env_vars){ #settings.skip_cron_env_vars }

:   This will prevent the extra environment variables from being added to your cron jobs.

    For a list of env vars, see [Cron Env Vars](#cron-env-vars).

    _Default_: `{{schema:default:Settings.skip_cron_env_vars}}`

[`skip_deployment_env_vars`](#settings.skip_deployment_env_vars){ #settings.skip_deployment_env_vars }

:   This will prevent the extra deployment environment variables from being added to your application.

    These are things like `ECSDEPLOYER_IMAGE_TAG`, `ECSDEPLOYER_PROJECT`, etc.

    For a list of env vars, see [Deployment Env Vars](#deployment-env-vars).

    _Default_: `{{schema:default:Settings.skip_deployment_env_vars}}`


[`disable_marker_tag`](#settings.disable_marker_tag){ #settings.disable_marker_tag }

:   This will disable the creation of the marker tag used by the deployer to track resources it creates.

    Disabling the marker tag also requires you to set [`keep_in_sync`](#settings.keep_in_sync) to `false`

    **It's recommended you do not disable this**

    _Default_: `false` (Marker tag will be used)

[`use_old_cron_eventbus`](#settings.use_old_cron_eventbus){ #settings.use_old_cron_eventbus }

:   Use the old and deprecated method of creating CronJobs. The new method uses EventBridge Scheduler, which is much better.

    This option is only really for legacy deployments that still use EventBridge targets/rules. You should not use this.

    _Default_: `false` (Will use the newer EventBridge Scheduler)

[`keep_in_sync`](#settings.keep_in_sync){ #settings.keep_in_sync }

:   See [Keeping Resources In-Sync](#keeping-in-sync)

[`wait_for_stable`](#settings.wait_for_stable){ #settings.wait_for_stable }

:   See [Service Stability Waiter](#service-stability-waiter)

----

## SSM Parameter Store Secrets Import {#ssm-import}

By default, this feature is **disabled**, but is highly recommended.

It will automatically look for environment variables in a specific path on [AWS Parameter Store](https://docs.aws.amazon.com/systems-manager/latest/userguide/systems-manager-parameter-store.html).

It's **highly** recommended that you use SSM Parameter Store to put sensitive values, rather than in plaintext.

=== "Shorthand Enabling"

    ```yaml
    settings:
      ssm_import: /path/to/project/parameters
    ```

    This will enable SSM Secrets Importing, and will set the path as provided.


=== "Advanced Configuration"

    ```yaml
    settings:
      ssm_import:
        enabled: true
        path: /path/to/project/parameters
    ```

    [`enabled`](#ssm_import.enabled){ #ssm_import.enabled }

    :   Enable importing secrets from SSM Parameter

        _Default_: `{{schema:default:SSMImport.enabled}}`

    [`path`](#ssm_import.path){ #ssm_import.path }

    :   This is the path prefix that will be searched to automatically add environment variables to your applications.
        If you do not want SSM Secrets, then set this to a blank string

        === "Default"
            ```
            {{schema:default:TplDefault.ssm_import__path}}
            ```
        
        === "Default (with Stage)"
            ```
            {{schema:default:TplDefaultStage.ssm_import__path}}
            ```

        === "Raw"
            ```
            {{schema:default:SSMImport.path}}
            ```

    [`recursive`](#ssm_import.recursive){ #ssm_import.recursive }

    :   If you have multiple levels of parameters nested underneath the [`path`](#ssm_import.path) provided above, then you will want to enable recursive import. If you do not, then some parameters will be ignored.

        _Default_: `{{schema:default:SSMImport.recursive}}`


----

## Keeping Resources In-Sync {#keeping-in-sync}

!!! tip "Important"
    This feature is reliant upon the [Marker Tag](naming.md#naming.marker_tag_key). If you modify or change the tag or value, then any resources using that tag will be "invisible" to the Deployer.

=== "Enable (default)"

    ```yaml
    settings:
      keep_in_sync: true
    ```

=== "Disable (not recommended)"

    ```yaml
    settings:
      keep_in_sync: false
    ```

    !!! warning "Warning: Not recommended"

        If you disable this, then **you** are responsible for deleting unused services, cronjobs, task definitions, etc.

        This means that services that you remove from the environment file will still be running old code after deploy.

=== "Customize"

    ```yaml
    settings:
      keep_in_sync:
        services: true
        cronjobs: false
    ```

    [`services`](#keep_in_sync.services){ #keep_in_sync.services }

    :   Removes services that no longer appear in the config file

    [`cronjobs`](#keep_in_sync.cronjobs){ #keep_in_sync.cronjobs }

    :   Removes cronjobs that no longer appear in the config file

    [`log_retention`](#keep_in_sync.log_retention){ #keep_in_sync.log_retention }

    :   Will update existing log groups to have the correct retention setting

    [`task_definitions`](#keep_in_sync.task_definitions){ #keep_in_sync.task_definitions }

    :   Deregister any task definitions that point to tasks/services/crons that are no longer listed in the config file.

----

## Service Stability Waiter

By default, the deployer will wait for your services to become stable according to ECS. If you do not want to wait, or need to customize that behavior, then you can do that here.

!!! example "Shorthand to disable waiter (not recommended)"

    ```yaml
    settings:
      wait_for_stable: false
    ```


[`disabled`](#wait_for_stable.disabled){ #wait_for_stable.disabled }

:   Disable the waiter. Deployment will succeed immediately after services have been updated. (Regardless of if they actually work / are not in a crash loop).

    _Default_: `{{schema:default:WaitForStable.disabled}}`

[`timeout`](#wait_for_stable.timeout){ #wait_for_stable.timeout }

:   How long should the deployer wait for a service to become stable

    _Default_: `{{schema:default:WaitForStable.timeout}}` (seconds)

<!-- [`individually`](#wait_for_stable.individually){ #wait_for_stable.individually }

:   X -->

## Deployment Env Vars

<div class="tbl-nowrap-key tbl-normal-font" markdown>

Variable Name             | Template
--------------------------|------
`ECSDEPLOYER_APP_VERSION` | `{{schema:default:DeploymentEnvVars.ECSDEPLOYER_APP_VERSION}}`
`ECSDEPLOYER_DEPLOYED_AT` | `{{schema:default:DeploymentEnvVars.ECSDEPLOYER_DEPLOYED_AT}}`
`ECSDEPLOYER_IMAGE_TAG`   | `{{schema:default:DeploymentEnvVars.ECSDEPLOYER_IMAGE_TAG}}`
`ECSDEPLOYER_PROJECT`     | `{{schema:default:DeploymentEnvVars.ECSDEPLOYER_PROJECT}}`
`ECSDEPLOYER_STAGE`       | `{{schema:default:DeploymentEnvVars.ECSDEPLOYER_STAGE}}`
`ECSDEPLOYER_TASK_NAME`   | `{{schema:default:DeploymentEnvVars.ECSDEPLOYER_TASK_NAME}}`

</div>

!!! note ""
    Any values that evaluate to a blank string will not be added to your environment.


## Cron Env Vars

!!! note ""
    These are only added to cron jobs. Values are [EventBridge Scheduler context attributes](https://docs.aws.amazon.com/scheduler/latest/UserGuide/managing-schedule-context-attributes.html)

<div class="tbl-nowrap-key tbl-normal-font" markdown>

Variable Name                     | Template
----------------------------------|------
`ECSDEPLOYER_CRON_SCHEDULE_ARN`   | `<aws.scheduler.schedule-arn>`
`ECSDEPLOYER_CRON_SCHEDULED_TIME` | `<aws.scheduler.scheduled-time>`
`ECSDEPLOYER_CRON_EXECUTION_ID`   | `<aws.scheduler.execution-id>`
`ECSDEPLOYER_CRON_ATTEMPT`        | `<aws.scheduler.attempt-number>`

</div>