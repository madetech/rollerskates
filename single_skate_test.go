package rollerskates

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type RemoveFromLoadBalancerState struct {
	removeFromLoadBalancerSpy chan [2]string
}

func (s RemoveFromLoadBalancerState) RemoveFromLoadBalancer(loadBalancerName string, instanceId string) {
	s.removeFromLoadBalancerSpy <- [2]string{loadBalancerName, instanceId}
}

func GetMockAwsAdapter() RemoveFromLoadBalancerState {
	return RemoveFromLoadBalancerState{removeFromLoadBalancerSpy: make(chan [2]string, 1)}
}

func GetLoadBalancerNameUsed(adapter RemoveFromLoadBalancerState) string {
	output := <-adapter.removeFromLoadBalancerSpy
	return output[0]
}

func GetInstanceIdUsed(adapter RemoveFromLoadBalancerState) string {
	output := <-adapter.removeFromLoadBalancerSpy
	return output[1]
}

func AssertLoadBalancerNameUsedIsEqualTo(expected string, t *testing.T, adapter RemoveFromLoadBalancerState) {
	assert.Equal(t, GetLoadBalancerNameUsed(adapter), expected)
}

func AssertInstanceIdUsedIsEqualTo(expected string, t *testing.T, adapter RemoveFromLoadBalancerState) {
	assert.Equal(t, GetInstanceIdUsed(adapter), expected)
}

func TestGivenInstanceThenInstanceShouldBeRemoved1(t *testing.T) {
	adapter := GetMockAwsAdapter()
	expectedLoadBalancerName := "load-balancer-name"
	RestartLoadBalancerInstance(adapter, expectedLoadBalancerName, "")
	AssertLoadBalancerNameUsedIsEqualTo(expectedLoadBalancerName, t, adapter)
}

func TestGivenInstanceThenInstanceShouldBeRemoved(t *testing.T) {
	adapter := GetMockAwsAdapter()
	expectedInstanceId := "instance-id-to-restart"
	RestartLoadBalancerInstance(adapter, "", expectedInstanceId)
	AssertInstanceIdUsedIsEqualTo(expectedInstanceId, t, adapter)
}
