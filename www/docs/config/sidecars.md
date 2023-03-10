# Sidecars

## Usage
```yaml
task_defaults:
  sidecars:
    - name: datadog-agent
      image: public.ecr.aws/datadog/agent:latest
      cpu: 128
      memory: 512
      environment:
        DD_API_KEY: 
          ssm: /secrets/global/DD_API_KEY
        ECS_FARGATE: true
```

## Fields

[`name`](#sidecars.name){ #sidecars.name } - **(required)**

:   Name for your task. This should be short. 
    This name must be unique across all services, tasks, cronjobs, sidecarss

[`inherit_env`](#sidecars.inherit_env){ #sidecars.inherit_env }

:   By default, a sidecar will not inherit the environment variables from the primary container.

    If you want all environment variables on the primary container to be copied to the sidecar, set this to `true`.

    If you also specify environment variables for this container, those will override any environment variables inherited from the primary container.

    _Default_: `{{schema:default:Sidecar.inherit_env}}`

[`essential`](#sidecars.essential){ #sidecars.essential }

:   Mark this container as essential (or not essential). A non-essential container can die without the task as a whole being considered unhealthy.

    _Default_: `{{schema:default:Sidecar.essential}}`

[`port_mappings`](#sidecars.port_mappings){ #sidecars.port_mappings }

:   A list of container ports to expose on the task.

    See [Specifying Port Mappings](#specifying-port-mappings) below for the required format.

[`memory_reservation`](#sidecars.memory_reservation){ #sidecars.memory_reservation }

:   You can specify a memory reservation (minimum) if you want to allow your sidecar to use more memory if available.

    You can specify this in the same format as [`memory`](common.md#resources.memory), but you cannot use multiplier values.

[`<most things from common>`](#sidecars.common){ #sidecars.common }

:   You can use any container level property from [Common Task Options](common.md).

    This would exclude things that only apply at the task level (network, tags, arch).


## Specifying Port Mappings

=== "Shorthand"

    ```yaml
    port_mappings:
      - 8080 # protocol is assumed as tcp
      - 1234/tcp
      - 5000/udp
    ```

=== "Detailed"

    ```yaml
    port_mappings:
      - port: 1234
        protocol: tcp
      
      - port: 5678
        protocol: udp
    ```
  
[`port`](#port_mappings.port){ #port_mappings.port } - **(required)**

:   The port number to open. Must be between 1 and 65535.

[`protocol`](#port_mappings.protocol){ #port_mappings.protocol }

:   The protocol for the mapping.

    _Default_: `{{schema:default:PortMapping.protocol}}`