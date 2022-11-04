---
hide:
  - toc
---

# Customize Naming

If you want to modify how the deployer names various resources, you can add a custom section with the following:

For more details on the options you can use in templates, view the [Templating](../templating.md) documentation.


!!! danger "Danger! Advanced Topic!"

    Under normal operation you should not need to modify any of these values.
    Incorrectly modifying the templates used could seriously harm your application and cause downtime.

    **If you do not understand what these are doing, you should not adjust them**


## Example Usage

```yaml
name_templates:
  log_group: "/applogs/{{ .ProjectName }}/{{ .Name }}"
```

## Fields

[`log_group`](#naming.log_group){ #naming.log_group } - **(required)**

:   This is the log group path that will be created for each task.
    This is not used if logging has been disabled.

    === "Default"
        ```
        {{schema:default:TplDefault.log_group}}
        ```

    === "Default (with Stage)"
        ```
        {{schema:default:TplDefaultStage.log_group}}
        ```

    === "Raw"
        ```
        {{schema:default:NameTemplates.log_group}}
        ```

[`log_stream_prefix`](#naming.log_stream_prefix){ #naming.log_stream_prefix } - **(required)**

:   Fargate tasks require a logging prefix.

    === "Default"
        ```
        {{schema:default:TplDefault.log_stream_prefix}}
        ```
    
    === "Default (with Stage)"
        ```
        {{schema:default:TplDefaultStage.log_stream_prefix}}
        ```

    === "Raw"
        ```
        {{schema:default:NameTemplates.log_stream_prefix}}
        ```

[`service_name`](#naming.service_name){ #naming.service_name } - **(required)**

:   The name used for ECS Services. If you need a custom prefix, you can modify this

    === "Default"
        ```
        {{schema:default:TplDefault.service_name}}
        ```
    
    === "Default (with Stage)"
        ```
        {{schema:default:TplDefaultStage.service_name}}
        ```

    === "Raw"
        ```
        {{schema:default:NameTemplates.service_name}}
        ```

[`task_family`](#naming.task_family){ #naming.task_family } - **(required)**

:   This will be the TaskDefinition Family Name

    === "Default"
        ```
        {{schema:default:TplDefault.task_family}}
        ```
    
    === "Default (with Stage)"
        ```
        {{schema:default:TplDefaultStage.task_family}}
        ```

    === "Raw"
        ```
        {{schema:default:NameTemplates.task_family}}
        ```

[`cron_rule`](#naming.cron_rule){ #naming.cron_rule } - **(required)**

:   The EventBridge rule name when creating [CronJobs](cronjobs.md)

    === "Default"
        ```
        {{schema:default:TplDefault.cron_rule}}
        ```
    
    === "Default (with Stage)"
        ```
        {{schema:default:TplDefaultStage.cron_rule}}
        ```

    === "Raw"
        ```
        {{schema:default:NameTemplates.cron_rule}}
        ```

[`cron_target`](#naming.cron_target){ #naming.cron_target } - **(required)**

:   The EventBridge target name when creating [CronJobs](cronjobs.md)

    === "Default"
        ```
        {{schema:default:TplDefault.cron_target}}
        ```
    
    === "Default (with Stage)"
        ```
        {{schema:default:TplDefaultStage.cron_target}}
        ```

    === "Raw"
        ```
        {{schema:default:NameTemplates.cron_target}}
        ```

[`cron_group`](#naming.cron_group){ #naming.cron_group }

:   This is the 'Group' field when a task is run via a cron job.
    If you have special tracking of this field, you can enter a value here

    If you enter something that evaluates to an empty string, then the group field will not be included.

    === "Default"
        ```
        {{schema:default:TplDefault.cron_group}}
        ```
    
    === "Default (with Stage)"
        ```
        {{schema:default:TplDefaultStage.cron_group}}
        ```

    === "Raw"
        ```
        {{schema:default:NameTemplates.cron_group}}
        ```

[`predeploy_group`](#naming.predeploy_group){ #naming.predeploy_group }

:   Same as the [`cron_group`](#naming.cron_group), this is only used when predeploy tasks are run

    If you enter something that evaluates to an empty string, then the group field will not be included.

    === "Default"
        ```
        {{schema:default:TplDefault.predeploy_group}}
        ```
    
    === "Default (with Stage)"
        ```
        {{schema:default:TplDefaultStage.predeploy_group}}
        ```

    === "Raw"
        ```
        {{schema:default:NameTemplates.predeploy_group}}
        ```

[`predeploy_started_by`](#naming.predeploy_started_by){ #naming.predeploy_started_by }

:   This is the 'StartedBy' field on RunTask.

    If you enter something that evaluates to an empty string, then the group field will not be included.

    === "Default"
        ```
        {{schema:default:TplDefault.predeploy_started_by}}
        ```
    
    === "Default (with Stage)"
        ```
        {{schema:default:TplDefaultStage.predeploy_started_by}}
        ```

    === "Raw"
        ```
        {{schema:default:NameTemplates.predeploy_started_by}}
        ```

[`marker_tag_key`](#naming.marker_tag_key){ #naming.marker_tag_key } - **(required)**

:   The marker tag is something added to all resources created by the deployer.
    This allows the deployer to delete resources that it created, but that have since been removed from the deployment configuration file.
  
    **Once you deploy a project, you should not modify the values used for the marker tag.**

    === "Default"
        ```
        {{schema:default:TplDefault.marker_tag_key}}
        ```
    
    === "Default (with Stage)"
        ```
        {{schema:default:TplDefaultStage.marker_tag_key}}
        ```

    === "Raw"
        ```
        {{schema:default:NameTemplates.marker_tag_key}}
        ```

[`marker_tag_value`](#naming.marker_tag_value){ #naming.marker_tag_value } - **(required)**

:   Value used for the tag defined by [`marker_tag_key`](#naming.marker_tag_key)

    === "Default"
        ```
        {{schema:default:TplDefault.marker_tag_value}}
        ```
    
    === "Default (with Stage)"
        ```
        {{schema:default:TplDefaultStage.marker_tag_value}}
        ```

    === "Raw"
        ```
        {{schema:default:NameTemplates.marker_tag_value}}
        ```
