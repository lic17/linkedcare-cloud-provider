/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cloud_provider

import (
	"fmt"
	"strconv"

	"github.com/golang/glog"
	v1 "k8s.io/api/core/v1"
)

const (
	// ServiceAnnotationLoadBalancerPrefix is the annotation prefix of LoadBalancer
	ServiceAnnotationLoadBalancerPrefix = "service.beta.kubernetes.io/linkedcare-load-balancer-"

	ServiceAnnotationLoadBalancerId = ServiceAnnotationLoadBalancerPrefix + "id"
)

const (
	// NodeAnnotationPrefix is the annotation prefix of Node
	NodeAnnotationPrefix = "node.alpha.kubernetes.io/"
	// NodeAnnotationVpcId is the annotation of VpcId on node
	NodeAnnotationVpcId = NodeAnnotationPrefix + "vpc-id"
	// NodeAnnotationVpcRouteTableId is the annotation of VpcRouteTableId on node
	NodeAnnotationVpcRouteTableId = NodeAnnotationPrefix + "vpc-route-table-id"
	// NodeAnnotationVpcRouteRuleId is the annotation of VpcRouteRuleId on node
	NodeAnnotationVpcRouteRuleId = NodeAnnotationPrefix + "vpc-route-rule-id"

	// NodeAnnotationCCMVersion is the version of CCM
	NodeAnnotationCCMVersion = NodeAnnotationPrefix + "ccm-version"

	// NodeAnnotationAdvertiseRoute indicates whether to advertise route to vpc route table
	NodeAnnotationAdvertiseRoute = NodeAnnotationPrefix + "advertise-route"
)

// ServiceAnnotation contains annotations from service
type ServiceAnnotation struct {
	/* BLB */
	Loadbalancerid string
}

// NodeAnnotation contains annotations from node
type NodeAnnotation struct {
	VpcId           string
	VpcRouteTableId string
	VpcRouteRuleId  string
	CCMVersion      string
	AdvertiseRoute  bool
}

// ExtractServiceAnnotation extract annotations from service
func ExtractServiceAnnotation(service *v1.Service) (*ServiceAnnotation, error) {
	glog.V(4).Infof("start to ExtractServiceAnnotation: %v", service.Annotations)
	result := &ServiceAnnotation{}
	annotation := make(map[string]string)
	for k, v := range service.Annotations {
		annotation[k] = v
	}

	loadBalancerId, exist := annotation[ServiceAnnotationLoadBalancerId]
	if exist {
		result.Loadbalancerid = loadBalancerId
	}

	return result, nil
}

// ExtractNodeAnnotation extract annotations from node
func ExtractNodeAnnotation(node *v1.Node) (*NodeAnnotation, error) {
	glog.V(4).Infof("start to ExtractNodeAnnotation: %v", node.Annotations)
	result := &NodeAnnotation{}
	annotation := make(map[string]string)
	for k, v := range node.Annotations {
		annotation[k] = v
	}

	vpcId, ok := annotation[NodeAnnotationVpcId]
	if ok {
		result.VpcId = vpcId
	}

	vpcRouteTableId, ok := annotation[NodeAnnotationVpcRouteTableId]
	if ok {
		result.VpcRouteTableId = vpcRouteTableId
	}

	vpcRouteRuleId, ok := annotation[NodeAnnotationVpcRouteRuleId]
	if ok {
		result.VpcRouteRuleId = vpcRouteRuleId
	}

	ccmVersion, ok := annotation[NodeAnnotationCCMVersion]
	if ok {
		result.CCMVersion = ccmVersion
	}

	advertiseRoute, ok := annotation[NodeAnnotationAdvertiseRoute]
	if ok {
		advertise, err := strconv.ParseBool(advertiseRoute)
		if err != nil {
			return nil, fmt.Errorf("NodeAnnotationAdvertiseRoute syntex error: %v", err)
		}
		result.AdvertiseRoute = advertise
	} else {
		result.AdvertiseRoute = true
	}

	return result, nil
}

func serviceAnnotation(service *v1.Service, annotate string) string {
	for k, v := range service.Annotations {
		if annotate == k {
			return v
		}
	}
	return ""
}
