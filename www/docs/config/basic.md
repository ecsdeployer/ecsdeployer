# Configuration Root Fields


```yaml
# .ecsdeployer/production.yml
project: myapp
cluster: mycluster
image: ...
role: ...
services:
  - ...
cronjobs:
  - ...
settings: ...
# ...
```

## Fields


[`project`](#root.project){ #root.project } - **(required)**

:   The project name will be used to construct the names of all tasks.
    It should be short, and should only contain letters, numbers, dashes

    ```yaml
    project: deployer-test
    ```

[`cluster`](#root.cluster){ #root.cluster } - **(required)**

:   You must specify a `cluster` key that provides either an [ECS Cluster](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/clusters.html) name or ARN.

    === "Using Name Only"

        ```yaml
        cluster: mycluster
        ```

    === "Using ARN"

        ```yaml
        cluster: arn:aws:ecs:us-east-1:1234567890:cluster/mycluster
        ```

[`image`](#root.image){ #root.image } - _(recommended)_

:   See [Specifying Container Image](#specifying-container-image) below.

[`role`](#root.role){ #root.role } - _(recommended)_

:   This is the application role (Task Role).
    This is what your application will use when it is running normally.
    This is highly recommended, although not required.

    You can specify roles using the full ARN or the name of the role. Roles must already exist, they will not be created for you.

[`execution_role`](#root.execution_role){ #root.execution_role } - **(required)**

:   This is the role that AWS ECS will use to run your tasks.

    Read more: [Amazon ECS task execution IAM](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task_execution_IAM_role.html) 

    You can specify roles using the full ARN or the name of the role. Roles must already exist, they will not be created for you.

[`cron_launcher_role`](#root.cron_launcher_role){ #root.cron_launcher_role } - **(required if using [CronJobs](cronjobs.md))**

:   This is the role that EventBridge will use to launch your task based on the schedule you specify.

    You can specify roles using the full ARN or the name of the role. Roles must already exist, they will not be created for you.

[`network`](#root.network){ #root.network } - **(required)**

:   _See [Networking](network.md)_

[`services`](#root.services){ #root.services } - _(recommended)_

:   _See [Services](services.md)_

    This is not required, but if you don't provide it... why even use this project?

[`cronjobs`](#root.cronjobs){ #root.cronjobs }

:   _See [CronJobs](cronjobs.md)_

[`predeploy`](#root.predeploy){ #root.predeploy }

:   _See [PreDeploy Tasks](predeploy.md)_

[`console`](#root.console){ #root.console }

:   _See [Remote Shell](console.md)_

[`environment`](#root.environment){ #root.environment }

:   _See [Environment Variables](envvars.md)_

[`task_defaults`](#root.task_defaults){ #root.task_defaults }

:   _See [Task Defaults](defaults.md)_

[`name_templates`](#root.name_templates){ #root.name_templates }

:   _See [Naming](naming.md)_

[`logging`](#root.logging){ #root.logging }

:   _See [Logging](logging.md)_

[`tags`](#root.tags){ #root.tags }

:   Specify a set of tags that will be applied to all resources.

    See [Tags Documentation](tags.md) for more information

[`settings`](#root.settings){ #root.settings }

:   _See [Settings](settings.md)_

## Specifying Container Image

=== "Specifying as an object"

    ```yaml
    image:
      ecr: myapp
      tag: "{{ .ImageTag }}"
    ```

    [`ecr`](#image.ecr){ #image.ecr }

    :   The [AWS Elastic Container Repository](https://aws.amazon.com/ecr/) to pull from.

        You can specify this as:
        
        * a name (which will be converted to a full URL for you)
        * a full hostname/path

        _Conflicts with [`docker`](#image.docker)_

    [`docker`](#image.docker){ #image.docker }

    :   Reference a [DockerHub](https://hub.docker.com/) repository.

        _Conflicts with [`ecr`](#image.ecr)_

    [`tag`](#image.tag){ #image.tag }

    :   Pulls an image using its tag. This is the most common usage.

        Examples: `latest`, `v1.2.3`

        _Default_: `{{ .ImageTag }}` (what you pass in as `--image-tag`)

        _Conflicts with [`digest`](#image.digest)_

    [`digest`](#image.digest){ #image.digest } 

    :   Pulls an image using its digest value.

        Should be in the format of `sha256:<hexvalue>`

        _Conflicts with [`tag`](#image.tag)_

    **Note:** You may use [templates](../templating.md) in all fields.

=== "Specifying as a string"

    You can specify the full URI to a container image as a string.


    You may use [templates](../templating.md) in the string.

    **Examples:**

    ```yaml
    image: "01234567890.dkr.ecr.REGION.amazonaws.com/thing/stuff:latest"
    ```

    ```yaml
    image: "01234567890.dkr.ecr.REGION.amazonaws.com/thing/stuff:{{ .ImageTag }}"
    ```

    ```yaml
    image: "nginx:1.2.3"
    ```

