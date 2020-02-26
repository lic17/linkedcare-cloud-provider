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
	fmt.Println(name)

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

func (s *LoadBalancerClient) ensureLoadBalancer(service *v1.Service, nodes []*v1.Node) (bool, *slb.LoadBalancer, error) {
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

		/*Create LoadBalancer*/

		name := cloudprovider.GetLoadBalancerName(service)
		create_lbs_request := slb.CreateCreateLoadBalancerRequest()
		//create_lbs_request.Scheme = "https"
		create_lbs_request.LoadBalancerName = name
		create_lbs_request.LoadBalancerSpec = "slb.s1.small"
		create_lbs_request.AddressType = "internet"

		lbs_response, err := s.c.CreateLoadBalancer(create_lbs_request)
		if err != nil {
			fmt.Print(err.Error())
		}
		fmt.Println(lbs_response)
		err = BuildVirturalGroupFromService(s, service, nodes, lbs_response.LoadBalancerId)
		if err != nil {
			fmt.Print(err.Error())
		}
		/*Create Vserver Group for nodes*/

		/*
			for _, port := range service.Spec.Ports {
				create_vsg_request := slb.CreateCreateVServerGroupRequest()
				create_vsg_request.VServerGroupName = service.Namespace + "-" + service.Name

				var servers []VsgBackendServer
				for _, node := range nodes {
					server := new(VsgBackendServer)
					nodeid, err := nodeFromProviderID(node.Spec.ProviderID)
					if err != nil {
						fmt.Print(err.Error())
						break
					}

					server.ServerId = nodeid
					server.Weight = "100"
					server.Type = "ecs"
					server.Port = strconv.Itoa(int(port.NodePort))
				    server.Description = "tcp-" + server.Port
					servers = append(servers, *server)
				}
				jsonServers, _ := json.Marshal(servers)
				//create_vsg_request.Scheme = "https"

				create_vsg_request.BackendServers = string(jsonServers)
				//create_vsg_request.BackendServers = "[{ \"ServerId\": \"i-xxxxxxxxx\", \"Weight\": \"100\", \"Type\": \"ecs\", \"Port\":\"80\",\"Description\":\"test-112\" }]"
				create_vsg_request.LoadBalancerId = lbs_response.LoadBalancerId

				vsg_response, err := s.c.CreateVServerGroup(create_vsg_request)
				if err != nil {
					fmt.Print(err.Error())
				}
				fmt.Println(vsg_response)

				//Create TCPListener of LoadBalancer
				create_tcp_request := slb.CreateCreateLoadBalancerTCPListenerRequest()
				//create_tcp_request.Scheme = "https"

				create_tcp_request.VServerGroupId = vsg_response.VServerGroupId
				create_tcp_request.ListenerPort = requests.Integer(port.Port)
				create_tcp_request.LoadBalancerId = lbs_response.LoadBalancerId
				create_tcp_request.BackendServerPort = requests.Integer(port.NodePort)

				tcp_response, err := s.c.CreateLoadBalancerTCPListener(create_tcp_request)
				if err != nil {
					fmt.Print(err.Error())
				}
				fmt.Println(tcp_response)
			}
		*/

		// From here, we need to create a new loadbalancer
		glog.V(5).Infof("alicloud: can not find a "+
			"loadbalancer with service name [%s/%s], creating a new one", service.Namespace, service.Name)
		exists, origined, err := s.findLoadBalancer(service)
		if err != nil {
			return false, nil, err
		}
		return exists, origined, err
		// If need created, double check if the resource id has been deleted
	} else {
		// Need to verify loadbalancer.
		// Reuse SLB is not allowed when the SLB is created by k8s service.
		return exists, origined, err
	}
	return exists, origined, err
}

func (s *LoadBalancerClient) updateLoadBalancer(service *v1.Service, nodes []*v1.Node) error {
	exists, lb, err := s.findLoadBalancer(service)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the loadbalance you specified by name [%s] does not exist", service.Name)
	}

	err = BuildVirturalGroupFromService(s, service, nodes, lb.LoadBalancerId)
	if err != nil {
		return err
	}
	//if err := EnsureVirtualGroups(vgs, nodes); err != nil {
	//	return fmt.Errorf("update backend servers: error %s", err.Error())
	//}

	//if !needUpdateDefaultBackend(service, lb) {
	//	return nil
	//}
	//utils.Logf(service, "update default backend server group")
	//return s.UpdateDefaultServerGroup(nodes, lb)
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
