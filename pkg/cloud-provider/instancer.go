package cloud_provider

import (
	"context"
	"errors"

	"github.com/golang/glog"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/kubernetes/pkg/cloudprovider"
)

func (c *Cloud) Instances() (cloudprovider.Instances, bool) {
	return c, true
}

func (c *Cloud) NodeAddresses(ctx context.Context, name types.NodeName) ([]v1.NodeAddress, error) {
	glog.V(2).Infof("Alicloud.NodeAddresses(\"%s\")", name)
	return c.climgr.Instances().findAddressByNodeName(name)
}

func (c *Cloud) NodeAddressesByProviderID(ctx context.Context, providerID string) ([]v1.NodeAddress, error) {
	return c.climgr.Instances().findAddressByProviderID(providerID)
}

func (c *Cloud) InstanceExistsByProviderID(ctx context.Context, providerID string) (bool, error) {
	_, err := c.climgr.Instances().findInstanceByProviderID(providerID)
	if err != nil {
		return false, err
	}
	return true, nil
}

//Doing
func (c *Cloud) CurrentNodeName(ctx context.Context, hostname string) (types.NodeName, error) {
	var name types.NodeName
	return name, nil
}

func (c *Cloud) InstanceTypeByProviderID(ctx context.Context, providerID string) (string, error) {
	ins, err := c.climgr.Instances().findInstanceByProviderID(providerID)
	if err == nil {
		return ins.InstanceType, nil
	}
	return "", err
}
func (c *Cloud) InstanceType(ctx context.Context, name types.NodeName) (string, error) {

	ins, err := c.climgr.Instances().findInstanceByNodeName(name)
	if err == nil {
		return ins.InstanceType, nil
	}
	return "", err
}
func (c *Cloud) InstanceID(ctx context.Context, nodeName types.NodeName) (string, error) {
	ins, err := c.climgr.Instances().findInstanceByNodeName(nodeName)
	if err == nil {
		return ins.InstanceId, nil
	}
	return "", err
}

func (c *Cloud) AddSSHKeyToAllInstances(ctx context.Context, user string, keyData []byte) error {
	return errors.New("Alicloud.AddSSHKeyToAllInstances() is not implemented")
}

// TODO
func (c *Cloud) InstanceShutdownByProviderID(ctx context.Context, providerID string) (bool, error) {
	return false, nil
}
