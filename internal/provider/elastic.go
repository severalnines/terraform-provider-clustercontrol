package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/severalnines/clustercontrol-client-sdk/go/pkg/openapi"
	"log/slog"
)

type Elastic struct {
	Common DbCommon
	Backup DbBackupCommon
}

func (m *Elastic) GetInputs(d *schema.ResourceData, jobData *openapi.JobsJobJobSpecJobData) error {
	funcName := "Elastic::GetInputs"
	slog.Debug(funcName)

	var err error

	// parent/super - get common attributes
	if err = m.Common.GetInputs(d, jobData); err != nil {
		return err
	}

	snapshotLocation := d.Get(TF_FIELD_CLUSTER_SNAPSHOT_LOC).(string)
	jobData.SetSnapshotLocation(snapshotLocation)

	snapshotRepo := d.Get(TF_FIELD_CLUSTER_SNAPSHOT_REPO).(string)
	jobData.SetSnapshotRepository(snapshotRepo)

	//snapshotHost := d.Get(TF_FIELD_CLUSTER_SNAPSHOT_HOST).(string)
	//jobData.SetSn(snapshotHost)

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
		protocol := f[TF_FIELD_CLUSTER_HOST_PROTO].(string)
		roles := f[TF_FIELD_CLUSTER_HOST_ROLES].(string)
		node.SetClassName(CMON_CLASS_NAME_ELASTIC_HOST)
		node.SetProtocol(protocol)
		node.SetRoles(roles)

		nodes = append(nodes, node)
	}
	jobData.SetNodes(nodes)

	return nil
}

func (c *Elastic) HandleRead(ctx context.Context, d *schema.ResourceData, m interface{}, clusterInfo *openapi.ClusterResponse) error {

	if err := c.Common.HandleRead(ctx, d, m, clusterInfo); err != nil {
		return err
	}

	return nil
}

func (c *Elastic) HandleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}, clusterInfo *openapi.ClusterResponse) error {

	if err := c.Common.HandleUpdate(ctx, d, m, clusterInfo); err != nil {
		return err
	}

	return nil
}

func (c *Elastic) GetBackupInputs(d *schema.ResourceData, jobData *openapi.JobsJobJobSpecJobData) error {
	funcName := "Elastic::GetBackupInputs"
	slog.Info(funcName)

	var err error

	// parent/super - get common attributes
	if err = c.Backup.GetBackupInputs(d, jobData); err != nil {
		return err
	}

	return err
}

func NewElastic() *Elastic {
	return &Elastic{}
}
