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
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/golang/glog"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	v1core "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/record"
	"k8s.io/kubernetes/pkg/cloudprovider"
	"k8s.io/kubernetes/pkg/controller"
	"k8s.io/kubernetes/pkg/version"
)

// ProviderName is the name of this cloud provider.
const ProviderName = "alicloud"

// CLUSTER_ID default cluster id if it is not specified.
var CLUSTER_ID = "clusterid"

// KUBERNETES_ALICLOUD_IDENTITY is for statistic purpose.
var KUBERNETES_ALICLOUD_IDENTITY = fmt.Sprintf("Kubernetes.Alicloud/%s", version.Get().String())

// Cloud defines the main struct
type Cloud struct {
	climgr           *ClientMgr
	cfg              *CloudConfig
	kubeClient       kubernetes.Interface
	eventBroadcaster record.EventBroadcaster
	eventRecorder    record.EventRecorder
}

var (
	// DEFAULT_CHARGE_TYPE default charge type

	// DEFAULT_BANDWIDTH default bandwidth
	DEFAULT_BANDWIDTH = 100

	DEFAULT_NODE_MONITOR_PERIOD = 120 * time.Second

	DEFAULT_NODE_ADDR_SYNC_PERIOD = 240 * time.Second

	// DEFAULT_REGION should be override in cloud initialize.
	//DEFAULT_REGION = common.Hangzhou
)

// CloudConfig is the cloud config
type CloudConfig struct {
	UID             string `json:"uid"`
	ClusterID       string `json:"ClusterId"`
	ClusterName     string `json:"ClusterName"`
	AccessKeyID     string `json:"AccessKeyID"`
	AccessKeySecret string `json:"AccessKeySecret"`
	Region          string `json:"Region"`
	VpcID           string `json:"VpcId"`
	SubnetID        string `json:"SubnetId"`
	MasterID        string `json:"MasterId"`
	Endpoint        string `json:"Endpoint"`
	NodeIP          string `json:"NodeIP"`
	Debug           bool   `json:"Debug"`
}

// CCMVersion is the version of CCM
var CCMVersion string
var cfg CloudConfig

func init() {
	cloudprovider.RegisterCloudProvider(ProviderName,
		func(config io.Reader) (cloudprovider.Interface, error) {
			var (
				keyid     = ""
				keysecret = ""
				regionid  = ""
			)
			if config != nil {
				if err := json.NewDecoder(config).Decode(&cfg); err != nil {
					return nil, err
				}
				if cfg.AccessKeyID != "" && cfg.AccessKeySecret != "" && cfg.Region != "" {
					key, err := b64.StdEncoding.DecodeString(cfg.AccessKeyID)
					if err != nil {
						return nil, err
					}
					keyid = string(key)
					secret, err := b64.StdEncoding.DecodeString(cfg.AccessKeySecret)
					if err != nil {
						return nil, err
					}
					keysecret = string(secret)
					region, err := b64.StdEncoding.DecodeString(cfg.Region)
					if err != nil {
						return nil, err
					}
					regionid = string(region)
					glog.V(2).Infof("Alicloud: Try Accesskey AccessKeySecret and Region from config file.")
				}
				if cfg.ClusterID != "" {
					CLUSTER_ID = cfg.ClusterID
					glog.Infof("use clusterid %s", CLUSTER_ID)
				}
			}
			if keyid == "" || keysecret == "" {
				glog.V(2).Infof("cloud config does not have keyid and keysecret . try environment ACCESS_KEY_ID ACCESS_KEY_SECRET REGION_ID")
				keyid = os.Getenv("ACCESS_KEY_ID")
				keysecret = os.Getenv("ACCESS_KEY_SECRET")
				regionid = os.Getenv("REGION_ID")
			}
			mgr, err := NewClientMgr(regionid, keyid, keysecret)
			if err != nil {
				return nil, err
			}
			// wait for client initialized
			//err = mgr.Start(RefreshToken)
			//f err != nil {
			//	panic(fmt.Sprintf("token not ready %s", err.Error()))
			//}
			fmt.Println(mgr)
			return newAliCloud(mgr)
		})

}

func newAliCloud(mgr *ClientMgr) (*Cloud, error) {

	return &Cloud{
		climgr: mgr,
		cfg:    &cfg,
	}, nil
}

// Initialize provides the cloud with a kubernetes client builder and may spawn goroutines
// to perform housekeeping activities within the cloud provider.
func (c *Cloud) Initialize(builder controller.ControllerClientBuilder) {
	c.kubeClient = builder.ClientOrDie(ProviderName)
	c.eventBroadcaster = record.NewBroadcaster()
	c.eventBroadcaster.StartLogging(glog.Infof)
	c.eventBroadcaster.StartRecordingToSink(&v1core.EventSinkImpl{Interface: c.kubeClient.CoreV1().Events("")})
	c.eventRecorder = c.eventBroadcaster.NewRecorder(scheme.Scheme, v1.EventSource{Component: "CCM"})
}
func (c *Cloud) ProviderName() string {
	return ProviderName
}
