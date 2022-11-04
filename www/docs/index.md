# ECS Deployer

ECS Deployer allows you to easily deploy containerized applications to AWS ECS Fargate.
It simplifies the process of creating task definitions, running pre-deployment tasks, setting up scheduled jobs, as well as service orchestration.

Applications can easily and securely be deployed with a simple GitHub Action.

## Features
* [Service creation/updating](config/services.md)
* [Scheduled Jobs/Cron](config/cronjobs.md)
* [Pre Deploy Tasks](config/predeploy.md)
* [Spot Containers](config/spot.md)
* [Automatic Logging Setup](config/logging.md)
* [GitHub Actions support](ci/github.md)
* Seamless integration with [Remote Shell](https://github.com/webdestroya/remote-shell-client) to easily debug your application


## Resources ECS Deployer Manages
* ECS Services
* ECS Task Definitions
* EventBridge Rules and Targets (For CronJobs)
* Creation only of CloudWatch Log Groups (if desired)

#### Resources ECS Deployer DOES NOT manage
* ECS Cluster
* SSM Secrets
* Load Balancers / Target Groups
* IAM Roles
* VPC Resources


## Example Config

```yaml title=".ecsdeployer.yml"
--8<-- "./docs/static/examples/simple_web.yml"
```

## Basic Flow

``` mermaid
flowchart TD
  A[Start] --> B{{Has PreDeploy?}};
  subgraph predeploy[Run PreDeploy Tasks]
    direction LR
    PD1 --> PD2[...]
    PD2 --> PDN
  end
  subgraph svc[Service Deployment]
    direction TB
    svc1[http] --> swait{{Wait for stability}}
    svc2[...] --> swait
    svc3[worker] --> swait
    cus{{Create/Update Services}} --> svc1
    cus --> svc2
    cus --> svc3
  end
  B -->|Yes| predeploy;
  B -->|No|C{{Cronjobs}};
  predeploy --> C;
  C --> svc;
  svc --> D{{Deregister Old Tasks}}
  D --> E[Success!]
```
