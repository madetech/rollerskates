package rollerskates

type AwsAdapter interface {
	GetInstanceIds(loadBalancerName string) []string
}

func ListLoadBalancerInstances(awsAdapter AwsAdapter, loadBalancerName string) []LoadBalancerInstance {
	instances := []LoadBalancerInstance{}

	for _, instanceId := range awsAdapter.GetInstanceIds(loadBalancerName) {
		instances = append(instances, LoadBalancerInstance{id: instanceId})
	}

	return instances
}

type LoadBalancerInstance struct {
	id string
}
