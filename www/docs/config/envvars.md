# Environment Variables


## Automatic import of SSM parameters
By default, ECS Deployer will attempt to import secrets from SSM Parameter store for you.

**For information on configuring SSM Parameter Import, see the [SSM Parameter Store Secrets Import](settings.md#ssm-import) section on the [Settings](settings.md) page.**

The basename of the parameter path is used as the Variable Name. The value is never read by the Deployer, as it is passed as a secret to your application.

!!! example
    * `/ecsdeployer/secrets/myapp/DATABASE_URL` added as `DATABASE_URL`
    * `/ecsdeployer/secrets/myapp/SECRET_API_KEY` added as `SECRET_API_KEY`


----

## Specifying Values

### SSM Parameter Store Reference
!!! success "This is secure and can be used for sensitive values"

```yaml
environment:
  # ...
  MY_VARIABLE_NAME:
    ssm: /path/or/arn/to/ssm/or/secrets-manager/value
```

You can specify values as either a path to an SSM Parameter, the full ARN for an SSM Parameter, or the full ARN of a SecretsManager secret.


### Insecure Values

!!! danger "Adding environment variables to your config file will keep them in plain text. This is not recommended for sensitive values."

=== "Plain String"

    You can specify the environment variable value as any primitive type that is castable to a string. (ie: string, boolean, number)

    These will be added to the task's environment as plaintext.

    ```yaml
    environment:
      PORT: 1234
      ALLOW_CLOWNS: false
      ALLOW_FUN: true
      PHASERS: stun
      GIGAWATTS: 1.21
    ```

=== "Templated Values"

    Same as Plain String, but you can now use [Templates](../templating.md) in the value.

    ```yaml
    environment:
      APP_DEPLOYED_AT: "{{ .Date }}"
      GITHUB_REF: "{{ .Env.GITHUB_REF_NAME }}"
    ```


## Example

```yaml
environment:
  PORT: 5000
  ENABLE_FUNTIME: true
  FOO: bar

  BUILD_REF_NAME: { template: "{{ .Env.GITHUB_REF_NAME }}" }

  SSM_VARIABLE: { ssm: "/path/to/some/ssm/variable" }

```