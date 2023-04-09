# Deprecations

Features and/or options that are slated to be removed will be listed below. If you use one of those features and are not willing to migrate, then you will need to ensure you version lock your ECSDeployer installation, as newer versions will remove the feature.


### Legacy Cron {#legacy-cron}
The Cloudwatch Event Rule/Target based CronJob system has been deprecated. It has been replaced with [EventBridge Scheduler](https://docs.aws.amazon.com/scheduler/latest/UserGuide/what-is-scheduler.html).

If you have already used ECSDeployer with Rule/Targets, you can use older versions or you can disable your cronjobs manually, and then deploy with ECS Deployer (which will set them up using the Scheduler). Once deployed, you can delete the old rules/targets.

For a few versions of ECSDeployer, you can use `settings.use_old_cron_eventbus: true` to force using the old versions. Keep in mind this will be removed soon.