package awsclients

func STSClient() STSClienter {
	initMutex.RLock()
	defer initMutex.RUnlock()

	return stsClient
}

func SSMClient() SSMClienter {
	initMutex.RLock()
	defer initMutex.RUnlock()

	return ssmClient
}

func ECSClient() ECSClienter {
	initMutex.RLock()
	defer initMutex.RUnlock()

	return ecsClient
}

func EC2Client() EC2Clienter {
	initMutex.RLock()
	defer initMutex.RUnlock()

	return ec2Client
}

func ELBv2Client() ELBv2Clienter {
	initMutex.RLock()
	defer initMutex.RUnlock()

	return elbv2Client
}

func LogsClient() LogsClienter {
	initMutex.RLock()
	defer initMutex.RUnlock()

	return logsClient
}

func EventsClient() EventsClienter {
	initMutex.RLock()
	defer initMutex.RUnlock()

	return eventsClient
}

func TaggingClient() TaggingClienter {
	initMutex.RLock()
	defer initMutex.RUnlock()

	return taggingClient
}

func SchedulerClient() SchedulerClienter {
	initMutex.RLock()
	defer initMutex.RUnlock()

	return schedulerClient
}
