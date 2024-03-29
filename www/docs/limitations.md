# Limitations

## Current Limitations

#### No Automatic Rollback.

:    If there is a problem during deployment (or a task is failing), then the deployment will halt, but it will not attempt to revert to a previous version. You will need to run the deployer again and point to a prior version.

#### Only Fargate/FargateSpot is supported.

:   You could use this to deploy to an EC2 based cluster, but all your tasks would need to conform to Fargate requirements. (So no placement strategies, local volumes, bridge network, etc)

#### Limited validations / Ability for invalid AWS deployments

:    To allow the most flexibility, the deployer does not verify every option you provide. It will allow invalid configurations that will then fail to deploy on AWS. (ex: using a 16 core cpu on old platform versions, referencing roles that do not exist, clusters that do not exist.) This flexibility is allowed so that if AWS updates settings, you will be able to use those immediately without an update. But the tradeoff is you are reponsible for ensuring you have a valid configuration. Obvious errors will still be validated. These are mostly limited to logical but unsupported configurations. If you try to make a network port be 70000, you'll still get an error.

#### Windows is not supported

:    Currently the "Operating System" option is forced to LINUX. Windows support will be added in the future.

#### Only manages deployment-related infrastructure

:    The focus of this app is to make repeated deployments easy and fast. It purposely does not try to manage the entirety of the infrastructure required for your app. You must create the constant resources (Load balancers, target groups, ECS cluster, ECR repository, roles, etc.). Managing those is outside the scope of this program for the time being.

## Planned Features
* [ ] ECS Service Registries
* [x] AppMesh/Envoy Proxies
* [x] Additional Sidecars
* [ ] Autoscaling 
* [x] EFS Volume Mounts
* [x] Custom Container Logging (splunk, fluentd, etc)
* [x] Linking Service with multiple load balancers
* [x] Ability to override every possible option on a task (User, Ulimits, WorkDir, etc)
* [ ] Windows containers
* [ ] Pre/Post deploy hooks (run arbitrary commands)

## Possible Future Features
* [ ] Ability to launch tasks using EventBridge events as a trigger
* [ ] Using custom CapacityProviders
* [ ] More safeguards around invalid configurations.