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
