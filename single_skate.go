package rollerskates

type AwsAdapterFunny interface {
	RemoveFromLoadBalancer(loadBalancerName string, instanceId string)
}

func RestartLoadBalancerInstance(awsAdapter AwsAdapterFunny, loadBalancerName string, instanceId string) {
	awsAdapter.RemoveFromLoadBalancer(loadBalancerName, instanceId)
}
