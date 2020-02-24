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
)

func NewMockClientInstanceMgr() (*ClientMgr, error) {

	keyid := "xxxxxxxxxxx"
	keysecret := "xxxxxxxxxxxxxxxxxx"
	regionid := "cn-hangzhou"

	mgr, err := NewClientMgr(regionid, keyid, keysecret)
	if err != nil {
		return nil, err
	}
	return mgr, nil
}

// ======================================= This begins the TESTS ============================================

func TestInstanceRefeshInstance(t *testing.T) {

	mgr, err := NewMockClientInstanceMgr()
	if err != nil {
		t.Fatal(fmt.Sprintf("create client manager fail. [%s]\n", err.Error()))
	}

	instanceid := "i-xxxxxxxxxxxxxxxxxxxx"
	providerid := "cn-hangzhou.i-xxxxxxxxxxxxxxxxxxxxxxx"

	//test getInstance
	ins, err := mgr.Instances().getInstance(instanceid)
	if err != nil {
		t.Errorf("TestInstanceRefeshInstance error: %s\n", err.Error())
	}
	fmt.Println(ins)

	//test findInstanceByProviderID
	insa, err := mgr.Instances().findInstanceByProviderID(providerid)
	if err != nil {
		t.Fatal(fmt.Sprintf("findInstanceByNode error: %s\n", err.Error()))
	}
	fmt.Println(insa)

	//test findAddressByProviderID
	ips, err := mgr.Instances().findAddressByProviderID(providerid)
	if err != nil {
		t.Fatal(fmt.Sprintf("findAddressByProviderID error: %s\n", err.Error()))
	}
	fmt.Println(ips)

	providerid = "cn-hangzhou.i-xxxxxxxxxxxxxxxxxxx"
	insa, err = mgr.Instances().findInstanceByProviderID(providerid)
	if err != nil {
		t.Fatal(fmt.Sprintf("findInstanceByNode error: %s\n", err.Error()))
	}
	fmt.Println(insa)
}
