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

type Redis struct {
	Common DbCommon
	Backup DbBackupCommon
}

func (m *Redis) GetInputs(d *schema.ResourceData, jobData *openapi.JobsJobJobSpecJobData) error {
	funcName := "Redis::GetInputs"
	slog.Debug(funcName)

	var err error

	// parent/super - get common attributes
	if err = m.Common.GetInputs(d, jobData); err != nil {
		return err
	}

	clusterType := jobData.GetClusterType()
	topLevelPort := jobData.GetPort()

	dataDirectory := d.Get(TF_FIELD_CLUSTER_DATA_DIR).(string)
	if err = CheckForEmptyAndSetDefault(&dataDirectory, gDefultDataDir, clusterType); err != nil {
		return err
	}
	jobData.SetDatadir(dataDirectory)

	sentinelPort := d.Get(TF_FIELD_CLUSTER_SENTINEL_PORT).(string)
	if sentinelPort == "" {
		sentinelPort = DEFAULT_MONGO_REDIS_SENTINEL_PORT
	}
	jobData.SetSentinelPort(sentinelPort)

	hosts := d.Get(TF_FIELD_CLUSTER_HOST)
	nodes := []openapi.JobsJobJobSpecJobDataNodesInner{}
	for _, ff := range hosts.([]any) {
		f := ff.(map[string]any)
		hostname := f[TF_FIELD_CLUSTER_HOSTNAME].(string)
		hostname_data := f[TF_FIELD_CLUSTER_HOSTNAME_DATA].(string)
		hostname_internal := f[TF_FIELD_CLUSTER_HOSTNAME_INT].(string)
		port := f[TF_FIELD_CLUSTER_HOST_PORT].(string)
		datadir := f[TF_FIELD_CLUSTER_HOST_DD].(string)

		if hostname == "" {
			return errors.New("Hostname cannot be empty")
		}
		var node = openapi.JobsJobJobSpecJobDataNodesInner{
			Hostname: &hostname,
		}

		node.SetClassName(CMON_CLASS_NAME_REDIS_HOST)

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

		var node2 = node
		node2.SetClassName(CMON_CLASS_NAME_REDIS_SENTNEL_HOST)
		node2.SetPort(sentinelPort)

		nodes = append(nodes, node, node2)
	}
	jobData.SetNodes(nodes)

	return nil
}

func (c *Redis) HandleRead(ctx context.Context, d *schema.ResourceData, m interface{}, clusterInfo *openapi.ClusterResponse) error {

	if err := c.Common.HandleRead(ctx, d, m, clusterInfo); err != nil {
		return err
	}

	return nil
}

func (c *Redis) IsUpdateBatchAllowed(d *schema.ResourceData) error {
	var err error

	if err = c.Common.IsUpdateBatchAllowed(d); err != nil {
		return err
	}

	updateClassA := d.HasChange(TF_FIELD_CLUSTER_HOST)
	updateClassAprime := d.HasChangeExcept(TF_FIELD_CLUSTER_HOST)
	if updateClassA && updateClassAprime {
		return errors.New(fmt.Sprintf("You are not allowed to update %s along with any other fields.", TF_FIELD_CLUSTER_HOST))
	}

	return nil
}

