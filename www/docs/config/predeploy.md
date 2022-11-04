# PreDeploy Tasks

PreDeploy tasks are tasks that are run **sequentially** before your services are deployed.
Tasks must succeed, otherwise the whole deployment will be halted.

```yaml
predeploy:
  - name: migrations
    command: ["bundle", "exec", "rake", "db:migrate"]
  
  - name: warm-caches
    command: ["bundle", "exec", "rake", "caches:warm"]
```


## Fields

[`name`](#predeploy.name){ #predeploy.name } - **(required)**

:   Name for your task. This should be short. 
    This name must be unique across all services, tasks, cronjobs, predeploys

[`command`](#predeploy.command){ #predeploy.command } - _(recommended)_

:   The command you want to run for this task.
    This field is not required, but it doesn't make much sense to exclude it.

    For details on how to specify a command, see [Specifying Command/EntryPoints](common.md#commandentry-point)

[`timeout`](#predeploy.timeout){ #predeploy.timeout }

:   Set a custom timeout for this task.
    if the task does not complete within the allotted time, it will be killed
    and considered as a failure.
    
    Note: this is the total launch time + runtime.
    so if you have a massive image, that will count against the timeout
    
    _Default_: inherit from [`settings.predeploy_timeout`](settings.md#settings.predeploy_timeout)

[`disabled`](#predeploy.disabled){ #predeploy.disabled }

:   Skip running this task. 
    This allows you to easily skip certain tasks while still keeping them in the configuration
    
    _Default_: `{{schema:default:PreDeployTask.disabled}}`

[`ignore_failure`](#predeploy.ignore_failure){ #predeploy.ignore_failure }

:   If you want to run a task, but do not care if it fails, then set this to true.
    Usually you will not want this.
    
    _Default_: `{{schema:default:PreDeployTask.ignore_failure}}`

[`<anything from common>`](#predeploy.command){ #predeploy.command }

:   See [Common Task Options](common.md).