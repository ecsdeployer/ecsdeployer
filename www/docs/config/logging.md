# Logging

Logging for your tasks will be automatically setup based upon your specified configuration.

If you do not specify any logging configuration, then CloudWatch logs will be used.

## Example Usage

=== "Using CloudWatch Logs (default)"

    ```yaml
    logging:
      awslogs:
        retention: 180
    ```


=== "Using FireLens"

    ```yaml
    logging:
      firelens:
        type: "fluentbit"
        options:
          enable-ecs-log-metadata: true
          Name: thing
          region: "us-east-1"
          delivery_stream: "my-stream"
          log-driver-buffer-limit: "2097152"
    ```

## Fields

[`awslogs`](#logging.awslogs){ #logging.awslogs }

:   See [Using CloudWatch Logs](#using-cloudwatch-logs)

    _Conflicts with [`firelens`](#logging.firelens)_

    **This is the default if you do not specify a `logging` block**

[`firelens`](#logging.firelens){ #logging.firelens }

:   See [Using FireLens](#using-firelens)

    _Conflicts with [`awslogs`](#logging.awslogs)_

[`disabled`](#logging.disabled){ #logging.disabled }

:   Disable all logging

    _Default_: `false`

!!! note ""
    **Note:** If you enable Firelens logging, then AwsLogs will be ignored. You cannot have both at the same time.

----

## Using CloudWatch Logs



!!! note ""
    Note: When using CloudWatch logs, the deployer will automatically create any missing log groups for you.

### Fields

[`retention`](#awslogs.retention){ #awslogs.retention }

:   How many days logs should be kept.

    To keep logs forever, specify `forever`

    Valid values can be found on the [PutRetentionPolicy](https://docs.aws.amazon.com/AmazonCloudWatchLogs/latest/APIReference/API_PutRetentionPolicy.html#API_PutRetentionPolicy_RequestSyntax) documentation.

    _Default_: `{{schema:default:AwsLogConfig.retention}}` (days)

[`options`](#awslogs.options){ #awslogs.options }

:   Allows you to specify extra options for the awslogs driver.

    Note: `awslogs-group`, `awslogs-stream-prefix` and `awslogs-region` are already set for you.

    Valid options: [awslogs driver options](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/using_awslogs.html#create_awslogs_logdriver_options)

[`disabled`](#awslogs.disabled){ #awslogs.disabled }

:   Disable CloudWatch logs.

    _Default_: `{{schema:default:AwsLogConfig.disabled}}`

**See Also**

* [`name_template.log_group`](naming.md#naming.log_group)
* [`name_template.log_stream_prefix`](naming.md#naming.log_stream_prefix)
* [AWS Documentation on using CloudWatch Logs](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/using_awslogs.html)

----

## Using FireLens

!!! note ""
    Note: When using firelens, a container dependency will be automatically created for you. Your primary container will depend on the firelens router.

### Fields

[`type`](#firelens.retention){ #firelens.retention }

:   The FireLens log router flavor you want to use.<br>Possible Options:

    * `fluentbit` - (preferred, default)
    * `fluentd`

    _Default_: `{{schema:default:FirelensConfig.type}}`

[`options`](#firelens.options){ #firelens.options }

:   Allows you to specify options to provide for the individual task container log configurations.

    Options can be specified identically to how [Environment Variables](envvars.md) are specified. Values that reference an SSM Parameter will be added to the `SecretOptions` field on the container. All others will be added to `Options`.

    Unless you are using a custom log router image that has these values already set for you, this field is most likely required.

    !!! example "Usage example"
        ```yaml
        logging:
          firelens:
            type: fluentbit
            options:
              Name: thing
              region: "us-east-1"
              delivery_stream: "my-stream"
              log-driver-buffer-limit: "2097152"
              whatever: {template: "{{ .Date }}"}
        ```

[`router_options`](#firelens.router_options){ #firelens.router_options }

:   Allows you to specify options that will be passed to the Firelens **router container only**. These are not applied to the task containers.

    Options can be specified identically to how [Environment Variables](envvars.md) are specified, **but you cannot specify any SSM parameters**.

    You generally will not need to configure this. This is meant for advanced customization of Firelens.

    !!! example "Usage example"
        ```yaml
        logging:
          firelens:
            type: fluentbit
            options:
              Name: thing
              region: "us-east-1"
              delivery_stream: "my-stream"
            router_options:
              enable-ecs-log-metadata: true
        ```    

[`memory`](#firelens.memory){ #firelens.memory }

:   Memory _reservation_ for this container.

    This will be used as the reservation setting for the container within this task.
    This should generally be low, as this will prevent the logging container from stealing memory from the primary container

    _Default_: `{{schema:default:FirelensConfig.memory}}` (megabytes)

### Advanced Fields

[`image`](#firelens.image){ #firelens.image }

:   Override the FireLens router image to use.

    If you are using **fluentbit**, then this will default to the [official fluentbit image](https://gallery.ecr.aws/aws-observability/aws-for-fluent-bit)

    If you are using **fluentd** then you must specify this.

    See [Specifying Images](basic.md#specifying-container-image) for more information

[`inherit_env`](#firelens.inherit_env){ #firelens.inherit_env }

:   Should the logging container inherit all the environment variables that were provided to the primary container?
    
    Note: if you have many environment variables, you might encounter task size limits.
    
    _Default_: `{{schema:default:FirelensConfig.inherit_env}}`

[`environment`](#firelens.environment){ #firelens.environment }

:   Add extra environment variables that are specific to this container (and will not be shared to the other containers in the task.)

    For more info, read the [Environment Variables](envvars.md) documentation.

[`container_name`](#firelens.container_name){ #firelens.container_name }

:   The name of the logging container. Normally you should not be changing this.

    _Default_: `{{schema:default:FirelensConfig.container_name}}`

[`log_to_awslogs`](#firelens.log_to_awslogs){ #firelens.log_to_awslogs }

:   Whether the firelens container itself should log to CloudWatch logs or not. This is helpful if you are debugging issues with your log router, but otherwise is not necessary.

    This is either:

    * Boolean `false` - disable logging to Cloudwatch Logs (default)
    * String (the log group you want to log to) - You are responsible for making this group.

    _Default_: `{{schema:default:FirelensConfig.log_to_awslogs}}`

[`credentials`](#firelens.credentials){ #firelens.credentials }

:   Optional [private registry](common.md#private-registry) credentials

[`disabled`](#firelens.disabled){ #firelens.disabled }

:   Disable FireLens entirely.

    _Default_: `{{schema:default:FirelensConfig.disabled}}`

