package rollerskates

type InstanceIdsRetriever interface {
	GetInstanceIds(loadBalancerName string) []string
}

func ListLoadBalancerInstances(retriever InstanceIdsRetriever, loadBalancerName string) []LoadBalancerInstance {
	instances := []LoadBalancerInstance{}

	for _, instanceId := range retriever.GetInstanceIds(loadBalancerName) {
		instances = append(instances, LoadBalancerInstance{id: instanceId})
	}

	return instances
}

type LoadBalancerInstance struct {
	id string
}
