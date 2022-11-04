# Services


Services are defined as a list of objects under the `services` key.

## Example
```yaml
services:
   # A web server that is attached to a target group
  - name: web
    desired: 3
    command: ["bundle", "exec", "puma", "-c", "config/puma.rb"]
    load_balancer:
      target_group: c87-deployer-test-web
      port: 5000

  # A background worker
  # or some other service that does not need a load balancer
  - name: worker
    desired: 5
    command: ["bundle", "exec", "sidekiq"]
```

## Fields

[`name`](#service.name){ #service.name } - **(required)**

:   Name for your service. This should be short. 
    This name must be unique across all services, tasks, cronjobs, predeploys

[`desired`](#service.desired){ #service.desired } - **(required)**

:   The number of containers you wish to deploy for this service.<br>
    To disable a service, set this to `0`.

    _Default_: `0` (disabled)

[`command`](#service.command){ #service.command }

:   The command you want to run for this service.
    This field is not required, but it doesn't make much sense to exclude it.

    For details on how to specify a command, see [Specifying Command/EntryPoints](common.md#commandentry-point)

[`load_balancer`](#service.load_balancer){ #service.load_balancer }

:   Connect your service to a load balancer.

    See [Load Balancer](#load-balancer) below for more information.

[`rollout`](#service.rollout){ #service.rollout }

:   You can customize how the rollout process is performed if you want. 
    Normally, you should just leave this key out of your configuration and an ideal rollout configuration will be provided

    See [Rollout Configuration](#rollout-configuration) below for more information.

[`spot`](#service.spot){ #service.spot }

:   Allows you to override the [Spot Configuration](spot.md) for this single service.

[`skip_wait_for_stable`](#service.skip_wait_for_stable){ #service.skip_wait_for_stable }

:   If set to true, then the deployer will make no attempt to wait for this service to be stable before marking the deployment successful.
    You probably do not want to set this to true unless you have a service that takes really really long to reach stability
    
    _Default_: `{{schema:default:Service.skip_wait_for_stable}}`

[`<anything from common>`](#service.common){ #service.common }

:   See [Common Task Options](common.md).

### Load Balancer

[`port`](#load_balancer.port){ #load_balancer.port } - **(required)**

:   The port to route traffic to for the target group.

[`target_group`](#load_balancer.target_group){ #load_balancer.target_group } - **(required)**

:   The target group to connect this service to.
    You can specify either the name or the ARN.
    If you specify the name then the deployer will lookup the correct ARN

[`grace`](#load_balancer.grace){ #load_balancer.grace }

:   Optional health check grace period.
    If your app takes a long time to boot up before it can serve traffic, then you can set a grace period so that the load balancer does not kill it while it is still booting.
    <br>
    Specify as integer seconds or a duration

    _Default_: use AWS default

<small>_If you require multiple load balancers on a single service, you can define this block as an array to add multiple. (This is uncommon)_</small>


### Rollout Configuration
You can specify the minimum and maximum percentages that should be possible during a deployment. 
By default this will be configured for you with sensible defaults.

[`min`](#rollout.min){ #rollout.min } - **(required)**

:   Lower limit (as a percentage of the service's desiredCount) of the number of running tasks that must remain running and healthy in a service during a deployment.

[`max`](#rollout.max){ #rollout.max } - **(required)**

:   Upper limit (as a percentage of the service's desiredCount) of the number of running tasks that can be running in a service during a deployment

For more information, take a look at the [AWS ECS Deployment Configuration](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ecs-service-deploymentconfiguration.html).
