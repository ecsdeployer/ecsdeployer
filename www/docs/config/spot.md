# Spot Configuration

To reduce the running costs of your application, you can easily launch containers using [Fargate Spot](https://aws.amazon.com/fargate/pricing/).

!!! danger "Make sure your application can handle Spot containers"
    Spot containers maybe terminated at any time, and you are only given 2 minutes notice before a task is killed.

    Containers will properly drain connections, so web applications will not suddenly throw errors to clients.

    If your application cannot handle tasks being killed suddenly (or you need the stability of a long running container), then do not enable Spot.

!!! tip "Read the Fargate Spot documentation"
    Please read the [Fargate Spot Documentation](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/fargate-capacity-providers.html) to make sure you understand any caveats.

    **Fargate Spot (as of Oct 2022) does not support ARM64 architecture**


### Enabling for entire application

```yaml
task_defaults:
  spot:
    enabled: true
```

!!! note ""
    **Note**: Fargate Spot is not used for [CronJobs](cronjobs.md), [PreDeploy Tasks](predeploy.md), or [Remote Shell](console.md). It is only used for [Services](services.md).


### Overriding a specific service
For services, you can override any of the spot configurations for that specific service. The parameters are the same as for the entire application, but you just need to put the block on the service.

```yaml title="Example"
project: example

task_defaults:
  spot:
    enabled: true
    minimum_ondemand: 0
  
services:
  - name: web
    # ...
    spot:
      minimum_ondemand: 5

```

## Advanced Capacity

You can also enable Spot containers, while still keeping a baseline (and/or a percentage) of your containers using OnDemand pricing.


```yaml
task_defaults:
  spot:
    enabled: true
    minimum_ondemand: 1
    minimum_ondemand_percent: 25
```

## Configuration Fields

[`enabled`](#spot.enabled){ #spot.enabled }

:   Enable Spot. If spot is disabled then the remaining configuration is irrelevant.

    _Default_: `false` (spot is disabled by default)

[`minimum_ondemand`](#spot.minimum_ondemand){ #spot.minimum_ondemand }

:   The minimum number of containers that should use OnDemand containers.
    
    This is equivalent to the `base` field on [CapacityProviderStategyItem](https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_CapacityProviderStrategyItem.html).
    
    _Default_: `0` (do not set an OnDemand baseline)

[`minimum_ondemand_percent`](#spot.minimum_ondemand_percent){ #spot.minimum_ondemand_percent }

:   Once the 'minimum_ondemand' count is met, what percentage of your service should be using OnDemand containers?
    
    This is equivalent to the `weight` field on [CapacityProviderStategyItem](https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_CapacityProviderStrategyItem.html).

    _Default_: `0`

See Also:

* [API CapacityProviderStategyItem](https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_CapacityProviderStrategyItem.html)
* [AWS Capacity Providers](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/cluster-capacity-providers.html)
