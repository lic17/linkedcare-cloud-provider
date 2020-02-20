package cloud_provider

import (
	"context"
	"fmt"

	"github.com/golang/glog"
	v1 "k8s.io/api/core/v1"
	"k8s.io/kubernetes/pkg/cloudprovider"
)

func (c *Cloud) LoadBalancer() (cloudprovider.LoadBalancer, bool) {
	return c, true
}

func (c *Cloud) GetLoadBalancer(ctx context.Context, clusterName string, service *v1.Service) (status *v1.LoadBalancerStatus, exists bool, err error) {

	exists, lb, err := c.climgr.LoadBalancers().findLoadBalancer(service)

	if err != nil || !exists {
		return nil, exists, err
	}

	return &v1.LoadBalancerStatus{
		Ingress: []v1.LoadBalancerIngress{{
			IP: lb.Address,
		}}}, true, nil
}

func (c *Cloud) GetLoadBalancerName(ctx context.Context, clusterName string, service *v1.Service) string {
	return cloudprovider.GetLoadBalancerName(service)
}

func (c *Cloud) EnsureLoadBalancer(ctx context.Context, clusterName string, service *v1.Service, nodes []*v1.Node) (*v1.LoadBalancerStatus, error) {

	glog.V(2).Infof("Alicloud.EnsureLoadBalancer(%v, %s/%s, %v, %v)",
		clusterName, service.Namespace, service.Name, c.region, NodeList(nodes))
	//defaulted, _ := ExtractAnnotationRequest(service)
	//if defaulted.AddressType == slb.InternetAddressType {
	//	if c.cfg != nil && c.cfg.Global.DisablePublicSLB {
	//		return nil, fmt.Errorf("PublicAddress SLB is Not allowed")
	//	}
	//}

	//ns, err := c.fileOutNode(nodes, service)
	//if err != nil {
	//	return nil, err
	//}

	if len(service.Spec.Ports) == 0 {
		return nil, fmt.Errorf("requested load balancer with no ports")
	}

	_, lb, err := c.climgr.LoadBalancers().ensureLoadBalancer(service)
	if err != nil {
		return nil, err
	}

	//pz, pzr, err := c.climgr.PrivateZones().EnsurePrivateZoneRecord(service, lb.Address, defaulted.AddressIPVersion)
	//if err != nil {
	//	return nil, err
	//}

	return &v1.LoadBalancerStatus{
		Ingress: []v1.LoadBalancerIngress{
			{
				IP: lb.Address,
			},
		},
	}, nil

}

func (c *Cloud) UpdateLoadBalancer(ctx context.Context, clusterName string, service *v1.Service, nodes []*v1.Node) error {
	glog.V(2).Infof("Alicloud.UpdateLoadBalancer(%v, %v, %v, %v, %v, %v, %v)",
		clusterName, service.Namespace, service.Name, c.region, service.Spec.LoadBalancerIP, service.Spec.Ports, NodeList(nodes))
	//ns, err := c.fileOutNode(nodes, service)
	//if err != nil {
	//	return err
	//}
	return c.climgr.LoadBalancers().updateLoadBalancer(service)

}

func (c *Cloud) EnsureLoadBalancerDeleted(ctx context.Context, clusterName string, service *v1.Service) error {
	glog.V(2).Infof("Alicloud.EnsureLoadBalancerDeleted(%v, %v, %v, %v, %v, %v)",
		clusterName, service.Namespace, service.Name, c.region, service.Spec.LoadBalancerIP, service.Spec.Ports)

	//defaulted, _ := ExtractAnnotationRequest(service)

	//	if len(service.Status.LoadBalancer.Ingress) > 0 {
	//		err := c.climgr.PrivateZones().EnsurePrivateZoneRecordDeleted(service, service.Status.LoadBalancer.Ingress[0].IP, defaulted.AddressIPVersion)
	//		if err != nil {
	//			return err
	//		}
	//	}

	return c.climgr.LoadBalancers().ensureLoadBalancerDeleted(service)

}
