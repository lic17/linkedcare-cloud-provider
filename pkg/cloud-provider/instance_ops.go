package cloud_provider

import (
	"github.com/denverdino/aliyungo/ecs"
)

type InstanceClient struct {
	c *ecs.Client
}
