package cloud_provider

import (
	"context"

	"k8s.io/apimachinery/pkg/types"
	"k8s.io/kubernetes/pkg/cloudprovider"
)

/*****************routes*******************/
func (c *Cloud) Routes() (cloudprovider.Routes, bool) {
	return c, true
}
func (c *Cloud) ListRoutes(ctx context.Context, clusterName string) ([]*cloudprovider.Route, error) {
	var routes []*cloudprovider.Route
	return routes, nil
}

func (c *Cloud) CreateRoute(ctx context.Context, clusterName string, nameHint string, route *cloudprovider.Route) error {
	return nil
}
func (c *Cloud) DeleteRoute(ctx context.Context, clusterName string, route *cloudprovider.Route) error {
	return nil
}

/*****************zones*******************/
func (c *Cloud) Zones() (cloudprovider.Zones, bool) {
	return c, true
}
func (c *Cloud) GetZone(ctx context.Context) (cloudprovider.Zone, error) {
	var zone cloudprovider.Zone
	return zone, nil
}
func (c *Cloud) GetZoneByProviderID(ctx context.Context, providerID string) (cloudprovider.Zone, error) {
	var zone cloudprovider.Zone
	return zone, nil
}
func (c *Cloud) GetZoneByNodeName(ctx context.Context, nodeName types.NodeName) (cloudprovider.Zone, error) {
	var zone cloudprovider.Zone
	return zone, nil
}
