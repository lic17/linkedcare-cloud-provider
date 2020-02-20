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

package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/golang/glog"
	"k8s.io/apiserver/pkg/server/healthz"
	"k8s.io/kubernetes/cmd/cloud-controller-manager/app"
	_ "k8s.io/kubernetes/pkg/client/metrics/prometheus" // for client metric registration
	_ "k8s.io/kubernetes/pkg/features"                  // add the kubernetes feature gates
	_ "k8s.io/kubernetes/pkg/version/prometheus"        // for version metric registration

	cloud_provider "github.com/lic17/linkedcare-cloud-provider/pkg/cloud-provider"
)

var version string

func init() {
	healthz.DefaultHealthz()
	cloud_provider.CCMVersion = version
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	c := app.NewCloudControllerManagerCommand()

	glog.V(1).Infof("CCE Cloud-Controller-Manager version: %s", version)

	if err := c.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func GetVersion() string {
	return version
}
