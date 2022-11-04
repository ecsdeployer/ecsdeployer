---
hide:
  toc: true
---
# ecsdeployer

Deploy applications to Fargate

## Synopsis

ECS Deployer allows you to easily deploy containerized applications to AWS ECS Fargate.
It simplifies the process of creating task definitions, running pre-deployment tasks, setting up scheduled jobs,
as well as service orchestration.

Applications can easily and securely be deployed with a simple GitHub Action.

Check out our website for more information, examples and documentation: https://ecsdeployer.com/


## Options

```
      --debug   Enable debug mode
  -h, --help    help for ecsdeployer
```

## See also

* [`ecsdeployer check`](ecsdeployer_check.md)	 - Checks if configuration is valid, validating it against the schema
* [`ecsdeployer clean`](ecsdeployer_clean.md)	 - Runs the cleanup step only. Skips actual deployment
* [`ecsdeployer deploy`](ecsdeployer_deploy.md)	 - Deploys application
* [`ecsdeployer info`](ecsdeployer_info.md)	 - Gives an overview of your project and what things are enabled
* [`ecsdeployer jsonschema`](ecsdeployer_jsonschema.md)	 - outputs ECS Deployer's JSON schema

