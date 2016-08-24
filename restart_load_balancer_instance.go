package rollerskates

type LoadBalancerInstanceRemover interface {
	RemoveFromLoadBalancer(loadBalancerName string, instanceId string)
}

type InstanceRestarter interface {
	RestartInstance(id string)
}

type RestartLoadBalancerInstanceDependencies struct {
	remover LoadBalancerInstanceRemover
	restarter InstanceRestarter
}

func RestartLoadBalancerInstance(deps RestartLoadBalancerInstanceDependencies, loadBalancerName string, instanceId string) {
	deps.restarter.RestartInstance(instanceId)
	deps.remover.RemoveFromLoadBalancer(loadBalancerName, instanceId)
}
