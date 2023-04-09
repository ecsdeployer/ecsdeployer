package deprecate

import "ecsdeployer.com/ecsdeployer/pkg/config"

func Deprecate_LegacyCron(ctx *config.Context) {
	NoticeCustom(ctx, "legacy-cron", "Legacy CronJobs (Eventbridge: rules/targets) are deprecated and will be removed. Please switch to Eventbridge Scheduler. Check {{ .URL }} for more info")
}
