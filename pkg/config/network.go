package config

import (
	"errors"
	"sync"
	"time"

	"ecsdeployer.com/ecsdeployer/internal/awsclients"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	eventTypes "github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
	schedulerTypes "github.com/aws/aws-sdk-go-v2/service/scheduler/types"
	"github.com/caarlos0/log"
)

var networkResolverMutex = sync.Mutex{}

type resolvedNetworkConfig struct {
	subnets        []string
	securityGroups []string
	publicIp       bool
}

type NetworkConfiguration struct {
	Subnets        []NetworkFilter `yaml:"subnets,omitempty" json:"subnets,omitempty" jsonschema:"description=List of SubnetIds or Subnet Filters"`
	SecurityGroups []NetworkFilter `yaml:"security_groups,omitempty" json:"security_groups,omitempty" jsonschema:"description=List of SecurityGroupIds or SecurityGroup Filters"`
	AllowPublicIp  *bool           `yaml:"public_ip,omitempty" json:"public_ip,omitempty" jsonschema:"description=Should the task be given a Public IP address?,default=false"`

	// if this was already resolved, we can cache it
	// resolved *ecstypes.NetworkConfiguration
	resolved       bool
	resolvedConfig *resolvedNetworkConfig
	resolveError   *error
}

func (nc *NetworkConfiguration) IsResolved() bool {
	return nc.resolvedConfig != nil || nc.IsResolveError()
}

func (nc *NetworkConfiguration) IsResolveError() bool {
	return nc.resolveError != nil
}

func (nc *NetworkConfiguration) resolve(ctx *Context) error {
	if nc.IsResolveError() {
		return *nc.resolveError
	}
	if nc.IsResolved() {
		return nil
	}

	networkResolverMutex.Lock()
	defer networkResolverMutex.Unlock()

	result, err := networkConfigurationResolver(ctx, nc)
	if err != nil {
		nc.resolveError = &err
		return err
	}

	nc.resolved = true
	nc.resolvedConfig = result

	return nil
}

// Good thing AWS has a bunch of different incompatible, yet identical types. Neat
func (nc *NetworkConfiguration) Resolve(ctx *Context, netConfRef any) error {
	err := nc.resolve(ctx)
	if err != nil {
		return err
	}

	switch ref := netConfRef.(type) {
	case *ecsTypes.NetworkConfiguration:
		*ref = ecsTypes.NetworkConfiguration{
			AwsvpcConfiguration: &ecsTypes.AwsVpcConfiguration{
				Subnets:        nc.resolvedConfig.subnets,
				AssignPublicIp: util.Ternary(nc.resolvedConfig.publicIp, ecsTypes.AssignPublicIpEnabled, ecsTypes.AssignPublicIpDisabled),
				SecurityGroups: nc.resolvedConfig.securityGroups,
			},
		}
	case *eventTypes.NetworkConfiguration:
		*ref = eventTypes.NetworkConfiguration{
			AwsvpcConfiguration: &eventTypes.AwsVpcConfiguration{
				Subnets:        nc.resolvedConfig.subnets,
				AssignPublicIp: util.Ternary(nc.resolvedConfig.publicIp, eventTypes.AssignPublicIpEnabled, eventTypes.AssignPublicIpDisabled),
				SecurityGroups: nc.resolvedConfig.securityGroups,
			},
		}
	case *schedulerTypes.NetworkConfiguration:
		*ref = schedulerTypes.NetworkConfiguration{
			AwsvpcConfiguration: &schedulerTypes.AwsVpcConfiguration{
				Subnets:        nc.resolvedConfig.subnets,
				AssignPublicIp: util.Ternary(nc.resolvedConfig.publicIp, schedulerTypes.AssignPublicIpEnabled, schedulerTypes.AssignPublicIpDisabled),
				SecurityGroups: nc.resolvedConfig.securityGroups,
			},
		}
	case nil:
		// do nothing, they just want to validate that the config is valid
	default:
		return errors.New("unknown network configuration type??")
	}

	return nil
}

func (a *NetworkConfiguration) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type tNetworkConfiguration NetworkConfiguration
	var obj = tNetworkConfiguration{}
	if err := unmarshal(&obj); err != nil {
		return err
	}
	*a = NetworkConfiguration(obj)

	a.ApplyDefaults()

	return a.Validate()
}

func (nc *NetworkConfiguration) ApplyDefaults() {
}

func (nc *NetworkConfiguration) Validate() error {
	if len(nc.Subnets) > 0 {
		for _, f := range nc.Subnets {
			if err := f.Validate(); err != nil {
				return err
			}
		}
	}

	if len(nc.SecurityGroups) > 0 {
		for _, f := range nc.SecurityGroups {
			if err := f.Validate(); err != nil {
				return err
			}
		}
	}
	return nil
}

