package provider

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/severalnines/clustercontrol-client-sdk/go/pkg/openapi"
	"log/slog"
)

type MySQLMaria struct {
	Common DbCommon
	Backup DbBackupCommon
}

func (m *MySQLMaria) GetInputs(d *schema.ResourceData, jobData *openapi.JobsJobJobSpecJobData) error {
	funcName := "GetInputs::GetInputs"
	slog.Info(funcName)
	//fmt.Fprintf(os.Stderr, "%s", funcName)

	var err error

	// parent/super - get common attributes
	if err = m.Common.GetInputs(d, jobData); err != nil {
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

func (c *MySQLMaria) HandleRead(ctx context.Context, d *schema.ResourceData, m interface{}, clusterInfo *openapi.ClusterResponse) error {

	if err := c.Common.HandleRead(ctx, d, m, clusterInfo); err != nil {
		return err
	}

	return nil
}

func (c *MySQLMaria) HandleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}, clusterInfo *openapi.ClusterResponse) error {

	if err := c.Common.HandleUpdate(ctx, d, m, clusterInfo); err != nil {
		return err
	}

	isHostChanged := false
	if d.HasChange(TF_FIELD_CLUSTER_HOST) {
		isHostChanged = true

		// Could be one of a few scenarios...

		// Scenario #1 - A new replica is to be added

		// Scenario #2 - A new replica is to be removed

		clusterType := d.Get(TF_FIELD_CLUSTER_TYPE)
		if clusterType == CLUSTER_TYPE_GALERA {

		} else {

		}

		//var clusterNameTf string
		//clusterNameTf = d.Get(TF_FIELD_CLUSTER_NAME).(string)
		//var tags = openapi.ClustersConfigurationInner{
		//	Name:  &CMON_CLUSTERS_OPERATION_SET_NAME,
		//	Value: &clusterNameTf,
		//}
		//configChanges = append(configChanges, tags)
		//isConfigUpdated = true

	}

	if d.HasChange(TF_FIELD_CLUSTER_TOPOLOGY) {
		// Could be one of a few scenarios...

		// Scenario #2 - Started with 1 node; Then, added a replica and specified the topology for future purposes (Scenario #1)
		if isHostChanged {
			// Noop: Just ignore. Nothing to do. This could possibly be on e of the following:
			// A. The addition of a new replica and the user has specified proper Master=>Slave links
			// B. The removal of a replica and the user has updated (or, if no replicas are present any longer, then completely removed) the topology def
		} else {
			// Scenario #1 - Role change: Slave promoted to Master
		}
	}

	return nil
}

func (m *MySQLMaria) GetBackupInputs(d *schema.ResourceData, jobData *openapi.JobsJobJobSpecJobData) error {
	funcName := "MySQL_Maria::GetBackupInputs"
	slog.Info(funcName)

	var err error

	// parent/super - get common attributes
	if err = m.Backup.GetBackupInputs(d, jobData); err != nil {
		return err
	}

	return err
}

func NewMySQLMaria() *MySQLMaria {
	return &MySQLMaria{}
}
