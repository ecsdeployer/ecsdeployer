package service

import ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"

func (b *Builder) applyLoadBalancers() error {

	if !b.entity.IsLoadBalanced() {
		return nil
	}

	primaryContainerName, err := b.tplEval(*b.templates.ContainerName)
	if err != nil {
		return err
	}

	b.serviceDef.LoadBalancers = make([]ecsTypes.LoadBalancer, 0, len(b.entity.LoadBalancers))

	for _, lbInfo := range b.entity.LoadBalancers {

		targetGroupArn, err := lbInfo.TargetGroup.Arn(b.ctx)
		if err != nil {
			return err
		}

		b.serviceDef.LoadBalancers = append(b.serviceDef.LoadBalancers, ecsTypes.LoadBalancer{
			ContainerName:  &primaryContainerName,
			ContainerPort:  lbInfo.PortMapping.Port,
			TargetGroupArn: &targetGroupArn,
		})
	}

	b.serviceDef.HealthCheckGracePeriodSeconds = b.entity.LoadBalancers.GetHealthCheckGracePeriod()

	return nil
}
