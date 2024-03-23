package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/severalnines/clustercontrol-client-sdk/go/pkg/openapi"
	"log/slog"
)

type PostgresSql struct {
	Common DbCommon
	Backup DbBackupCommon
}

func (m *PostgresSql) GetInputs(d *schema.ResourceData, jobData *openapi.JobsJobJobSpecJobData) error {
	funcName := "Postgres::GetInputs"
	slog.Debug(funcName)

	var err error

	// parent/super - get common attributes
	if err = m.Common.GetInputs(d, jobData); err != nil {
		return err
	}

	iPort := int(jobData.GetPort())
	clusterType := jobData.GetClusterType()
	hosts := d.Get(TF_FIELD_CLUSTER_HOST)
	nodes := []openapi.JobsJobJobSpecJobDataNodesInner{}
	for _, ff := range hosts.([]any) {
		f := ff.(map[string]any)

		var node = openapi.JobsJobJobSpecJobDataNodesInner{}
		var memHost = memberHosts{
			vanillaNode: &node,
		}
		getCommonHostAttributes(f, iPort, clusterType, memHost)
		sync_replication := f[TF_FIELD_CLUSTER_SYNC_REP].(bool)
		node.SetSynchronous(sync_replication)

		slog.Debug(funcName, TF_FIELD_CLUSTER_SYNC_REP, sync_replication)

		nodes = append(nodes, node)
	}
	jobData.SetNodes(nodes)

	return nil
}

func (c *PostgresSql) HandleRead(ctx context.Context, d *schema.ResourceData, m interface{}, clusterInfo *openapi.ClusterResponse) error {

	if err := c.Common.HandleRead(ctx, d, m, clusterInfo); err != nil {
		return err
	}

	return nil
}

func (c *PostgresSql) HandleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}, clusterInfo *openapi.ClusterResponse) error {

	if err := c.Common.HandleUpdate(ctx, d, m, clusterInfo); err != nil {
		return err
	}

	return nil
}

func (c *PostgresSql) GetBackupInputs(d *schema.ResourceData, jobData *openapi.JobsJobJobSpecJobData) error {
	funcName := "PG::GetBackupInputs"
	slog.Info(funcName)

	var err error

	// parent/super - get common attributes
	if err = c.Backup.GetBackupInputs(d, jobData); err != nil {
		return err
	}

	return err
}

func NewPostgres() *PostgresSql {
	return &PostgresSql{}
}
