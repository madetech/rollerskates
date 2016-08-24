package rollerskates

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type MockAwsAdapterFunnyState struct {
	removeFromLoadBalancerSpy chan [2]string
}

func (s MockAwsAdapterFunnyState) RemoveFromLoadBalancer(loadBalancerName string, instanceId string) {
	s.removeFromLoadBalancerSpy <- [2]string{loadBalancerName, instanceId}
}

func GetMockAwsAdapter() MockAwsAdapterFunnyState {
	return MockAwsAdapterFunnyState{removeFromLoadBalancerSpy: make(chan [2]string, 1)}
}

func GetLoadBalancerNameUsed(adapter MockAwsAdapterFunnyState) string {
	output := <-adapter.removeFromLoadBalancerSpy
	return output[0]
}

func GetInstanceIdUsed(adapter MockAwsAdapterFunnyState) string {
	output := <-adapter.removeFromLoadBalancerSpy
	return output[1]
}

func AssertLoadBalancerNameUsedIsEqualTo(expected string, t *testing.T, adapter MockAwsAdapterFunnyState) {
	assert.Equal(t, GetLoadBalancerNameUsed(adapter), expected)
}

func AssertInstanceIdUsedIsEqualTo(expected string, t *testing.T, adapter MockAwsAdapterFunnyState) {
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
