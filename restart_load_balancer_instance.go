package rollerskates

type LoadBalancerInstanceRemover interface {
	RemoveFromLoadBalancer(loadBalancerName string, instanceId string) bool
}

type InstanceRestarter interface {
	RestartInstance(id string) bool
}

type RestartLoadBalancerInstanceDependencies struct {
	remover   LoadBalancerInstanceRemover
	restarter InstanceRestarter
}

func RestartLoadBalancerInstance(deps RestartLoadBalancerInstanceDependencies, loadBalancerName string, instanceId string) {
	removalSuccessful := deps.remover.RemoveFromLoadBalancer(loadBalancerName, instanceId)
	if removalSuccessful {
		deps.restarter.RestartInstance(instanceId)
	}
}
