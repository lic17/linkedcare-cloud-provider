package cloud_provider

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	v1 "k8s.io/api/core/v1"
)

// NodeList return nodes list in string
func NodeList(nodes []*v1.Node) []string {
	ns := []string{}
	for _, node := range nodes {
		ns = append(ns, node.Name)
	}
	return ns
}

/*
	service := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:   "default",
			Name:        "service-test",
			UID:         types.UID("1f11ce6d-5782-11ea-ae49-00163f00bfd3"),
			Annotations: map[string]string{},
		},
		Spec: v1.ServiceSpec{
			Type: "LoadBalancer",
		},
	}
	node := []*v1.Node{
		{
			ObjectMeta: metav1.ObjectMeta{Name: prid},
			Spec: v1.NodeSpec{
				ProviderID: prid,
			},
		},
	}
*/

func BuildVirturalGroupFromService(s *LoadBalancerClient, service *v1.Service, nodes []*v1.Node, lbsId string) error {

	/*Create Vserver Group for nodes*/
	for _, port := range service.Spec.Ports {
		vsg_name := service.Namespace + "-" + service.Name + "-" + strconv.Itoa(int(port.NodePort))
		hadVsg := false
		hadListener := false

		//first to get vsg
		get_vsg_request := slb.CreateDescribeVServerGroupsRequest()
		//get_vsg_request.Scheme = "https"
		get_vsg_request.LoadBalancerId = lbsId

		get_vsg_response, err := s.c.DescribeVServerGroups(get_vsg_request)
		if err != nil {
			fmt.Print(err.Error())
			return err
		}
		for _, vsg := range get_vsg_response.VServerGroups.VServerGroup {
			if !hadVsg {
				if vsg.VServerGroupName == vsg_name {
					hadVsg = true
					break
				}
			}
		}

		vsgId := ""

		//if don't  have vsg to create vsg
		if !hadVsg {
			create_vsg_request := slb.CreateCreateVServerGroupRequest()
			create_vsg_request.VServerGroupName = vsg_name

			var servers []VsgBackendServer
			for i, node := range nodes {
				if i < 5 {
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
			}
			jsonServers, _ := json.Marshal(servers)
			//create_vsg_request.Scheme = "https"

			create_vsg_request.BackendServers = string(jsonServers)

			//create_vsg_request.BackendServers = "[{ \"ServerId\": \"i-xxxxxxxxx\", \"Weight\": \"100\", \"Type\": \"ecs\", \"Port\":\"80\",\"Description\":\"test-112\" }]"
			create_vsg_request.LoadBalancerId = lbsId

			vsg_response, err := s.c.CreateVServerGroup(create_vsg_request)
			if err != nil {
				fmt.Print(err.Error())
				return err
			}
			vsgId = vsg_response.VServerGroupId
			fmt.Println(vsg_response)
		}

		/***************get Listener**********************/
		listener_request := slb.CreateDescribeLoadBalancerTCPListenerAttributeRequest()
		//listener_request.Scheme = "https"

		listener_request.ListenerPort = requests.NewInteger(int(port.Port))
		listener_request.LoadBalancerId = lbsId

		_, err = s.c.DescribeLoadBalancerTCPListenerAttribute(listener_request)
		if err != nil {
			fmt.Print(err.Error())
			if clientError, ok := err.(*errors.ClientError); ok {
				if clientError.ErrorCode() == "ListenerNotFound" {
					hadListener = false
				}
			}
		} else {
			hadListener = true
		}

		//if don't  have TCPListener to create TCPListener
		if !hadListener {
			/*Create TCPListener of LoadBalancer*/
			create_tcp_request := slb.CreateCreateLoadBalancerTCPListenerRequest()
			//create_tcp_request.Scheme = "https"

			create_tcp_request.VServerGroupId = vsgId
			create_tcp_request.ListenerPort = requests.NewInteger(int(port.Port))
			create_tcp_request.Bandwidth = requests.NewInteger(-1)
			create_tcp_request.LoadBalancerId = lbsId
			create_tcp_request.BackendServerPort = requests.NewInteger(int(port.NodePort))

			tcp_response, err := s.c.CreateLoadBalancerTCPListener(create_tcp_request)
			if err != nil {
				fmt.Print(err.Error())
				return err
			}
			fmt.Println(tcp_response)
		}
	}
	start_request := slb.CreateSetLoadBalancerStatusRequest()
	//start_request.Scheme = "https"
	start_request.LoadBalancerStatus = "active"
	start_request.LoadBalancerId = lbsId
	start_response, err := s.c.SetLoadBalancerStatus(start_request)
	if err != nil {
		fmt.Print(err.Error())
	}
	fmt.Println(start_response)
	return nil
}
