package rollerskates

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
)

func GetInstanceIds(loadBalancerName string) []string {
	svc := elb.New(session.New(), aws.NewConfig())

	params := &elb.DescribeLoadBalancersInput{
		LoadBalancerNames: []*string{aws.String("testing")},
	}

	resp, err := svc.DescribeLoadBalancers(params)
	if err != nil {
		return []string{}
	}

	instances := []string{}
	for _, instance := range resp.LoadBalancerDescriptions[0].Instances {
		instances = append(instances, *instance.InstanceId)
	}

	return instances
}

func RestartLoadBalancerInstance(loadBalancerName string, instanceId string) bool {
	deregistered := DeregisterInstancesFromLoadBalancer(loadBalancerName, instanceId)
	if !deregistered {
		return false
	}

	drained := DeregisterInstancesFromLoadBalancer(loadBalancerName, instanceId)
	if !drained {
		return false
	}

  restarted := RebootInstance(instanceId)
  if !restarted {
    return false
  }

	return true
}

func DeregisterInstancesFromLoadBalancer(loadBalancerName string, instanceId string) bool {
	svc := elb.New(session.New(), aws.NewConfig())

	params := &elb.DeregisterInstancesFromLoadBalancerInput{
		Instances: []*elb.Instance{
			{
				InstanceId: aws.String(instanceId),
			},
		},
		LoadBalancerName: aws.String(loadBalancerName),
	}

	_, err := svc.DeregisterInstancesFromLoadBalancer(params)

	if err != nil {
		return false
	}

	return true
}

func WaitForInstancesToBeDeregisteredFromLoadBalancer(loadBalancerName string, instanceId string) bool {
	svc := elb.New(session.New(), aws.NewConfig())

	params := &elb.DescribeInstanceHealthInput{
		Instances: []*elb.Instance{
			{
				InstanceId: aws.String(instanceId),
			},
		},
		LoadBalancerName: aws.String(loadBalancerName),
	}

	err := svc.WaitUntilInstanceDeregistered(params)

	if err != nil {
		return false
	}

	return true
}

func RebootInstance(instanceId string) bool {
	svc := ec2.New(session.New(), aws.NewConfig())

	params := &ec2.RebootInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceId),
		},
	}

  _, err := svc.RebootInstances(params)

	if err != nil {
		return false
	}

	return true
}

func WaitForInstaceToBeRebooted(instanceId string) bool {
	svc := ec2.New(session.New(), aws.NewConfig())

	params := &ec2.DescribeInstancesInput{
    InstanceIds: []*string{
			aws.String(instanceId),
		},
	}

	err := svc.WaitUntilInstanceRunning(params)

	if err != nil {
		return false
	}

	return true
}
