package provider

import (
	"context"
	"errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/severalnines/clustercontrol-client-sdk/go/pkg/openapi"
	"log/slog"
	"strconv"
)

type MsSql struct {
	Common DbCommon
	Backup DbBackupCommon
}

func (m *MsSql) GetInputs(d *schema.ResourceData, jobData *openapi.JobsJobJobSpecJobData) error {
	funcName := "MsSql::GetInputs"
	slog.Debug(funcName)

	var err error

	// parent/super - get common attributes
	if err = m.Common.GetInputs(d, jobData); err != nil {
		return err
	}

	clusterType := jobData.GetClusterType()

	var iPort int
	port := d.Get(TF_FIELD_CLUSTER_MSSQL_SERVER_PORT).(string)
	if err = CheckForEmptyAndSetDefault(&port, gDefultDbPortMap, clusterType); err != nil {
		return err
	}
	if iPort, err = strconv.Atoi(port); err != nil {
		slog.Error(funcName, "ERROR", "Non-numeric database port")
		return err
	}
	jobData.SetPort(int32(iPort))

	topLevelPort := jobData.GetPort()

	hosts := d.Get(TF_FIELD_CLUSTER_HOST)
	nodes := []openapi.JobsJobJobSpecJobDataNodesInner{}
	numHosts := 0
	for _, ff := range hosts.([]any) {
		f := ff.(map[string]any)
		hostname := f[TF_FIELD_CLUSTER_HOSTNAME].(string)
		hostname_data := f[TF_FIELD_CLUSTER_HOSTNAME_DATA].(string)
		hostname_internal := f[TF_FIELD_CLUSTER_HOSTNAME_INT].(string)
		//port := f[TF_FIELD_CLUSTER_HOST_PORT].(string)
		datadir := f[TF_FIELD_CLUSTER_HOST_DD].(string)

		if hostname == "" {
			return errors.New("Hostname cannot be empty")
		}
		var node = openapi.JobsJobJobSpecJobDataNodesInner{
			Hostname: &hostname,
		}
		node.SetClassName(CMON_CLASS_NAME_MSSQL_HOST)
		if hostname_data != "" {
			node.SetHostnameData(hostname_data)
		}
		if hostname_internal != "" {
			node.SetHostnameInternal(hostname_internal)
		}
		if datadir != "" {
			node.SetDatadir(datadir)
		} else {
			dataDir := gDefultDataDir[clusterType]
			node.SetDatadir(dataDir)
		}
		node.SetPort(strconv.Itoa(int(topLevelPort)))
		//if port == "" {
		//	node.SetPort(strconv.Itoa(int(topLevelPort)))
		//} else {
		//	node.SetPort(strconv.Itoa(int(convertPortToInt(port, topLevelPort))))
		//}

		configFile := gDefaultHostConfigFile[clusterType]
		node.SetConfigfile(configFile)

		nodes = append(nodes, node)

		numHosts++
	}

	// For MSSQL, two internal cluster-types map to one external cluster-type. Fix it
	if numHosts > 1 {
		jobData.SetClusterType(CLUSTER_TYPE_MSSQL_AO_ASYNC)
	}

	jobData.SetNodes(nodes)

	return nil
}

func (c *MsSql) HandleRead(ctx context.Context, d *schema.ResourceData, apiClient *openapi.APIClient, clusterInfo *openapi.ClusterResponse) error {

	if err := c.Common.HandleRead(ctx, d, apiClient, clusterInfo); err != nil {
		return err
	}

	return nil
}

func (c *MsSql) IsUpdateBatchAllowed(d *schema.ResourceData) error {
	var err error

	if err = c.Common.IsUpdateBatchAllowed(d); err != nil {
		return err
	}

	return nil
}

func (c *MsSql) HandleUpdate(ctx context.Context, d *schema.ResourceData, apiClient *openapi.APIClient, clusterInfo *openapi.ClusterResponse) error {

	if err := c.Common.HandleUpdate(ctx, d, apiClient, clusterInfo); err != nil {
		return err
	}

	return nil
}

func (c *MsSql) GetBackupInputs(d *schema.ResourceData, jobData *openapi.JobsJobJobSpecJobData) error {
	funcName := "MsSql::GetBackupInputs"
	slog.Debug(funcName)

	var err error

	// parent/super - get common attributes
	if err = c.Backup.GetBackupInputs(d, jobData); err != nil {
		return err
	}

	jobData.SetHostname(STINRG_AUTO)

	backupSystemDb := d.Get(TF_FIELD_BACKUP_SYSTEM_DB).(bool)
	jobData.SetBackupSystemDb(backupSystemDb)

	return err
}

func (c *MsSql) IsValidBackupOptions(vendor string, clusterType string, jobData *openapi.JobsJobJobSpecJobData) error {
	return c.Backup.IsValidBackupOptions(vendor, clusterType, jobData)
}

func (c *MsSql) SetBackupJobData(jobData *openapi.JobsJobJobSpecJobData) error {
	return c.Backup.SetBackupJobData(jobData)
}

func (c *MsSql) IsBackupRemovable(clusterInfo *openapi.ClusterResponse, jobData *openapi.JobsJobJobSpecJobData) bool {
	return true
}

func NewMsSql() *MsSql {
	return &MsSql{}
}
