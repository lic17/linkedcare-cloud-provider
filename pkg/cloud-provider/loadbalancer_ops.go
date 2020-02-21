package cloud_provider

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/golang/glog"
	v1 "k8s.io/api/core/v1"
	"k8s.io/kubernetes/pkg/cloudprovider"
)

type LoadBalancerClient struct {
	c   *slb.Client
	ins *ecs.Client
}

func (s *LoadBalancerClient) findLoadBalancer(service *v1.Service) (bool, *slb.LoadBalancer, error) {
	def, _ := ExtractServiceAnnotation(service)
	if def.Loadbalancerid != "" {
		return s.findLoadBalancerByID(def.Loadbalancerid)
	}
	// if not, find by slb name
	return s.findLoadBalancerByName(service)

}

func (s *LoadBalancerClient) findLoadBalancerByID(lbid string) (bool, *slb.LoadBalancer, error) {

	request := slb.CreateDescribeLoadBalancersRequest()
	//request.Scheme = "https"
	request.LoadBalancerId = lbid

	res, err := s.c.DescribeLoadBalancers(request)
	lbs := res.LoadBalancers.LoadBalancer
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

	fmt.Println(lbs[0])
	return err == nil, &lbs[0], err
}

func (s *LoadBalancerClient) findLoadBalancerByName(service *v1.Service) (bool, *slb.LoadBalancer, error) {
	if service.UID == "" {
		return false, nil, fmt.Errorf("unexpected empty service uid")
	}
	name := cloudprovider.GetLoadBalancerName(service)

	request := slb.CreateDescribeLoadBalancersRequest()
	//request.Scheme = "https"
	request.LoadBalancerName = name
	res, err := s.c.DescribeLoadBalancers(request)
	lbs := res.LoadBalancers.LoadBalancer
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
	fmt.Println(lbs[0])
	return err == nil, &lbs[0], err
}

func (s *LoadBalancerClient) ensureLoadBalancer(service *v1.Service) (bool, *slb.LoadBalancer, error) {
	exists, origined, err := s.findLoadBalancer(service)
	if err != nil {
		return false, nil, err
	}
	request, _ := ExtractServiceAnnotation(service)

	if !exists {
		//if isServiceDeleted(service) {
		//	glog.V(2).Infof("alicloud: isServiceDeleted report that this service has been " +
		//		"deleted before. see issue: https://github.com/kubernetes/kubernetes/issues/59084")
		//	os.Exit(1)
		//}
		if request.Loadbalancerid != "" {
			return false, nil, fmt.Errorf("alicloud: user specified "+
				"loadbalancer[%s] does not exist. pls check", request.Loadbalancerid)
		}

		// From here, we need to create a new loadbalancer
		glog.V(5).Infof("alicloud: can not find a "+
			"loadbalancer with service name [%s/%s], creating a new one", service.Namespace, service.Name)
		// If need created, double check if the resource id has been deleted
	} else {
		// Need to verify loadbalancer.
		// Reuse SLB is not allowed when the SLB is created by k8s service.
		return exists, origined, err
	}
	return exists, origined, err
}

func (s *LoadBalancerClient) updateLoadBalancer(service *v1.Service) error {
	return nil
}

/*********************/
/*delete loadbalancer*/
/*********************/
func (s *LoadBalancerClient) ensureLoadBalancerDeleted(service *v1.Service) error {
	// need to save the resource version when deleted event
	//err := keepResourceVersion(service)
	//if err != nil {
	//	glog.Warningf(service, "Warning: failed to save deleted service resourceVersion,due to [%s] ", err.Error())
	//}
	exists, lb, err := s.findLoadBalancer(service)
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}
	// skip delete user defined loadbalancer
	if isUserDefinedLoadBalancer(service) {
		glog.Warningf("user managed loadbalancer will not be deleted by cloudprovider.")
		return nil
	}
	request := slb.CreateDeleteLoadBalancerRequest()
	//request.Scheme = "https"

	request.LoadBalancerId = lb.LoadBalancerId
	_, err = s.c.DeleteLoadBalancer(request)

	return err
}

// check to see if user has assigned any loadbalancer
func isUserDefinedLoadBalancer(svc *v1.Service) bool {
	return serviceAnnotation(svc, ServiceAnnotationLoadBalancerId) != ""
}
