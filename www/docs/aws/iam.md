# AWS IAM Permissions

Below are two default policies that you can use for ECS Deployer.

=== "Normal Policy (recommended)"

    !!! note ""
        This policy is less restrictive, but will let you reuse a single role for all projects using the ECS Deployer.

    ```json
    --8<-- "./data/iam/lax.json"
    ```

=== "Restrictive Policy"

    !!! danger "Make sure to replace the placeholders!"

        These examples contain placeholders meant for you to replace with values for your environment. These are the placeholders:

        * `REGION` - The AWS region short code (`us-east-1`, `us-west-2`)
        * `ACCOUNTID` - Your numerical AWS Account ID
        * `CLUSTER_NAME` - Name of the ECS cluster specified for [`cluster`](../config/basic.md#root.cluster)
        * `PROJECT_NAME` - Value of [`project`](../config/basic.md#root.project)
        * `APP_ROLE` - Role you used for [`role`](../config/basic.md#root.role)
        * `ECS_EXECUTION_ROLE` - Role you used for [`execution_role`](../config/basic.md#root.execution_role)
        * `CRON_LAUNCHER_ROLE` - Role you used for [`cron_launcher_role`](../config/basic.md#root.cron_launcher_role)

    ```json
    --8<-- "./data/iam/restricted.json"
    ```