func networkConfigurationResolver(ctx *Context, network *NetworkConfiguration) (*resolvedNetworkConfig, error) {

	// NOTE: Don't parallelize this. The security group finder sometimes depends on the subnet's response (to add the vpc filter)

	log.Debug("resolving network configuration")
	startTime := time.Now()

	defaultNetwork := ctx.Project.TaskDefaults.Network

	if network == nil {
		network = defaultNetwork
	}

	if network == nil {
		return nil, errors.New("you must provide network configuration information")
	}

	result := &resolvedNetworkConfig{
		subnets:        []string{},
		securityGroups: []string{},
		publicIp:       false,
	}

	pubIpNetwork, ok := util.CoalesceWithFunc(func(val *NetworkConfiguration) bool {
		return val.AllowPublicIp != nil
	}, network, defaultNetwork, ctx.Project.Network)
	if ok && *pubIpNetwork.AllowPublicIp {
		result.publicIp = true
	}

	// TODO: maybe do subnets and security groups in parallel to be faster

	subnetNetwork, ok := util.CoalesceWithFunc(func(val *NetworkConfiguration) bool {
		return len(val.Subnets) > 0
	}, network, defaultNetwork, ctx.Project.Network)
	if !ok {
		return nil, errors.New("could not determine network subnets. No configuration provided")
	}

	subnetIds, vpcId, err := calculateSubnets(ctx, *subnetNetwork)
	if err != nil {
		return nil, err
	}
	result.subnets = subnetIds

	securityGroupNetwork, ok := util.CoalesceWithFunc(func(val *NetworkConfiguration) bool {
		return len(val.SecurityGroups) > 0
	}, network, defaultNetwork, ctx.Project.Network)
	if !ok {
		return nil, errors.New("could not determine network security groups. No configuration provided")
	}

	securityGroupIds, err := calculateSecurityGroups(ctx, *securityGroupNetwork, vpcId)
	if err != nil {
		return nil, err
	}
	result.securityGroups = securityGroupIds

	log.WithField("duration", time.Since(startTime).Truncate(time.Second)).Debug("resolved network configuration")

	return result, nil
}

func calculateSubnets(ctx *Context, network NetworkConfiguration) ([]string, *string, error) {

	var subnetIds []string = []string{}

	idFilters, nfFilters := splitNetworkFiltersByType(network.Subnets)

	// if they gave ids, then add those to the list
	if len(network.Subnets) > 0 {
		subnetIds = idFilters
	}

	// no filters, not worth wasting time
	if len(nfFilters) == 0 {
		return subnetIds, nil, nil
	}

	ec2Client := awsclients.EC2Client()

	request := &ec2.DescribeSubnetsInput{}

	request.Filters = make([]ec2Types.Filter, len(nfFilters))
	for i, filter := range nfFilters {
		request.Filters[i] = filter.ToAws()
	}

	paginator := ec2.NewDescribeSubnetsPaginator(ec2Client, request, func(o *ec2.DescribeSubnetsPaginatorOptions) {
		o.Limit = 10
	})

	var vpcId *string

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx.Context)
		if err != nil {
			return nil, nil, err
		}
		for _, value := range output.Subnets {
			vpcId = value.VpcId
			subnetIds = append(subnetIds, *value.SubnetId)
		}
	}

	return subnetIds, vpcId, nil
}

func calculateSecurityGroups(ctx *Context, network NetworkConfiguration, vpcId *string) ([]string, error) {

	securityGroupIds, nfFilters := splitNetworkFiltersByType(network.SecurityGroups)

	if len(nfFilters) == 0 {
		return securityGroupIds, nil
	}

	ec2Client := awsclients.EC2Client()

	request := &ec2.DescribeSecurityGroupsInput{}
	hasVpcFilter := false
	request.Filters = make([]ec2Types.Filter, len(nfFilters), len(nfFilters)+1)
	for i, filter := range nfFilters {
		request.Filters[i] = filter.ToAws()
		if *filter.Name == "vpc-id" {
			hasVpcFilter = true
		}
	}

	if !hasVpcFilter && vpcId != nil {
		request.Filters = append(request.Filters, ec2Types.Filter{
			Name:   aws.String("vpc-id"),
			Values: []string{*vpcId},
		})
	}

	paginator := ec2.NewDescribeSecurityGroupsPaginator(ec2Client, request, func(o *ec2.DescribeSecurityGroupsPaginatorOptions) {
		o.Limit = 10
	})

	// var securityGroupIds []string = []string{}

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx.Context)
		if err != nil {
			return nil, err
		}
		for _, value := range output.SecurityGroups {
			securityGroupIds = append(securityGroupIds, *value.GroupId)
		}
	}

	return securityGroupIds, nil
}
