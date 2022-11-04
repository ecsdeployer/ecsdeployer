# Remote Shell (Console Task)

You can easily open an ephemeral shell into your application using the console task.

This is provided by [Remote SSH Docker Shell](https://github.com/webdestroya/remote-shell).
You can then easily connect to your application using the [Remote Shell Client](https://github.com/webdestroya/remote-shell-client).


## Client Usage

```sh
$ remote-shell -a myapp
spawning remote shell...
root@container$ 
```

And you're in! Now you can run commands in your app's actual environment. When you exit, the container will be terminated, and all files written will be deleted. (Note: If you make database modifications, or external changes, those will still persist.)


!!! warning ""
    You need to follow the [Usage](https://github.com/webdestroya/remote-shell#usage) instructions for the RemoteShell app in order to add it to your application.

This is not a task that runs normally. Only a task definition is created, and the task will be run on-demand when you try to launch a shell. Once your session finishes, the task will be killed. Any modifications to the filesystem will be lost once you exit the shell. You cannot modify already deployed code.


## Enabling Remote Shell

By default, the remote shell is not enabled. If you would like it to be created, you must explicitly enable it.


=== "Enable with defaults"

    ```yaml
    console: true
    ```

=== "Enable with customizations"

    ```yaml
    console:
      enabled: true
      port: 1234
    ```

=== "Disable"

    ```yaml
    console: false
    ```

## Fields

[`enabled`](#console.enabled){ #console.enabled }

:   Enables the Remote Shell. Obviously, this is required if you want the console.

    _Default_: `false` (disabled by default)

[`port`](#console.enabled){ #console.enabled }

:   Set the port that will be opened for this task. 

    _Default_: `8722`

[`path`](#console.enabled){ #console.enabled }

:   If you placed the RemoteShell binary in a non-standard location, then you should specify it here.

    _Default_: _not set_ (will use client default)

[`<anything from common>`](#console.common){ #console.common }

:   See [Common Task Options](common.md).

    You can override things like `cpu`, `memory` or `storage`, etc
