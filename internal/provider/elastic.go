package provider

import (
	"context"
	"errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/severalnines/clustercontrol-client-sdk/go/pkg/openapi"
	"log/slog"
	"strconv"
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
	if snapshotLocation != "" {
		jobData.SetSnapshotLocation(snapshotLocation)
	}

	snapshotRepo := d.Get(TF_FIELD_CLUSTER_SNAPSHOT_REPO).(string)
	if snapshotRepo != "" {
		jobData.SetSnapshotRepository(snapshotRepo)
	}

	snapshotRepoType := d.Get(TF_FIELD_CLUSTER_SNAPSHOT_REPO_TYPE).(string)
	if snapshotRepoType != "" {
		jobData.SetSnapshotRepositoryType(snapshotRepoType)
	}

	storageHost := d.Get(TF_FIELD_CLUSTER_STORAGE_HOST).(string)
	if storageHost != "" {
		jobData.SetStorageHost(storageHost)
	}

	clusterType := jobData.GetClusterType()

	var iPort int
	port := d.Get(TF_FIELD_CLUSTER_ELASTIC_HTTP_PORT).(string)
	if err = CheckForEmptyAndSetDefault(&port, gDefultDbPortMap, clusterType); err != nil {
		return err
	}
	if iPort, err = strconv.Atoi(port); err != nil {
		slog.Error(funcName, "ERROR", "Non-numeric database port")
		return err
	}
	jobData.SetPort(int32(iPort))
	// topLevelPort := jobData.GetPort()

	//clusterType := jobData.GetClusterType()
	hosts := d.Get(TF_FIELD_CLUSTER_HOST)
	nodes := []openapi.JobsJobJobSpecJobDataNodesInner{}
	for _, ff := range hosts.([]any) {
		f := ff.(map[string]any)
		hostname := f[TF_FIELD_CLUSTER_HOSTNAME].(string)
		hostname_data := f[TF_FIELD_CLUSTER_HOSTNAME_DATA].(string)
		hostname_internal := f[TF_FIELD_CLUSTER_HOSTNAME_INT].(string)
		//port := f[TF_FIELD_CLUSTER_HOST_PORT].(string)
		protocol := f[TF_FIELD_CLUSTER_HOST_PROTO].(string)
		roles := f[TF_FIELD_CLUSTER_HOST_ROLES].(string)

		if hostname == "" {
			return errors.New("Hostname cannot be empty")
		}
		var node = openapi.JobsJobJobSpecJobDataNodesInner{
			Hostname: &hostname,
		}
		node.SetClassName(CMON_CLASS_NAME_ELASTIC_HOST)
		if hostname_data != "" {
			node.SetHostnameData(hostname_data)
		}
		if hostname_internal != "" {
			node.SetHostnameInternal(hostname_internal)
		}
		// node.SetPort(strconv.Itoa(int(topLevelPort)))
		//if port == "" {
		//	node.SetPort(strconv.Itoa(int(topLevelPort)))
		//} else {
		//	node.SetPort(strconv.Itoa(int(convertPortToInt(port, topLevelPort))))
		//}
		if protocol == "" {
			node.SetProtocol(ES_PROTO_ELASTIC)
		} else {
			node.SetProtocol(protocol)
		}
		if roles == "" {
			node.SetRoles(ES_ROLES_MASTER_DATA)
		} else {
			node.SetRoles(roles)
		}
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

func (c *Elastic) IsUpdateBatchAllowed(d *schema.ResourceData) error {
	var err error

	if err = c.Common.IsUpdateBatchAllowed(d); err != nil {
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
	slog.Debug(funcName)

	var err error

	// parent/super - get common attributes
	if err = c.Backup.GetBackupInputs(d, jobData); err != nil {
		return err
	}

	//jobData.SetBackupMethod("")

	snapshotRepo := d.Get(TF_FIELD_CLUSTER_SNAPSHOT_REPO).(string)
	jobData.SetSnapshotRepository(snapshotRepo)

	return err
}

func (c *Elastic) IsValidBackupOptions(vendor string, clusterType string, jobData *openapi.JobsJobJobSpecJobData) error {
	return c.Backup.IsValidBackupOptions(vendor, clusterType, jobData)
}

func (c *Elastic) SetBackupJobData(jobData *openapi.JobsJobJobSpecJobData) error {
	return c.Backup.SetBackupJobData(jobData)
}

func (c *Elastic) IsBackupRemovable(clusterInfo *openapi.ClusterResponse, jobData *openapi.JobsJobJobSpecJobData) bool {
	return true
}

func NewElastic() *Elastic {
	return &Elastic{}
}
