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
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func NewMockClientLoadBalancerMgr() (*ClientMgr, error) {

	keyid := "XXXXXXXXXXXXXXX"
	keysecret := "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
	regionid := "cn-hangzhou"

	mgr, err := NewClientMgr(regionid, keyid, keysecret)
	if err != nil {
		return nil, err
	}
	return mgr, nil
}

// ======================================= This begins the TESTS ============================================

func TestLoadBalancer(t *testing.T) {

	mgr, err := NewMockClientLoadBalancerMgr()
	if err != nil {
		t.Fatal(fmt.Sprintf("create client manager fail. [%s]\n", err.Error()))
	}

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
	prid := "cn-hangzhou.i-bp15ekjuuvrwuxjowxcc"
	node := []*v1.Node{
		{
			ObjectMeta: metav1.ObjectMeta{Name: prid},
			Spec: v1.NodeSpec{
				ProviderID: prid,
			},
		},
	}
	//test findLoadBalancer
	exists, lbs, err := mgr.LoadBalancers().findLoadBalancer(service)
	if err != nil {
		t.Errorf("findLoadBalancer error: %s\n", err.Error())
	}
	if exists {
		fmt.Println("findLoadBalancer", lbs)
	} else {
		fmt.Println("findLoadBalancer: no loadbalancer")
	}

	def, _ := ExtractServiceAnnotation(service)
	fmt.Println(def)
	fmt.Println(def.Loadbalancerid)

	//test findLoadBalancerByID
	lbid := "lb-bp1umlml75qdkig2ggf5y"
	exists, lbs, err = mgr.LoadBalancers().findLoadBalancerByID(lbid)
	if err != nil {
		t.Errorf("findLoadBalancer by id error: %s\n", err.Error())
	}
	if exists {
		fmt.Println("findLoadBalancer by id :", lbs)
	} else {
		fmt.Println("findLoadBalancer by id: no loadbalancer")
	}

	//test ensureLoadBalancer
	exists, lbs, err = mgr.LoadBalancers().ensureLoadBalancer(service, node)
	if err != nil {
		t.Errorf("ensureLoadBalancer error: %s\n", err.Error())
	}
	if exists {
		fmt.Println("ensureLoadBalancer:", lbs)
	} else {
		fmt.Println("ensureLoadBalancer: no loadbalancer")
	}
}
