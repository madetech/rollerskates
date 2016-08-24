package rollerskates

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/joho/godotenv"
	_ "github.com/orchestrate-io/dvr"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCallToApi(t *testing.T) {
	godotenv.Load()

	ids := GetInstanceIds("testing")
	assert.Equal(t, ids[0], "i-e0708cd1")
}

// func TestReturnsTrueWhenConnectionsDrainedAndInstanceHasRestarted(t *testing.T) {
//   assert.Equal(t, RestartLoadBalancerInstance("testing", "i-e0708cd1"), true)
// }

func TestDeregisterInstancesFromLoadBalancer(t *testing.T) {
	assert.Equal(t, DeregisterInstancesFromLoadBalancer("testing", "i-e0708cd1"), true)
}

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

	svc := elb.New(session.New(), aws.NewConfig())

	paramsHealth := &elb.DescribeInstanceHealthInput{
		Instances: []*elb.Instance{
			{
				InstanceId: aws.String(instanceId),
			},
		},
		LoadBalancerName: aws.String(loadBalancerName),
	}

	errHealth := svc.WaitUntilInstanceDeregistered(paramsHealth)

	if errHealth != nil {
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
