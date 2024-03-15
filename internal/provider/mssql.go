package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/severalnines/clustercontrol-client-sdk/go/pkg/openapi"
	"log/slog"
)

type MsSql struct {
	common DbCommon
}

func (m *MsSql) GetInputs(d *schema.ResourceData, jobData *openapi.JobsJobJobSpecJobData) error {
	funcName := "MsSql::GetInputs"
	slog.Debug(funcName)

	var err error

	// parent/super - get common attributes
	if err = m.common.GetInputs(d, jobData); err != nil {
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
		node.SetClassName(CMON_CLASS_NAME_MSSQL_HOST)
		dataDir := gDefultDataDir[clusterType]
		configFile := gDefaultHostConfigFile[clusterType]
		node.SetDatadir(dataDir)
		node.SetConfigfile(configFile)

		nodes = append(nodes, node)
	}
	jobData.SetNodes(nodes)

	return nil
}

func NewMsSql() *MsSql {
	return &MsSql{}
}
