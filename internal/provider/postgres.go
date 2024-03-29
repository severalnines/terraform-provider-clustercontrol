package provider

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/severalnines/clustercontrol-client-sdk/go/pkg/openapi"
	"log/slog"
	"strconv"
	"strings"
)

type PostgresSql struct {
	Common DbCommon
	Backup DbBackupCommon
}

func (m *PostgresSql) GetInputs(d *schema.ResourceData, jobData *openapi.JobsJobJobSpecJobData) error {
	funcName := "PostgresSql::GetInputs"
	slog.Debug(funcName)

	var err error

	// parent/super - get common attributes
	if err = m.Common.GetInputs(d, jobData); err != nil {
		return err
	}

	topLevelPort := jobData.GetPort()

	dataDirectory := d.Get(TF_FIELD_CLUSTER_DATA_DIR).(string)
	if dataDirectory != "" {
		jobData.SetDatadir(dataDirectory)
	}

	timescaleExt := d.Get(TF_FIELD_CLUSTER_PG_TIMESALE_EXT).(bool)
	jobData.SetInstallTimescaledb(timescaleExt)

	hosts := d.Get(TF_FIELD_CLUSTER_HOST)
	nodes := []openapi.JobsJobJobSpecJobDataNodesInner{}
	for _, ff := range hosts.([]any) {
		f := ff.(map[string]any)

		hostname := f[TF_FIELD_CLUSTER_HOSTNAME].(string)
		hostname_data := f[TF_FIELD_CLUSTER_HOSTNAME_DATA].(string)
		hostname_internal := f[TF_FIELD_CLUSTER_HOSTNAME_INT].(string)
		port := f[TF_FIELD_CLUSTER_HOST_PORT].(string)
		datadir := f[TF_FIELD_CLUSTER_HOST_DD].(string)
		sync_replication := f[TF_FIELD_CLUSTER_SYNC_REP].(bool)

		if hostname == "" {
			return errors.New("Hostname cannot be empty")
		}
		var node = openapi.JobsJobJobSpecJobDataNodesInner{
			Hostname: &hostname,
		}
		if hostname_data != "" {
			node.SetHostnameData(hostname_data)
		}
		if hostname_internal != "" {
			node.SetHostnameInternal(hostname_internal)
		}
		if port == "" {
			node.SetPort(strconv.Itoa(int(topLevelPort)))
		} else {
			node.SetPort(strconv.Itoa(int(convertPortToInt(port, topLevelPort))))
		}
		if datadir != "" {
			node.SetDatadir(datadir)
		}
		node.SetSynchronous(sync_replication)

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
	var err error

	if err = c.Common.IsUpdateBatchAllowed(d); err != nil {
		return err
	}

	updateClassA := d.HasChange(TF_FIELD_CLUSTER_HOST)
	updateClassAprime := d.HasChangeExcept(TF_FIELD_CLUSTER_HOST)
	if updateClassA && updateClassAprime {
		return errors.New(fmt.Sprintf("You are not allowed to update %s along with any other fields.", TF_FIELD_CLUSTER_HOST))
	}

	updateClassA = d.HasChange(TF_FIELD_CLUSTER_ENABLE_PGBACKREST_AGENT)
	updateClassAprime = d.HasChangeExcept(TF_FIELD_CLUSTER_ENABLE_PGBACKREST_AGENT)
	if updateClassA && updateClassAprime {
		err = errors.New(fmt.Sprintf("You are not allowed to update %s along with any other fields.", TF_FIELD_CLUSTER_ENABLE_PGM_AGENT))
		return err
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

			// NOTE: host is guaranteed to be non-nil.
			hostTfRec := c.Common.findHostEntry(nodeFromTf.GetHostname(), d.Get(TF_FIELD_CLUSTER_HOST))
			hostname_data := hostTfRec[TF_FIELD_CLUSTER_HOSTNAME_DATA].(string)
			hostname_internal := hostTfRec[TF_FIELD_CLUSTER_HOSTNAME_INT].(string)
			port := hostTfRec[TF_FIELD_CLUSTER_HOST_PORT].(string)
			datadir := hostTfRec[TF_FIELD_CLUSTER_HOST_DD].(string)
			syncReplication := hostTfRec[TF_FIELD_CLUSTER_SYNC_REP].(bool)

			if hostname_data != "" {
				node.SetHostnameData(hostname_data)
			} else {
				node.SetHostnameData(node.GetHostname())
			}
			if hostname_internal != "" {
				node.SetHostnameInternal(hostname_internal)
			}
			if port != "" {
				node.SetPort(convertPortToInt(port, tmpJobData.GetPort()))
			} else {
				node.SetPort(tmpJobData.GetPort())
			}
			if datadir != "" {
				node.SetDatadir(datadir)
			}

			node.SetSynchronous(syncReplication)

		} else if isRemoveNode {

			node.SetHostname(nodeToAddOrRemove.GetHostname())
			node.SetPort(tmpJobData.GetPort())
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
	} // d.HasChange(TF_FIELD_CLUSTER_HOST)

	if d.HasChange(TF_FIELD_CLUSTER_ENABLE_PGBACKREST_AGENT) {
		enablePgB := d.Get(TF_FIELD_CLUSTER_ENABLE_PGBACKREST_AGENT).(bool)
		if !enablePgB {
			// Don't support disabling pgbackrest at this time
			return errors.New("Disabling PgBackRest is not currently supported")
		}
		apiClient := m.(*openapi.APIClient)

		enablePbmJob := NewCCJob(CMON_JOB_CREATE_JOB)
		job := enablePbmJob.GetJob()
		jobSpec := job.GetJobSpec()
		jobData := jobSpec.GetJobData()
		enablePbmJob.SetClusterId(clusterInfo.GetClusterId())
		jobSpec.SetCommand(CMON_JOB_PGBACKREST_COMMAND)
		jobData.SetAction(JOB_ACTION_SETUP)

		var nodes = []openapi.JobsJobJobSpecJobDataNodesInner{}
		var node = openapi.JobsJobJobSpecJobDataNodesInner{}
		node.SetClassName(CMON_CLASS_NAME_PGBACKREST_HOST)
		node.SetHostname("*")
		nodes = append(nodes, node)
		jobData.SetNodes(nodes)

		jobSpec.SetJobData(jobData)
		job.SetJobSpec(jobSpec)
		enablePbmJob.SetJob(job)

		if err = SendAndWaitForJobCompletion(ctx, apiClient, enablePbmJob); err != nil {
			slog.Error(err.Error())
		}
	} // d.HasChange(TF_FIELD_CLUSTER_ENABLE_PGBACKREST_AGENT)

	return nil
}

func (c *PostgresSql) GetBackupInputs(d *schema.ResourceData, jobData *openapi.JobsJobJobSpecJobData) error {
	funcName := "PostgresSql::GetBackupInputs"
	slog.Info(funcName)

	var err error

	// parent/super - get common attributes
	if err = c.Backup.GetBackupInputs(d, jobData); err != nil {
		return err
	}

	return err
}

func (c *PostgresSql) IsValidBackupOptions(vendor string, clusterType string, jobData *openapi.JobsJobJobSpecJobData) error {
	return c.Backup.IsValidBackupOptions(vendor, clusterType, jobData)
}

func (c *PostgresSql) SetBackupJobData(jobData *openapi.JobsJobJobSpecJobData) error {
	var err error
	if err = c.Backup.SetBackupJobData(jobData); err != nil {
		return err
	}

	// Need to set port if backup-method is pgbackrest
	backupMethod := jobData.GetBackupMethod()
	if strings.Contains(backupMethod, "pgbackrest") {
		// unfortunately had to had code it here due to inability to unmarshall getcluster
		// result into data structures. Why? Because the `port` field data type is inconsistent
		// in the returned response from CMON. Sometime int, other times string !!!
		jobData.SetPort(5432)
	}
	return nil
}

func (c *PostgresSql) IsBackupRemovable(clusterInfo *openapi.ClusterResponse, jobData *openapi.JobsJobJobSpecJobData) bool {
	backupMethod := jobData.GetBackupMethod()
	if strings.Contains(backupMethod, "pgbackrest") {
		return false
	}
	return true
}

func NewPostgres() *PostgresSql {
	return &PostgresSql{}
}
