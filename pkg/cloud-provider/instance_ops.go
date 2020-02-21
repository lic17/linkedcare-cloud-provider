package cloud_provider

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
)

type InstanceClient struct {
	c *ecs.Client
}
