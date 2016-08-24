package rollerskates

type LoadBalancerInstanceRemover interface {
	RemoveFromLoadBalancer(loadBalancerName string, instanceId string)
}

func RestartLoadBalancerInstance(remover LoadBalancerInstanceRemover, loadBalancerName string, instanceId string) {
	remover.RemoveFromLoadBalancer(loadBalancerName, instanceId)
}
