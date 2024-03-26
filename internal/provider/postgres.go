package provider

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/severalnines/clustercontrol-client-sdk/go/pkg/openapi"
	"log/slog"
	"strings"
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
		m.Common.getCommonHostAttributes(f, iPort, clusterType, memHost)
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

func (c *PostgresSql) IsUpdateBatchAllowed(d *schema.ResourceData) error {

	if err := c.Common.IsUpdateBatchAllowed(d); err != nil {
		return err
	}

	updateClassA := d.HasChange(TF_FIELD_CLUSTER_HOST)
	updateClassAprime := d.HasChangeExcept(TF_FIELD_CLUSTER_HOST)
	if updateClassA && updateClassAprime {
		return errors.New(fmt.Sprintf("You are not allowed to update %s along with any other fields.", TF_FIELD_CLUSTER_HOST))
	}

	return nil
}

func (c *PostgresSql) HandleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}, clusterInfo *openapi.ClusterResponse) error {
	funcName := "PostgresSql::HandleUpdate"
	slog.Info(funcName)

	var err error

	// handle things like cluster-name, tags, and toggling cluster-auto-covery in base ...
	if err = c.Common.HandleUpdate(ctx, d, m, clusterInfo); err != nil {
		return err
	}

	// handle other big updates such as add-replication-slave, remove-node, etc
	tmpJobData := openapi.NewJobsJobJobSpecJobData()
	if err = c.GetInputs(d, tmpJobData); err != nil {
		return err
	}

	if d.HasChange(TF_FIELD_CLUSTER_HOST) {
		var nodesToAdd []openapi.JobsJobJobSpecJobDataNodesInner
		var nodesToRemove []openapi.JobsJobJobSpecJobDataNodesInner

		hostClassName := CMON_CLASS_NAME_PROSGRESQL_HOST
		command := CMON_JOB_ADD_REPLICATION_SLAVE_COMMAND

		// Compare Terraform and CMON to determine whether adding node, remove node or promoting standby/slave
		if nodesToAdd, nodesToRemove, err = c.Common.determineNodesDelta(d, clusterInfo, hostClassName); err != nil {
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
			command = CMON_JOB_PROMOTE_REPLICAION_SLAVE_COMMAND
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

		apiClient := m.(*openapi.APIClient)
		addOrRemoveNodeJob := NewCCJob(CMON_JOB_CREATE_JOB)
		addOrRemoveNodeJob.SetClusterId(clusterInfo.GetClusterId())
		job := addOrRemoveNodeJob.GetJob()
		jobSpec := job.GetJobSpec()
		jobData := jobSpec.GetJobData()
		jobSpec.SetCommand(command)

		var primaryInCmon *openapi.ClusterResponseHostsInner
		// Find the Primary/Master node in CMON
		if primaryInCmon, err = c.Common.findMasterNode(clusterInfo, hostClassName, CMON_DB_HOST_ROLE_MASTER); err != nil {
			return err
		}

		// Set job data fields
		jobData.SetConfigTemplate(tmpJobData.GetConfigTemplate())
		jobData.SetMasterAddress(fmt.Sprintf("%s:%v", primaryInCmon.GetHostname(), tmpJobData.GetPort()))
		jobData.SetInstallSoftware(tmpJobData.GetInstallSoftware())
		jobData.SetDisableSelinux(tmpJobData.GetDisableSelinux())
		jobData.SetDisableFirewall(tmpJobData.GetDisableFirewall())
		jobData.SetUpdateLb(true)
		jobData.SetDatadir(tmpJobData.GetDatadir())
		jobData.SetUsePackageForDataDir(true) // for PG
		//jobData.SetVersion("")

		var node openapi.JobsJobJobSpecJobDataNode

		if isAddNode {
			node.SetHostname(nodeFromTf.GetHostname())
			node.SetHostnameData(nodeFromTf.GetHostnameData())
			node.SetHostnameInternal(nodeFromTf.GetHostnameInternal())
			node.SetPort(convertPortToInt(nodeFromTf.GetPort(), tmpJobData.GetPort()))
			node.SetDatadir(nodeFromTf.GetDatadir())
			node.SetSynchronous(nodeFromTf.GetSynchronous())
		} else if isRemoveNode {
			node.SetHostname(nodeToAddOrRemove.GetHostname())
			jobData.SetEnableUninstall(true)
			jobData.SetUnregisterOnly(false)
		} else {
			// Here we are dealing with a Role change (slave promotion to master)
			return errors.New("Standby promotion is yet to be supported")
		}

		jobData.SetNode(node)
		jobSpec.SetJobData(jobData)
		job.SetJobSpec(jobSpec)
		addOrRemoveNodeJob.SetJob(job)

		if err = SendAndWaitForJobCompletion(ctx, apiClient, addOrRemoveNodeJob); err != nil {
			slog.Error(err.Error())
			return err
		}

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
