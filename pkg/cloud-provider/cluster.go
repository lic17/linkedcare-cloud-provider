package cloud_provider

import (
	"context"
	"fmt"

	"k8s.io/kubernetes/pkg/cloudprovider"
)

// Clusters returns a clusters interface.  Also returns true if the interface is supported, false otherwise.
func (bc *Cloud) Clusters() (cloudprovider.Clusters, bool) {
	return nil, false
}

// HasClusterID returns true if a ClusterID is required and set
func (bc *Cloud) HasClusterID() bool {
	return true
}

// ListClusters lists the names of the available clusters.
func (bc *Cloud) ListClusters(ctx context.Context) ([]string, error) {
	return nil, fmt.Errorf("ListClusters unimplemented")
}

// Master gets back the address (either DNS name or IP address) of the master node for the cluster.
func (bc *Cloud) Master(ctx context.Context, clusterName string) (string, error) {
	return "", fmt.Errorf("Master unimplemented")
}
