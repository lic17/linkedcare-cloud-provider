package cloud_provider

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/nas"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
)

type ClientMgr struct {
	stop <-chan struct{}

	instance     *InstanceClient
	loadbalancer *LoadBalancerClient
}

// NewClientMgr return a new client manager
func NewClientMgr(region, key, secret string) (*ClientMgr, error) {
	ecsclient, err := ecs.NewClientWithAccessKey(region, key, secret)
	if err != nil {
		return nil, err
	}
	slbclient, err := slb.NewClientWithAccessKey(region, key, secret)
	if err != nil {
		return nil, err
	}
	nasclient, err := nas.NewClientWithAccessKey(region, key, secret)
	if err != nil {
		return nil, err
	}
	mgr := &ClientMgr{
		stop: make(<-chan struct{}, 1),
		instance: &InstanceClient{
			c:   ecsclient,
			nas: nasclient,
		},
		loadbalancer: &LoadBalancerClient{
			ins: ecsclient,
			c:   slbclient,
		},
	}

	return mgr, nil
}

// Instances return instance client
func (mgr *ClientMgr) Instances() *InstanceClient { return mgr.instance }

// LoadBalancers return loadbalancer client
func (mgr *ClientMgr) LoadBalancers() *LoadBalancerClient { return mgr.loadbalancer }
