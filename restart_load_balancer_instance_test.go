package rollerskates

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type RemoveFromLoadBalancerState struct {
	removeSuccess bool
	removeSpy     chan [2]string
	orderSpy      chan string
}

type RestartInstanceState struct {
	restartSpy chan string
	orderSpy   chan string
}

func (s RemoveFromLoadBalancerState) RemoveFromLoadBalancer(loadBalancerName string, instanceId string) bool {
	s.orderSpy <- "RemoveFromLoadBalancer"
	s.removeSpy <- [2]string{loadBalancerName, instanceId}
	return s.removeSuccess
}

func (s RestartInstanceState) RestartInstance(id string) bool {
	s.orderSpy <- "RestartInstance"
	s.restartSpy <- id
	return true
}

type RestartLoadBalancerInstanceMockDependencies struct {
	remover   RemoveFromLoadBalancerState
	restarter RestartInstanceState
}

func (s RestartLoadBalancerInstanceMockDependencies) ConvertToProduction() RestartLoadBalancerInstanceDependencies {
	return RestartLoadBalancerInstanceDependencies{
		remover:   s.remover,
		restarter: s.restarter,
	}
}

func GetDependencies(removeSuccess bool) RestartLoadBalancerInstanceMockDependencies {
	globalSpyChannel := make(chan string, 3)

	return RestartLoadBalancerInstanceMockDependencies{
		remover: RemoveFromLoadBalancerState{
			removeSuccess: removeSuccess,
			removeSpy:     make(chan [2]string, 1),
			orderSpy:      globalSpyChannel,
		},
		restarter: RestartInstanceState{
			restartSpy: make(chan string, 1),
			orderSpy:   globalSpyChannel,
		},
	}
}

func GetLoadBalancerNameRemoved(s RemoveFromLoadBalancerState) string {
	output := <-s.removeSpy
	return output[0]
}

func GetInstanceIdRemoved(s RemoveFromLoadBalancerState) string {
	output := <-s.removeSpy
	return output[1]
}

func GetInstanceIdRestarted(s RestartInstanceState) string {
	output := <-s.restartSpy
	return output
}

func AssertLoadBalancerNameRemovedIsEqualTo(expected string, t *testing.T, deps RestartLoadBalancerInstanceMockDependencies) {
	assert.Equal(t, GetLoadBalancerNameRemoved(deps.remover), expected)
}

func AssertInstanceIdRemovedIsEqualTo(expected string, t *testing.T, deps RestartLoadBalancerInstanceMockDependencies) {
	assert.Equal(t, GetInstanceIdRemoved(deps.remover), expected)
}

func AssertInstanceIdRestartedIsEqualTo(expected string, t *testing.T, deps RestartLoadBalancerInstanceMockDependencies) {
	assert.Equal(t, GetInstanceIdRestarted(deps.restarter), expected)
}

func AssertChannelHasNoMoreMessages(t *testing.T, channel chan string) {
	select {
	case <-channel:
		t.Fail()
	default:
		//noop
	}
}

func ExecuteRestartLoadBalancerInstance(deps RestartLoadBalancerInstanceMockDependencies, loadBalancerName string, expectedInstanceId string) {
	RestartLoadBalancerInstance(deps.ConvertToProduction(), loadBalancerName, expectedInstanceId)
}

func TestGivenInstanceThenInstanceShouldBeRemoved1(t *testing.T) {
	deps := GetDependencies(true)
	expectedLoadBalancerName := "load-balancer-name"
	ExecuteRestartLoadBalancerInstance(deps, expectedLoadBalancerName, "")
	AssertLoadBalancerNameRemovedIsEqualTo(expectedLoadBalancerName, t, deps)
}

func TestGivenInstanceThenInstanceShouldBeRemoved(t *testing.T) {
	deps := GetDependencies(true)
	expectedInstanceId := "instance-id-to-restart"
	ExecuteRestartLoadBalancerInstance(deps, "", expectedInstanceId)
	AssertInstanceIdRemovedIsEqualTo(expectedInstanceId, t, deps)
}

func TestGivenInstanceRemovedThenShouldRestartInstance(t *testing.T) {
	deps := GetDependencies(true)
	expectedInstanceId := "instance-id-to-restart"
	ExecuteRestartLoadBalancerInstance(deps, "", expectedInstanceId)
	AssertInstanceIdRestartedIsEqualTo(expectedInstanceId, t, deps)
}

func TestInstanceRemovedBeforeRestart(t *testing.T) {
	deps := GetDependencies(true)
	expectedInstanceId := "instance-id-to-restart"
	ExecuteRestartLoadBalancerInstance(deps, "", expectedInstanceId)
	assert.Equal(t, "RemoveFromLoadBalancer", <-deps.remover.orderSpy)
	assert.Equal(t, "RestartInstance", <-deps.remover.orderSpy)
}

func TestGivenInstanceNotRemovedThenInstanceShouldNotBeRestarted(t *testing.T) {
	deps := GetDependencies(false)
	expectedInstanceId := "instance-id-to-restart"
	ExecuteRestartLoadBalancerInstance(deps, "", expectedInstanceId)
	assert.Equal(t, "RemoveFromLoadBalancer", <-deps.remover.orderSpy)
	AssertChannelHasNoMoreMessages(t, deps.remover.orderSpy)
}