func (c *Redis) HandleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}, clusterInfo *openapi.ClusterResponse) error {
	funcName := "Redis::HandleUpdate"
	slog.Info(funcName)

	var err error

	// handle things like cluster-name, tags, and toggling cluster-auto-covery in base ...
	if err := c.Common.HandleUpdate(ctx, d, m, clusterInfo); err != nil {
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

		hostClassName := CMON_CLASS_NAME_REDIS_HOST
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
			//command = CMON_JOB_PROMOTE_REPLICAION_SLAVE_COMMAND
			// Here we are dealing with a Role change (slave promotion to master)
			return errors.New("Standby promotion is is not supported for Redis")
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
		slog.Debug(funcName, "Master:", primaryInCmon.GetHostname())

		jobData.SetInstallSoftware(tmpJobData.GetInstallSoftware())
		jobData.SetDisableSelinux(tmpJobData.GetDisableSelinux())
		jobData.SetEnableUninstall(true /*tmpJobData.GetEnableUninstall()*/)
		jobData.SetDisableFirewall(tmpJobData.GetDisableFirewall())

		if isAddNode {
			var node openapi.JobsJobJobSpecJobDataNodesInner
			var node2 openapi.JobsJobJobSpecJobDataNodesInner
			var nodes []openapi.JobsJobJobSpecJobDataNodesInner
			node.SetClassName(hostClassName)
			node.SetHostname(nodeFromTf.GetHostname())
			node2.SetClassName(CMON_CLASS_NAME_REDIS_SENTNEL_HOST)
			node2.SetHostname(nodeFromTf.GetHostname())

			hostTfRec := c.Common.findHostEntry(nodeFromTf.GetHostname(), d.Get(TF_FIELD_CLUSTER_HOST))
			hostname_data := hostTfRec[TF_FIELD_CLUSTER_HOSTNAME_DATA].(string)
			hostname_internal := hostTfRec[TF_FIELD_CLUSTER_HOSTNAME_INT].(string)
			port := hostTfRec[TF_FIELD_CLUSTER_HOST_PORT].(string)
			datadir := hostTfRec[TF_FIELD_CLUSTER_HOST_DD].(string)

			if hostname_data != "" {
				node.SetHostnameData(hostname_data)
			} else {
				node.SetHostnameData(node.GetHostname())
			}
			node2.SetHostnameData(node.GetHostnameData())

			if hostname_internal != "" {
				node.SetHostnameInternal(hostname_internal)
				node2.SetHostnameInternal(node.GetHostnameInternal())
			}
			if port != "" {
				node.SetPort(strconv.Itoa(int(convertPortToInt(port, tmpJobData.GetPort()))))
			} else {
				node.SetPort(strconv.Itoa(int(tmpJobData.GetPort())))
			}
			if datadir != "" {
				node.SetDatadir(datadir)
			}

			node2.SetPort(tmpJobData.GetSentinelPort())

			nodes = append(nodes, node, node2)
			jobData.SetNodes(nodes)
		} else if isRemoveNode {
			var node openapi.JobsJobJobSpecJobDataNodesInner
			var node2 openapi.JobsJobJobSpecJobDataNodesInner
			var nodes []openapi.JobsJobJobSpecJobDataNodesInner

			node.SetHostname(nodeToAddOrRemove.GetHostname())
			node.SetPort(strconv.Itoa(int(tmpJobData.GetPort())))

			node2.SetHostname(nodeToAddOrRemove.GetHostname())
			node2.SetPort(tmpJobData.GetSentinelPort())

			jobData.SetEnableUninstall(true)
			jobData.SetUnregisterOnly(false)

			nodes = append(nodes, node, node2)

			jobData.SetNodes(nodes)
		}

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

func (c *Redis) GetBackupInputs(d *schema.ResourceData, jobData *openapi.JobsJobJobSpecJobData) error {
	funcName := "Redis::GetBackupInputs"
	slog.Info(funcName)

	var err error

	// parent/super - get common attributes
	if err = c.Backup.GetBackupInputs(d, jobData); err != nil {
		return err
	}

	jobData.SetHostname(STINRG_AUTO)

	return err
}

func (c *Redis) IsValidBackupOptions(vendor string, clusterType string, jobData *openapi.JobsJobJobSpecJobData) error {
	return c.Backup.IsValidBackupOptions(vendor, clusterType, jobData)
}

func (c *Redis) SetBackupJobData(jobData *openapi.JobsJobJobSpecJobData) error {
	return c.Backup.SetBackupJobData(jobData)
}

func NewRedis() *Redis {
	return &Redis{}
}
