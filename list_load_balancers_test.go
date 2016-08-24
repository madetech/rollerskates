package rollerskates

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type AwsMockAdapterState struct {
	instanceIds []string
}

func ( s AwsMockAdapterState ) GetInstanceIds( loadBalancerName string ) [] string {
	return s.instanceIds
}

func Execute(loadBalancerName string, instanceIds []string) [] LoadBalancerInstance {
	awsAdapter := AwsMockAdapterState{instanceIds: instanceIds}
	return ListLoadBalancerInstances(awsAdapter, loadBalancerName)
}

func TestGivenOneLoadBalanceInstanceThenReturnsOneInstance(t *testing.T) {
	instances := Execute("some-client-loadbalance-name", []string{"fake-instance"})
	assert.True(t, len(instances) == 1)
}

func TestGivenNoLoadBalancerInstancesThenReturnsNoInstances(t *testing.T) {
	instances := Execute("some-client-loadbalance-name", []string{})
	assert.True(t, len(instances) == 0)
}

func TestGivenTwoLoadBalancerInstancesThenReturnsTwoInstances(t *testing.T) {
	instances := Execute("some-client-loadbalance-name", []string{"fake-instance", "second-fake-instance"})
	assert.True(t, len(instances) == 2)
}

func TestGivenThreeLoadBalancerInstancesThenReturnsThreeInstances(t *testing.T) {
	instances := Execute("some-client-loadbalance-name", []string{"fake-instance", "second-fake-instance", "the-last-slime"})
	assert.True(t, len(instances) == 3)
}

func TestGivenOneLoadBalancerInstanceThenInstanceShouldHaveId(t *testing.T) {
	instances := Execute("some-client-loadbalance-name", []string{"fake-instance"})
	assert.Equal(t, instances[0].id, "fake-instance")
}

