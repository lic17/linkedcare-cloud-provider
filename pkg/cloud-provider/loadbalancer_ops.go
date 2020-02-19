package cloud_provider

import (
	"context"
	"fmt"

	"github.com/denverdino/aliyungo/slb"
	"github.com/golang/glog"
	v1 "k8s.io/api/core/v1"
	"k8s.io/kubernetes/pkg/cloudprovider"
)

type LoadBalancerClient struct {
	c slb.Client
}

func (s *LoadBalancerClient) findLoadBalancer(ctx context.Context, clusterName string, service *v1.Service) (status *v1.LoadBalancerStatus, exists bool, err error) {
	def, _ := ExtractServiceAnnotation(service)
	if def.Loadbalancerid != "" {
		return s.findLoadBalancerByID(def.Loadbalancerid)
	}
	// if not, find by slb name
	return s.findLoadBalancerByName(service)

}

func (s *LoadBalancerClient) findLoadBalancerByID(lbid string) (bool, *slb.LoadBalancerType, error) {

	lbs, err := s.c.DescribeLoadBalancers(
		&slb.DescribeLoadBalancersArgs{
			RegionId:       DEFAULT_REGION,
			LoadBalancerId: lbid,
		},
	)
	glog.Infof("find loadbalancer with id [%s], %d found.", lbid, len(lbs))
	if err != nil {
		return false, nil, err
	}

	if lbs == nil || len(lbs) == 0 {
		return false, nil, nil
	}
	if len(lbs) > 1 {
		glog.Warningf("multiple loadbalancer returned with id [%s], using the first one with IP=%s", lbid, lbs[0].Address)
	}
	lb, err := s.c.DescribeLoadBalancerAttribute(lbs[0].LoadBalancerId)
	return err == nil, lb, err
}

func (s *LoadBalancerClient) findLoadBalancerByName(service *v1.ServiceGetLoadBalancerName) (bool, *slb.LoadBalancerType, error) {
	if service.UID == "" {
		return false, nil, fmt.Errorf("unexpected empty service uid")
	}
	name := cloudprovider.GetLoadBalancerName(service)
	lbs, err := s.c.DescribeLoadBalancers(
		&slb.DescribeLoadBalancersArgs{
			RegionId:         DEFAULT_REGION,
			LoadBalancerName: name,
		},
	)
	glog.V(2).Infof("fallback to find loadbalancer by name [%s]", name)
	if err != nil {
		return false, nil, err
	}

	if lbs == nil || len(lbs) == 0 {
		return false, nil, nil
	}
	if len(lbs) > 1 {
		glog.Warningf("alicloud: multiple loadbalancer returned with name [%s], "+
			"using the first one with IP=%s", name, lbs[0].Address)
	}
	lb, err := s.c.DescribeLoadBalancerAttribute(lbs[0].LoadBalancerId)
	return err == nil, lb, err
}

func EnsureLoadBalancer(ctx context.Context, clusterName string, service *v1.Service, nodes []*v1.Node) (*v1.LoadBalancerStatus, error) {
}

func UpdateLoadBalancer(ctx context.Context, clusterName string, service *v1.Service, nodes []*v1.Node) error {
}

func EnsureLoadBalancerDeleted(ctx context.Context, clusterName string, service *v1.Service) error {

}
