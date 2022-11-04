---
hide:
  - toc
---

# GitHub Actions

ECS Deployer can also be used on [GitHub Actions](https://github.com/features/actions) using our [official action](https://github.com/ecsdeployer/ecsdeployer-action).

## Usage

!!! tip
    **For the most up-to-date info, please visit the [ECS Deployer GitHub Action](https://github.com/ecsdeployer/ecsdeployer-action) page.**

### Recommend Directory Structure

```
<repo root>
├── .ecsdeployer
│   ├── beta.yml
│   ├── production.yml
│   └── staging.yml
└── .github
    └── workflows
      └── deploy.yml
```

## Configuring

### Action Inputs

[`config`](#action.config){ #action.config }

:   Path to the configuration file that will be used for deploying your application.

    _Default_: `.ecsdeployer.yml`

[`image`](#action.image){ #action.image }

:   Sets the container image URI to use as the primary image.

    If you pass this, you should **not** define an [`image`](../config/basic.md#root.image) section on your project root.
    If provide this parameter and have an [`image`](../config/basic.md#root.image) section defined, the section in the file will win.

[`image-tag`](#action.image-tag){ #action.image-tag }

:   Optional value passed as `--image-tag` to deployment. This is used for the [`{{ .Tag }}`](../templating.md#common-fields) template variable.

    _Default_: (not set)

[`app-version`](#action.app-version){ #action.app-version }

:   Optional value passed as `--app-version` to deployment. This is used for the [`{{ .AppVersion }}`](../templating.md#common-fields) template variable.

    _Default_: (not set)

[`extra-args`](#action.extra-args){ #action.extra-args }

:   Additional arguments to pass to the deploy command

    _Default_: (none)

[`workdir`](#action.workdir){ #action.workdir }

:   Working directory (below repository root)

    _Default_: `.`

[`timeout`](#action.timeout){ #action.timeout }

:   Sets the timeout for the entire deployment process.

    _Default_: (not set - use the default timeout for ECS Deployer)

[`install-only`](#action.install-only){ #action.install-only }

:   Just install ECS Deployer, and then exit

    _Default_: `false`

[`deployer-version`](#action.deployer-version){ #action.deployer-version }

:   Override the version of ECSDeployer to run.

    _Default_: `latest`


## Example Workflow

=== "`.github/workflows/deploy.yml`"

    ```yaml title=".github/workflows/deploy.yml" 
    --8<-- "./docs/static/github/workflow.yml"
    ```
    [Download this file](../static/github/workflow.yml)

## Further Reading
* [GitHub Actions Workflow Syntax](https://help.github.com/en/articles/workflow-syntax-for-github-actions#About-yaml-syntax-for-workflows)
* [aws-actions/configure-aws-credentials](https://github.com/aws-actions/configure-aws-credentials)
* [Configuring OpenID Connect in Amazon Web Services](https://docs.github.com/en/actions/deployment/security-hardening-your-deployments/configuring-openid-connect-in-amazon-web-services)