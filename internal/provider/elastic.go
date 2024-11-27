package provider

import (
	"context"
	"errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/severalnines/clustercontrol-client-sdk/go/pkg/openapi"
	"log/slog"
	"strconv"
	"strings"
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

func (c *Elastic) HandleRead(ctx context.Context, d *schema.ResourceData, apiClient *openapi.APIClient, clusterInfo *openapi.ClusterResponse) error {

	if err := c.Common.HandleRead(ctx, d, apiClient, clusterInfo); err != nil {
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

func (c *Elastic) HandleUpdate(ctx context.Context, d *schema.ResourceData, apiClient *openapi.APIClient, clusterInfo *openapi.ClusterResponse) error {
	funcName := "Elastic::HandleUpdate"
	slog.Debug(funcName)

	var err error

	if err := c.Common.HandleUpdate(ctx, d, apiClient, clusterInfo); err != nil {
		return err
	}

	tmpJobData := openapi.NewJobsJobJobSpecJobData()
	if err = c.GetInputs(d, tmpJobData); err != nil {
		return err
	}

	if d.HasChange(TF_FIELD_CLUSTER_HOST) {
		var nodesToAdd []openapi.JobsJobJobSpecJobDataNodesInner
		var nodesToRemove []openapi.JobsJobJobSpecJobDataNodesInner

		hostClassName := CMON_CLASS_NAME_ELASTIC_HOST
		command := CMON_JOB_ADD_NODE_COMMAND

		// Compare Terraform and CMON to determine whether adding node, remove node or promoting standby/slave
		nodes, _ := c.Common.getHosts(d)
		if nodesToAdd, nodesToRemove, err = c.Common.determineNodesDelta(nodes, clusterInfo, hostClassName); err != nil {
			return err
		}

		isAddNode := len(nodesToAdd) > 0
		isRemoveNode := len(nodesToRemove) > 0

		if isAddNode && len(nodesToAdd) > 1 {
			return errors.New("Can't add more than one node at a time")
		}

		if isRemoveNode && len(nodesToRemove) > 1 {
			return errors.New("Can't remove more than one node at a time")
		}

		var nodeToAddOrRemove *openapi.JobsJobJobSpecJobDataNodesInner
		if isAddNode {
			nodeToAddOrRemove = &nodesToAdd[0]
		} else if isRemoveNode {
			nodeToAddOrRemove = &nodesToRemove[0]
			command = CMON_JOB_REMOVE_NODE_COMMAND
		} else {
			return errors.New("Unsupported Elasticsearch operation. Neither Add nor Remove node.")
		}

		// From Terraform
		tmpJobDataNodes := tmpJobData.GetNodes()
		var nodeFromTf *openapi.JobsJobJobSpecJobDataNodesInner
		for i := 1; i < len(tmpJobDataNodes) && nodeToAddOrRemove != nil; i++ {
			tmpJobDataNode := tmpJobDataNodes[i]
			if strings.EqualFold(tmpJobDataNode.GetHostname(), nodeToAddOrRemove.GetHostname()) {
				nodeFromTf = &tmpJobDataNode
				break
			}
		}
		// No need to error check as the node must be in the list

		// variables at this point
		// tmpJobData: clustercontrol_db_cluster TF resource data
		// tmpJobDataNodes: the db_hosts portion of TF resource data
		// nodeFromTf: "the" node from the resource data. It is this node which is to be added or removed to the cluster
		// nodeToAddOrRemove: contains hostname of the node to be added or removed. Use it in the remove case

		//apiClient := m.(*openapi.APIClient)
		addOrRemoveNodeJob := NewCCJob(CMON_JOB_CREATE_JOB)
		addOrRemoveNodeJob.SetClusterId(clusterInfo.GetClusterId())
		job := addOrRemoveNodeJob.GetJob()
		jobSpec := job.GetJobSpec()
		jobData := jobSpec.GetJobData()
		jobSpec.SetCommand(command)

		// There's no concept of Primary node in Elastic, so we can skip that step here

		// Set job data fields
		jobData.SetInstallSoftware(tmpJobData.GetInstallSoftware())
		jobData.SetDisableSelinux(tmpJobData.GetDisableSelinux())
		jobData.SetDisableFirewall(tmpJobData.GetDisableFirewall())

		if isAddNode {

			var nodes []openapi.JobsJobJobSpecJobDataNodesInner
			var node openapi.JobsJobJobSpecJobDataNodesInner

			node.SetClassName(hostClassName)
			node.SetHostname(nodeFromTf.GetHostname())

			// NOTE: host is guaranteed to be non-nil.
			hostTfRec := c.Common.findHostEntry(nodeFromTf.GetHostname(), d.Get(TF_FIELD_CLUSTER_HOST))
			hostname_data := hostTfRec[TF_FIELD_CLUSTER_HOSTNAME_DATA].(string)
			hostname_internal := hostTfRec[TF_FIELD_CLUSTER_HOSTNAME_INT].(string)
			protocol := hostTfRec[TF_FIELD_CLUSTER_HOST_PROTO].(string)
			roles := hostTfRec[TF_FIELD_CLUSTER_HOST_ROLES].(string)

			if hostname_data != "" {
				node.SetHostnameData(hostname_data)
			} else {
				node.SetHostnameData(node.GetHostname())
			}
			if hostname_internal != "" {
				node.SetHostnameInternal(hostname_internal)
			}
			node.SetPort(strconv.Itoa(int(tmpJobData.GetPort())))

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
			jobData.SetNodes(nodes)
		} else if isRemoveNode {
			var node openapi.JobsJobJobSpecJobDataNode
			node.SetHostname(nodeToAddOrRemove.GetHostname())
			node.SetPort(tmpJobData.GetPort())
			jobData.SetEnableUninstall(true)
			jobData.SetUnregisterOnly(false)
			jobData.SetNode(node)
		} else {
			// Here we are dealing with a Role change (slave promotion to master)
			return errors.New("Standby promotion is yet to be supported")
		}

		jobSpec.SetJobData(jobData)
		job.SetJobSpec(jobSpec)
		addOrRemoveNodeJob.SetJob(job)

		if err = SendAndWaitForJobCompletion(ctx, apiClient, addOrRemoveNodeJob); err != nil {
			slog.Error(err.Error())
			return err
		}

	} // d.HasChange(TF_FIELD_CLUSTER_HOST)

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
