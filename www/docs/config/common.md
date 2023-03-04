# Shared Task Attributes

The attributes below are available in the following sections:

* [Services](services.md)
* [Remote Shell](console.md)
* [Task Defaults](defaults.md)
* [PreDeploy Tasks](predeploy.md)
* [CronJobs](cronjobs.md)

<!--
<div class="tbl-nowrap-key" markdown>

Key | Description
----|-----------
`arch`              | Override the architecture of this task. Values are `amd64` or `arm64`.
`command`           | See [Command/Entry Point](#commandentry-point)
`cpu`               | See [CPU/Memory Resources](#cpumemory-resources)
`credentials`       | See [Private Registry](#private-registry)
`entrypoint`        | See [Command/Entry Point](#commandentry-point)
`environment`       | Allows extra task-specific environment variables to be added. For details, see the [Environment Variables](envvars.md) page.
`image`             | Allows you to override the image used for this task. See [Image Documentation](basic.md#specifying-container-image)
`labels`            | Add Docker Labels. Same structure as used for tags below.
`memory`            | See [CPU/Memory Resources](#cpumemory-resources)
`network`           | Override network settings. See [Network](network.md)
`platform_version`  | Override the [Fargate Platform Version](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/platform_versions.html). Default is `LATEST`.
`start_timeout`<br>`stop_timeout` | See [Start/Stop Timeouts](#startstop-timeouts)
`storage`           | See [AWS Documentation on Fargate Storage](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/fargate-task-storage.html)
`tags`              | Additional tags to apply to this task. Same format as in the [Tags Documentation](tags.md)

</div>

`depends_on`        | XXX
`healthcheck`       | XXX
`logging`           | XXX
-->

## Usage

=== "Task Defaults"

    ```yaml
    task_defaults:
      cpu: 1024
      <something from Fields>: value
    ```

=== "Services"

    ```yaml
    task_defaults:
      cpu: 1024

    services:
      - name: web
        <something from Fields>: value

      - name: important-jobs
        cpu: 4096
        memory: 4x
    ```

=== "Cronjobs"

    ```yaml
    cronjobs:
      - name: reaper
        ...
        <something from Fields>: value
    ```

=== "PreDeploy Tasks"

    ```yaml
    predeploy:
      - name: dbmigrate
        ...
        <something from Fields>: value
    ```

=== "Remote Shell"

    ```yaml
    console:
      enabled: true
      ...
      <something from Fields>: value
    ```


## Fields

[`arch`](#common.arch){ #common.arch }

:   Override the architecture of this task. Values are `amd64` or `arm64`.

    _Default_: `{{schema:default:FargateDefaults.arch}}`

[`command`](#common.command){ #common.command }

:   See [Command/Entry Point](#commandentry-point)

[`cpu`](#common.cpu){ #common.cpu }

:   How many CPU shares are given to the task

    See [CPU/Memory Resources](#cpumemory-resources)

    _Default_: `{{schema:default:FargateDefaults.cpu}}`

[`credentials`](#common.credentials){ #common.credentials }

:   See [Private Registry](#private-registry)

[`entrypoint`](#common.entrypoint){ #common.entrypoint }

:   See [Command/Entry Point](#commandentry-point)

[`environment`](#common.environment){ #common.environment }

:   Allows extra task-specific environment variables to be added.

    See [Environment Variables](envvars.md)

[`image`](#common.image){ #common.image }

:   Allows you to override the image used for this task. See [Image Documentation](basic.md#specifying-container-image)

[`labels`](#common.labels){ #common.labels }

:   Add Docker Labels. Same structure as in the [Tags Documentation](tags.md).

<!--
    ??? example "Example: Datadog label configuration"
        ```yaml
        task_defaults:
          labels:
            - name: com.datadoghq.ad.instances
              value: '[{"host": "%%host%%", "port": <PORT_NUMBER>}]'
            - name: com.datadoghq.ad.check_names
              value: '["<CHECK_NAME>"]'
            - name: com.datadoghq.ad.init_configs
              value: "[{}]"
        ```
-->

[`memory`](#common.memory){ #common.memory }

:   See [CPU/Memory Resources](#cpumemory-resources)

    _Default_: `{{schema:default:FargateDefaults.memory}}`

[`mounts`](#common.mounts){ #common.mounts }

:   Specify mount points. See [Volumes/Mounts](volumes.md).

[`network`](#common.network){ #common.network }

:   Override network settings. See [Network](network.md)

[`platform_version`](#common.platform_version){ #common.platform_version }

:   Override the [Fargate Platform Version](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/platform_versions.html).

    _Default_: `{{schema:default:FargateDefaults.platform_version}}`

[`proxy`](#common.proxy){ #common.proxy }

:   See [Proxy Configuration](#proxy-configuration)

[`start_timeout`](#common.start_timeout){ #common.start_timeout }

:   See [Start/Stop Timeouts](#startstop-timeouts)

[`stop_timeout`](#common.stop_timeout){ #common.stop_timeout }

:   See [Start/Stop Timeouts](#startstop-timeouts)

[`storage`](#common.storage){ #common.storage }

:   Amount of storage to attach to the cluster (in GiB)

    See [AWS Documentation on Fargate Storage](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/fargate-task-storage.html)

[`healthcheck`](#common.healthcheck){ #common.healthcheck }

:   See [Container Health Checks](#container-health-checks)

[`tags`](#common.tags){ #common.tags }

:   Additional tags to apply to this task. Same format as in the [Tags Documentation](tags.md).

[`ulimits`](#common.ulimits){ #common.ulimits }

:   Specify limit overrides per container. See [Ulimits](#ulimits) for more details.

[`user`](#common.user){ #common.user }

:   Override the user to run your container as. Specify as username or UID or UID:GID.

[`volumes`](#common.volumes){ #common.volumes }

:   Specify volumes that can be mounted. See [Volumes/Mounts](volumes.md).

[`workdir`](#common.workdir){ #common.workdir }

:   Override the working directory to run your container in.

----

## CPU/Memory Resources

[`cpu`](#resources.cpu){ #resources.cpu }

:   CPU Shares specified as an integer. 1 cpu core is equal to `1024` shares. 

[`memory`](#resources.memory){ #resources.memory }

:   Memory allocation for the task. Can be specified in multiple ways:

    * An integer, denoting megabytes (ex: `4096`)
    * A multiplier of CPU shares in the format of `#.#x` (ex: `2x` or `0.25x`)
    * As gigabytes with the format of `#.# GB` (ex: `8GB` or `0.5 GB`)


!!! warning
    Fargate places restrictions on the allowable CPU/Memory requirements. (See more on [Fargate task sizes](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/AWS_Fargate.html#fargate-tasks-size))

    If you specify CPU/Memory combinations that do not exist on Fargate, the deployer will automatically pick the smallest size that will contain the CPU/Memory requirements you have specified.

    Example: If you specify cpu=512/memory=7000, then the deployer will use cpu=1024,memory=7168 as that is the smallest Fargate size that meets the requirements.

----

## Command/Entry Point

Commands and EntryPoints can be specified in two ways:

=== "Array Based (preferred)"

    ```yaml
    command: ["bundle", "exec", "puma", "-c", "config/puma.rb"]
    ```

=== "String Based"

    ```yaml
    command: "bundle exec puma -c config/puma.rb"
    ```

!!! warning
    Commands are always represented as an array of strings. If you provide it as a string, that string will be split into an array, (possibly incorrectly). For that reason we recommend you only use the **Array based** method for specifying commands.

    **Also**: Environment/Shell interpolation is not available on ECS. Do not use environment references in your commands.

----

## Start/Stop Timeouts

[`start_timeout`](#timeouts.start_timeout){ #timeouts.start_timeout }

:   How long ECS will wait for your _container_ to start before giving up.
    
    Note that this is for container dependency resolution. You probably do not need to modify this

[`stop_timeout`](#timeouts.stop_timeout){ #timeouts.stop_timeout }

:   Time to wait before the container is forcefully killed if it doesn't exit normally on its own, when requested to terminate.

    Note that for spot containers, they will be killed after 2 minutes regardless of what you put here.

You can specify all durations as either seconds (as an integer) or using [Go Duration Format](https://pkg.go.dev/time#ParseDuration) (eg `2m` or `30s`)

----

## Container Health Checks

!!! warning
    AWS recommends that you do not use a docker health check

For the official documentation on specifying health checks, see [AWS ECS HealthCheck Documentation](https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_HealthCheck.html)

[`command`](#healthcheck.command){ #healthcheck.command } - **(required)**

:   The command to be run for the healthcheck. You must specify this as an array of strings, with the first element being either `CMD` or `CMD-SHELL`.

    See [ECS HealthCheck Command Docs](https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_HealthCheck.html#ECS-Type-HealthCheck-command)

[`interval`](#healthcheck.interval){ #healthcheck.interval }

:   The time period in seconds between each health check execution.

[`retries`](#healthcheck.retries){ #healthcheck.retries }

:   The number of times to retry a failed health check before the container is considered unhealthy.

[`start_period`](#healthcheck.start_period){ #healthcheck.start_period }

:   The optional grace period to provide containers time to bootstrap before failed health checks count towards the maximum number of retries.

[`timeout`](#healthcheck.timeout){ #healthcheck.timeout }

:   The time period in seconds to wait for a health check to succeed before it is considered a failure.

----

## Ulimits

Specify ulimits as an array of objects with the following properties:

[`name`](#ulimits.name){ #ulimits.name } - **(required)**

:   The name of the ulimit to adjust. Possible values are shown in the [AWS ECS Documentation](https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_Ulimit.html). Note that Fargate may not allow you to adjust all ulimits.

[`soft`](#ulimits.soft){ #ulimits.soft }

:   The value for the soft limit. Specified as an integer.

[`hard`](#ulimits.hard){ #ulimits.hard }

:   The value for the hard limit. Specified as an integer.


----

## Proxy Configuration

!!! warning ""
    You must specify the proxy container as a sidecar. If you do not specify it, then your tasks will fail.

=== "Example"

    ```yaml
    task_defaults:
      proxy:
        container_name: envoy
        properties:
          AppPorts: 5000
          IgnoredUID: 1000
          ...
    ```

=== "Shorthand Disable"

    ```yaml
    proxy: false
    ```

[`type`](#proxy.type){ #proxy.type }

:   The only acceptable value for this is `APPMESH`. This is the default, so it's recommended that you just leave this blank.

    _Default_: `{{schema:default:ProxyConfig.type}}`

[`container_name`](#proxy.container_name){ #proxy.container_name }

:   The name of the container that is providing the proxy.

    _Default_: `{{schema:default:ProxyConfig.container_name}}`

[`properties`](#proxy.properties){ #proxy.properties } - **(required)**

:   The properties for the proxy configuration.

    You can use the same syntax used for [defining environment variables](envvars.md), except you cannot specify any SSM parameters.

    For a list of properties, see the [AWS ECS ProxyConfiguration](https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_ProxyConfiguration.html) docs.

    _Default_: _empty_ (You must provide the required properties)

[`disabled`](#proxy.disabled){ #proxy.disabled }

:   If this is true, then the proxy configuration will be disabled for this task.

    _Default_: `{{schema:default:ProxyConfig.disabled}}`

----

## Private Registry
You can utilize ECS to connect to a private registry by providing the ARN of a SecretsManager secret for the `credentials` key.

This is normally not needed. If you are using ECR to host images, you do not need this. If you are unsure if you need this, then you don't need it.

For more details, look at the AWS documentation for [Private registry authentication for tasks](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/private-auth.html)
