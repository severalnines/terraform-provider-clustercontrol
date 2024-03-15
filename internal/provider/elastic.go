package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/severalnines/clustercontrol-client-sdk/go/pkg/openapi"
	"log/slog"
)

type Elastic struct {
	common DbCommon
}

func (m *Elastic) GetInputs(d *schema.ResourceData, jobData *openapi.JobsJobJobSpecJobData) error {
	funcName := "Elastic::GetInputs"
	slog.Debug(funcName)

	var err error

	// parent/super - get common attributes
	if err = m.common.GetInputs(d, jobData); err != nil {
		return err
	}

	snapshotLocation := d.Get(TF_FIELD_CLUSTER_SNAPSHOT_LOC).(string)
	jobData.SetSnapshotLocaiton(snapshotLocation)

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

func NewElastic() *Elastic {
	return &Elastic{}
}
