package cloud_provider

import (
	"fmt"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/golang/glog"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

type InstanceClient struct {
	c *ecs.Client
}

func (s *InstanceClient) findAddressByNodeName(nodeName types.NodeName) ([]v1.NodeAddress, error) {
	instance, err := s.findInstanceByNodeName(nodeName)
	if err != nil {
		glog.Errorf("alicloud: error getting instance by nodeName. providerID='%s', message=[%s]\n", nodeName, err.Error())
		return nil, err
	}
	return s.findAddressByInstance(instance), nil
}

func (s *InstanceClient) findAddressByInstance(instance *ecs.DescribeInstanceAttributeResponse) []v1.NodeAddress {
	addrs := []v1.NodeAddress{}

	if len(instance.PublicIpAddress.IpAddress) > 0 {
		for _, ipaddr := range instance.PublicIpAddress.IpAddress {
			addrs = append(addrs, v1.NodeAddress{Type: v1.NodeExternalIP, Address: ipaddr})
		}
	}

	if instance.EipAddress.IpAddress != "" {
		addrs = append(addrs, v1.NodeAddress{Type: v1.NodeExternalIP, Address: instance.EipAddress.IpAddress})
	}

	if len(instance.InnerIpAddress.IpAddress) > 0 {
		for _, ipaddr := range instance.InnerIpAddress.IpAddress {
			addrs = append(addrs, v1.NodeAddress{Type: v1.NodeInternalIP, Address: ipaddr})
		}
	}

	if len(instance.VpcAttributes.PrivateIpAddress.IpAddress) > 0 {
		for _, ipaddr := range instance.VpcAttributes.PrivateIpAddress.IpAddress {
			addrs = append(addrs, v1.NodeAddress{Type: v1.NodeInternalIP, Address: ipaddr})
		}
	}

	return addrs
}

// findAddressByProviderID returns an address slice by it's providerID.
func (s *InstanceClient) findAddressByProviderID(providerID string) ([]v1.NodeAddress, error) {

	instance, err := s.findInstanceByProviderID(providerID)
	if err != nil {
		glog.Errorf("alicloud: error getting instance by providerID. providerID='%s', message=[%s]\n", providerID, err.Error())
		return nil, err
	}

	return s.findAddressByInstance(instance), nil
}

func (s *InstanceClient) findInstanceByNodeName(nodeName types.NodeName) (*ecs.DescribeInstanceAttributeResponse, error) {
	return s.findInstanceByProviderID(string(nodeName))
}

func (s *InstanceClient) findInstanceByProviderID(providerID string) (*ecs.DescribeInstanceAttributeResponse, error) {
	nodeid, err := nodeFromProviderID(providerID)
	if err != nil {
		return nil, err
	}
	ins, err := s.getInstance(nodeid)
	if err != nil {
		glog.Errorf("alicloud: InstanceInspectError, instanceid=[%s]. message=[%s]\n", providerID, err.Error())

		return nil, err
	}

	return ins, nil
}

func (s *InstanceClient) getInstance(id string) (*ecs.DescribeInstanceAttributeResponse, error) {
	request := ecs.CreateDescribeInstanceAttributeRequest()
	//	request.Scheme = "https"
	request.InstanceId = id

	instance, err := s.c.DescribeInstanceAttribute(request)

	if err != nil {
		glog.Errorf("alicloud: calling DescribeInstances error. , "+
			"instanceid=%s, message=[%s].\n", id, err.Error())
		return nil, err
	}
	return instance, nil
}

// Use '.' to separate providerID which looks like 'cn-hangzhou.i-v98dklsmnxkkgiiil7'. The format of "REGION.NODEID"
func nodeFromProviderID(providerID string) (string, error) {
	name := strings.Split(providerID, ".")
	if len(name) < 2 {
		return "", fmt.Errorf("alicloud: unable to split instanceid and region from providerID, error unexpected providerID=%s", providerID)
	}
	return name[1], nil
}
