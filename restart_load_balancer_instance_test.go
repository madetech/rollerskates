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
	restartSpy     chan string
	orderSpy       chan string
	restartSuccess bool
}

type AddToLoadBalancerState struct {
	orderSpy chan string
}

func (s RemoveFromLoadBalancerState) RemoveFromLoadBalancer(loadBalancerName string, instanceId string) bool {
	s.orderSpy <- "RemoveFromLoadBalancer"
	s.removeSpy <- [2]string{loadBalancerName, instanceId}
	return s.removeSuccess
}

func (s RestartInstanceState) RestartInstance(id string) bool {
	s.orderSpy <- "RestartInstance"
	s.restartSpy <- id
	return s.restartSuccess
}

func (s AddToLoadBalancerState) AddInstanceToLoadBalancer(loadBalancerName string, instanceId string) bool {
	s.orderSpy <- "AddInstanceToLoadBalancer"
	return true
}

type RestartLoadBalancerInstanceMockDependencies struct {
	orderSpy  chan string
	remover   RemoveFromLoadBalancerState
	restarter RestartInstanceState
	adder     AddToLoadBalancerState
}

func (s RestartLoadBalancerInstanceMockDependencies) ConvertToProduction() RestartLoadBalancerInstanceDependencies {
	return RestartLoadBalancerInstanceDependencies{
		remover:   s.remover,
		restarter: s.restarter,
		adder:     s.adder,
	}
}

func GetDependencies(removeSuccess bool, restartSuccess bool) RestartLoadBalancerInstanceMockDependencies {
	globalSpyChannel := make(chan string, 3)

	return RestartLoadBalancerInstanceMockDependencies{
		orderSpy: globalSpyChannel,
		remover: RemoveFromLoadBalancerState{
			removeSuccess: removeSuccess,
			removeSpy:     make(chan [2]string, 1),
			orderSpy:      globalSpyChannel,
		},
		restarter: RestartInstanceState{
			restartSpy: make(chan string, 1),
			orderSpy:   globalSpyChannel,
			restartSuccess: restartSuccess,
		},
		adder: AddToLoadBalancerState{
			orderSpy: globalSpyChannel,
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

func AssertChannelNeverReceives(t *testing.T, never string, channel chan string) {
	select {
	case msg := <-channel:
		if msg == never {
			t.Fail()
		}
	default:
	//noop
	}
}

func ExecuteRestartLoadBalancerInstance(deps RestartLoadBalancerInstanceMockDependencies, loadBalancerName string, expectedInstanceId string) {
	RestartLoadBalancerInstance(deps.ConvertToProduction(), loadBalancerName, expectedInstanceId)
}

func TestGivenInstanceThenInstanceShouldBeRemoved1(t *testing.T) {
	deps := GetDependencies(true, true)
	expectedLoadBalancerName := "load-balancer-name"
	ExecuteRestartLoadBalancerInstance(deps, expectedLoadBalancerName, "")
	AssertLoadBalancerNameRemovedIsEqualTo(expectedLoadBalancerName, t, deps)
}

func TestGivenInstanceThenInstanceShouldBeRemoved(t *testing.T) {
	deps := GetDependencies(true, true)
	expectedInstanceId := "instance-id-to-restart"
	ExecuteRestartLoadBalancerInstance(deps, "", expectedInstanceId)
	AssertInstanceIdRemovedIsEqualTo(expectedInstanceId, t, deps)
}

func TestGivenInstanceRemovedThenShouldRestartInstance(t *testing.T) {
	deps := GetDependencies(true, true)
	expectedInstanceId := "instance-id-to-restart"
	ExecuteRestartLoadBalancerInstance(deps, "", expectedInstanceId)
	AssertInstanceIdRestartedIsEqualTo(expectedInstanceId, t, deps)
}

func TestInstanceRemovedBeforeRestart(t *testing.T) {
	deps := GetDependencies(true, true)
	expectedInstanceId := "instance-id-to-restart"
	ExecuteRestartLoadBalancerInstance(deps, "", expectedInstanceId)
	assert.Equal(t, "RemoveFromLoadBalancer", <-deps.orderSpy)
	assert.Equal(t, "RestartInstance", <-deps.orderSpy)
}

func TestGivenInstanceNotRemovedThenInstanceShouldNotBeRestarted(t *testing.T) {
	deps := GetDependencies(false, true)
	expectedInstanceId := "instance-id-to-restart"
	ExecuteRestartLoadBalancerInstance(deps, "", expectedInstanceId)
	assert.Equal(t, "RemoveFromLoadBalancer", <-deps.orderSpy)
	AssertChannelNeverReceives(t, "RestartInstance", deps.orderSpy)
}

func TestInstanceIsAddedToLoadBalancer(t *testing.T) {
	deps := GetDependencies(true, true)
	expectedInstanceId := "instance-id-to-restart"
	ExecuteRestartLoadBalancerInstance(deps, "", expectedInstanceId)
	assert.Equal(t, "RemoveFromLoadBalancer", <-deps.orderSpy)
	assert.Equal(t, "RestartInstance", <-deps.orderSpy)
	assert.Equal(t, "AddInstanceToLoadBalancer", <-deps.orderSpy)
}

func TestGivenRestartInstanceFailsThenInstanceIsNotAddedToLoadBalancer(t *testing.T) {
	deps := GetDependencies(true, false)
	expectedInstanceId := "instance-id-to-restart"
	ExecuteRestartLoadBalancerInstance(deps, "", expectedInstanceId)
	assert.Equal(t, "RemoveFromLoadBalancer", <-deps.orderSpy)
	assert.Equal(t, "RestartInstance", <-deps.orderSpy)
	AssertChannelNeverReceives(t, "AddInstanceToLoadBalancer", deps.orderSpy)
}


