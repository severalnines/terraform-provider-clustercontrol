package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/severalnines/clustercontrol-client-sdk/go/pkg/openapi"
	"log/slog"
)

type Redis struct {
	Common DbCommon
	Backup DbBackupCommon
}

func (m *Redis) GetInputs(d *schema.ResourceData, jobData *openapi.JobsJobJobSpecJobData) error {
	funcName := "Redis::GetInputs"
	slog.Debug(funcName)

	var err error

	// parent/super - get common attributes
	if err = m.Common.GetInputs(d, jobData); err != nil {
		return err
	}

	clusterType := jobData.GetClusterType()
	iPort := int(jobData.GetPort())
	hosts := d.Get(TF_FIELD_CLUSTER_HOST)
	nodes := []openapi.JobsJobJobSpecJobDataNodesInner{}
	for _, ff := range hosts.([]any) {
		f := ff.(map[string]any)

		var node = openapi.JobsJobJobSpecJobDataNodesInner{}
		var memHost = memberHosts{
			vanillaNode: &node,
		}
		getCommonHostAttributes(f, iPort, clusterType, memHost)
		var node2 = node
		node.SetClassName(CMON_CLASS_NAME_REDIS_HOST)
		nodes = append(nodes, node)

		node2.SetClassName(CMON_CLASS_NAME_REDIS_SENTNEL_HOST)
		node2.SetPort(DEFAULT_MONGO_REDIS_SENTINEL_PORT)
		nodes = append(nodes, node2)
	}
	jobData.SetNodes(nodes)

	return nil
}

func (c *Redis) HandleRead(ctx context.Context, d *schema.ResourceData, m interface{}, clusterInfo *openapi.ClusterResponse) error {

	if err := c.Common.HandleRead(ctx, d, m, clusterInfo); err != nil {
		return err
	}

	return nil
}

func (c *Redis) HandleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}, clusterInfo *openapi.ClusterResponse) error {

	if err := c.Common.HandleUpdate(ctx, d, m, clusterInfo); err != nil {
		return err
	}

	return nil
}

func (c *Redis) GetBackupInputs(d *schema.ResourceData, jobData *openapi.JobsJobJobSpecJobData) error {
	funcName := "Redis::GetBackupInputs"
	slog.Info(funcName)

	var err error

	// parent/super - get common attributes
	if err = c.Backup.GetBackupInputs(d, jobData); err != nil {
		return err
	}

	return err
}

func NewRedis() *Redis {
	return &Redis{}
}
