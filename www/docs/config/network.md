
# Network Config

Tasks running on Fargate must be launched within a VPC Subnet. The network configuration block will allow you to specify which subnets and security groups will be


## Usage
```yaml
network:

  subnets:
    - subnet-00000000000
    - subnet-11111111111

    - name: state
      values:
        - available

    - name: tag:cloud87/subnet_class
      value: host

  security_groups:
    - sg-1234567890
    - sg-9876543210
  
    - name: group-name
      values: ["cloud87-ecs"]

  public_ip: false
```

## Fields

[`subnets`](#network.subnets){ #network.subnets }

:   Specify SubnetIDs or a list of filters that can be used to locate the desired Subnets.

    For more information:

    * [Specifying Filters](#specifying-filters) (below)
    * [AWS Subnet Filter Values](https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_DescribeSubnets.html)

[`security_groups`](#network.security_groups){ #network.security_groups }

:   Specify SecurityGroupIDs or a list of filters that can be used to locate the desired SecurityGroups.

    For more information:

    * [Specifying Filters](#specifying-filters) (below)
    * [AWS Security Group Filter Values](https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_DescribeSecurityGroups.html)

[`public_ip`](#network.public_ip){ #network.public_ip }

:   Whether or not this task should be giving a public IP address.

    _Default_: `false`


## Specifying Filters

Filters are specified using a list of filter objects. Each object should have the following structure:


[`name`](#filters.name){ #filters.name } - **(required)**

:   The "Name" of the filter to use. Possible names dependent on whether you are using it for SecurityGroups or Subnets.

    You can view possible values for filter names at:

    * [Subnet Filters](https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_DescribeSubnets.html)
    * [Security Group Filters](https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_DescribeSecurityGroups.html)


[`values`](#filters.values){ #filters.values } - **(required)**

:   The filter values. Filter values are case-sensitive. If you specify multiple values for a filter, the values are joined with an OR, and the request returns all results that match any of the specified values.

    === "Single Value"

        ```yaml
        network:
          subnets:
            - name: vpc-id
              values: vpc-1234567
        ```

    === "Multi Value"

        ```yaml
        network:
          subnets:
            - name: availability-zones
              values:
                - us-east-1a
                - us-east-1b
                - us-east-1c
        ```
<!--
    === "Single Value Shorthand"

        ```yaml
        network:
          subnets:
            - vpc-id=vpc-1234567
        ```

    === "Multi Value Shorthand"

        ```yaml
        network:
          subnets:
            - availability-zones=us-east-1a,us-east-1b,us-east-1c
        ```

!!! note
    If you specify filters for subnets, then the VPC-ID of the first subnet will be used to constrain any security group filters. (It will add a `vpc-id` filter)

A filter set consists of a list of individual filters, each of which can be specified using any of the following formats:

### Explicit ID

If you want to hardcode your network information, you can.

=== "Format"

    ```yaml
    - subnet-111111111
    - sg-11111111
    ```

=== "Example"

    ```yaml
    network:
      subnets:
        - subnet-111111111
        - subnet-222222222

      security_groups:
        - sg-111111111
        - sg-222222222
    ```



### Simple Filter

This is the simplest method for specifying a filter, with the format of:


=== "Format"

    ```yaml
    - FILTER_KEY=FILTER_VALUE[,FILTER_VALUE[,FILTER_VALUE]]
    ```

=== "Example"

    ```yaml
    network:
      subnets:
        - "tag:cloud87/network=private"
        - state=available
    ```
-->


<!--

### Single Value Filter


=== "Format"

    ```yaml
    - name: FILTER_KEY
      value: FILTER_VALUE
    ```

=== "Example"

    ```yaml
    network:
      subnets:
        - name: "tag:cloud87/network"
          value: "private"
        - name: state
          value: available
    ```

### Multi-Value Filter

=== "Format"

    ```yaml
    - name: FILTER_KEY
      values:
        - FILTER_VALUE
        - FILTER_VALUE
        - FILTER_VALUE
    ```

=== "Example"

    ```yaml
    network:
      subnets:
        - name: "tag:cloud87/network"
          values: private
        - name: availability-zone
          value:
            - us-east-1a
            - us-east-1b
            - us-east-1c
    ```
-->

### Filter Reference

* [Subnet Filters](https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_DescribeSubnets.html)
* [Security Group Filters](https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_DescribeSecurityGroups.html)
* [AWS EC2 API::Filter](https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_Filter.html)