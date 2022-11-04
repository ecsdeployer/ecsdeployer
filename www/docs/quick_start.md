---
hide:
  - toc

---
# Quick Start

## 1. Install ECS Deployer
This step is only necessary if you will be running ECS Deployer locally or somewhere other than GitHub Actions.

If you are using GitHub Actions, then please read our [GitHub Actions](ci/github.md) documentation.

Otherwise, see the [Install](install.md) page for more information.

## 2. Create a config file
Create your [Configuration File](config/index.md) for your application.

For applications with multiple environments/stages, you should place them in `.ecsdeployer/<EnvironmentName>.yml`.

If your app only has a single environment/stage, you can use `.ecsdeployer.yml` (or you can use the extended path above)


## 3. Ensure your user/role has sufficient permission
You will need the permissions listed on the [IAM Permissions](aws/iam.md) page.

## 4. Deploy!

```shell
ecsdeployer deploy --config path/to/your/config.yml --image-tag v1.2.3
```
