package cloud_provider

type ClientMgr struct {
	//instance InstanceClient
	loadbalancer LoadBalancerClient
}

// Instances return instance client
//func (mgr *ClientMgr) Instances() *InstanceClient { return mgr.instance }

// LoadBalancers return loadbalancer client
func (mgr *ClientMgr) LoadBalancers() *LoadBalancerClient { return mgr.loadbalancer }
