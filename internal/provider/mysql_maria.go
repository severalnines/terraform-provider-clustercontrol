package provider

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/severalnines/clustercontrol-client-sdk/go/pkg/openapi"
	"log/slog"
)

type MySQLMaria struct {
	common DbCommon
}

func (m *MySQLMaria) GetInputs(d *schema.ResourceData, jobData *openapi.JobsJobJobSpecJobData) error {
	funcName := "GetInputs::GetInputs"
	slog.Info(funcName)
	//fmt.Fprintf(os.Stderr, "%s", funcName)

	var err error

	// parent/super - get common attributes
	if err = m.common.GetInputs(d, jobData); err != nil {
		return err
	}

	clusterType := jobData.GetClusterType()
	dbVendor := jobData.GetVendor()
	dbVersion := jobData.GetVersion()
	iPort := int(jobData.GetPort())

	vendorMap, ok := gDbConfigTemplate[dbVendor]
	if !ok {
		return errors.New(fmt.Sprintf("Map doesn't support DB vendor: %s", dbVendor))
	}
	clusterTypeMap, ok := vendorMap[clusterType]
	if !ok {
		return errors.New(fmt.Sprintf("Map doesn't support DB vendor: %s, ClusterType: %s", dbVendor, clusterTypeMap))
	}
	cfgTemplate, ok := clusterTypeMap[dbVersion]
	if !ok {
		return errors.New(fmt.Sprintf("Map doesn't support DB vendor: %s, ClusterType: %s, DbVersion: %s", dbVendor, clusterTypeMap, dbVendor))
	}
	// TODO remove commented code later
	//configTemplate := d.Get("db_config_template").(string)
	//CheckForEmptyAndSet(&configTemplate, cfgTemplate)
	jobData.SetConfigTemplate(cfgTemplate)

	dataDirectory := d.Get(TF_FIELD_CLUSTER_DATA_DIR).(string)
	if err = CheckForEmptyAndSetDefault(&dataDirectory, gDefultDataDir, clusterType); err != nil {
		return err
	}
	jobData.SetDataDir(dataDirectory)

	semiSyncReplication := d.Get(TF_FIELD_CLUSTER_SEMISYNC_REP).(bool)
	jobData.SetMysqlSemiSync(semiSyncReplication)

	hosts := d.Get(TF_FIELD_CLUSTER_HOST)
	nodes := []openapi.JobsJobJobSpecJobDataNodesInner{}
	for _, ff := range hosts.([]any) {
		f := ff.(map[string]any)
		var node = openapi.JobsJobJobSpecJobDataNodesInner{}
		var memHost = memberHosts{
			vanillaNode: &node,
		}
		//var node = openapi.JobsJobJobSpecJobDataNodesInner{}
		getCommonHostAttributes(f, iPort, clusterType, memHost)
		nodes = append(nodes, node)
	}
	jobData.SetNodes(nodes)

	topology := d.Get(TF_FIELD_CLUSTER_TOPOLOGY)
	ccTopo := openapi.JobsJobJobSpecJobDataTopology{}
	for _, ff := range topology.([]any) {
		f := ff.(map[string]any)

		primary := f[TF_FIELD_CLUSTER_PRIMARY].(string)
		replica := f[TF_FIELD_CLUSTER_REPLICA].(string)

		slog.Debug(funcName, TF_FIELD_CLUSTER_PRIMARY, primary, TF_FIELD_CLUSTER_REPLICA, replica)

		var msLink = map[string]string{
			primary: replica,
		}
		ccTopo.MasterSlaveLinks = append(ccTopo.MasterSlaveLinks, msLink)
	}
	jobData.SetTopology(ccTopo)

	return nil
}

func NewMySQLMaria() *MySQLMaria {
	return &MySQLMaria{}
}
