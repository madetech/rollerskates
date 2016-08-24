package rollerskates

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type RemoveFromLoadBalancerState struct {
	removeSpy chan [2]string
}

type RestartInstanceState struct {
	restartSpy chan string
}

func (s RemoveFromLoadBalancerState) RemoveFromLoadBalancer(loadBalancerName string, instanceId string) {
	s.removeSpy <- [2]string{loadBalancerName, instanceId}
}

func (s RestartInstanceState) RestartInstance(id string) {
	s.restartSpy <- id
}

type RestartLoadBalancerInstanceMockDependencies struct {
	remover RemoveFromLoadBalancerState
	restarter RestartInstanceState
}

func (s RestartLoadBalancerInstanceMockDependencies) ConvertToProduction() RestartLoadBalancerInstanceDependencies {
	return RestartLoadBalancerInstanceDependencies{
		remover: s.remover,
		restarter: s.restarter,
	}
}

func GetDependencies() RestartLoadBalancerInstanceMockDependencies {
	return RestartLoadBalancerInstanceMockDependencies{
		remover: RemoveFromLoadBalancerState{
			removeSpy: make(chan [2]string, 1),
		},
		restarter: RestartInstanceState{
			restartSpy: make(chan string, 1),
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

func ExecuteRestartLoadBalancerInstance(deps RestartLoadBalancerInstanceMockDependencies, loadBalancerName string, expectedInstanceId string ) {
	RestartLoadBalancerInstance( deps.ConvertToProduction(), loadBalancerName, expectedInstanceId )
}

func TestGivenInstanceThenInstanceShouldBeRemoved1(t *testing.T) {
	deps := GetDependencies()
	expectedLoadBalancerName := "load-balancer-name"
	ExecuteRestartLoadBalancerInstance(deps, expectedLoadBalancerName, "")
	AssertLoadBalancerNameRemovedIsEqualTo(expectedLoadBalancerName, t, deps)
}

func TestGivenInstanceThenInstanceShouldBeRemoved(t *testing.T) {
	deps := GetDependencies()
	expectedInstanceId := "instance-id-to-restart"
	ExecuteRestartLoadBalancerInstance(deps, "", expectedInstanceId)
	AssertInstanceIdRemovedIsEqualTo(expectedInstanceId, t, deps)
}

func TestGivenInstanceRemovedThenShouldRestartInstance(t *testing.T) {
	deps := GetDependencies()
	expectedInstanceId := "instance-id-to-restart"
	ExecuteRestartLoadBalancerInstance(deps, "", expectedInstanceId)
	AssertInstanceIdRestartedIsEqualTo(expectedInstanceId, t, deps)
}

