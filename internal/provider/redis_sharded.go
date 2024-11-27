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

type RedisSharded struct {
	//Common DbCommon
	//Backup DbBackupCommon
	RedisStuff Redis
}

func (m *RedisSharded) GetInputs(d *schema.ResourceData, jobData *openapi.JobsJobJobSpecJobData) error {
	funcName := "RedisSharded::GetInputs"
	slog.Debug(funcName)

	var err error

	// parent/super - get common attributes
	if err = m.RedisStuff.Common.GetInputs(d, jobData); err != nil {
		return err
	}

	clusterType := jobData.GetClusterType()

	var iPort int
	port := d.Get(TF_FIELD_CLUSTER_REDIS_PORT).(string)
	if err = CheckForEmptyAndSetDefault(&port, gDefultDbPortMap, clusterType); err != nil {
		return err
	}
	if iPort, err = strconv.Atoi(port); err != nil {
		slog.Error(funcName, "ERROR", "Non-numeric database port")
		return err
	}
	jobData.SetPort(int32(iPort))
	jobData.SetRedisShardedPort(int32(iPort))
	jobData.SetValkeyShardedPort(int32(iPort))

	nodeTimeoutMs := d.Get(TF_FIELD_CLUSTER_REDIS_NODE_TIMEOUT_MS).(int)
	jobData.SetNodeTimeoutMs(int32(nodeTimeoutMs))

	replicaValidityFactor := d.Get(TF_FIELD_CLUSTER_REDIS_REPLICA_VALIDITY_FACTOR).(int)
	jobData.SetRedisClusterReplicaValidityFactor(int32(replicaValidityFactor))
	jobData.SetValkeyClusterReplicaValidityFactor(int32(replicaValidityFactor))

	dbVendor := jobData.GetVendor()
	dbVersion := jobData.GetVersion()
	topLevelPort := jobData.GetPort()

	vendorMap, ok := gDbConfigTemplate[dbVendor]
	if !ok {
		return errors.New(fmt.Sprintf("Map doesn't support DB vendor: %s", dbVendor))
	}
	clusterTypeMap, ok := vendorMap[clusterType]
	if !ok {
		return errors.New(fmt.Sprintf("Map doesn't support DB vendor: %s, ClusterType: %s", dbVendor, clusterType))
	}
	cfgTemplate, ok := clusterTypeMap[dbVersion]
	if !ok {
		return errors.New(fmt.Sprintf("Map doesn't support DB vendor: %s, ClusterType: %s, DbVersion: %s", dbVendor, clusterType, dbVersion))
	}
	jobData.SetConfigTemplate(cfgTemplate)

	busPort := d.Get(TF_FIELD_CLUSTER_REDIS_BUS_PORT).(string)
	if busPort == "" {
		busPort = DEFAULT_REDIS_BUS_PORT
	}
	var iBusPort int
	if iBusPort, err = strconv.Atoi(busPort); err != nil {
		slog.Error(funcName, "ERROR", "Non-numeric database bus port")
		return err
	}
	jobData.SetRedisShardedBusPort(int32(iBusPort))
	jobData.SetValkeyShardedBusPort(int32(iBusPort))

	dataDirectory := d.Get(TF_FIELD_CLUSTER_DATA_DIR).(string)
	if err = CheckForEmptyAndSetDefault(&dataDirectory, gDefultDataDir, clusterType); err != nil {
		return err
	}
	jobData.SetDatadir(dataDirectory)

	hosts := d.Get(TF_FIELD_CLUSTER_HOST)
	nodes := []openapi.JobsJobJobSpecJobDataNodesInner{}
	for _, ff := range hosts.([]any) {
		f := ff.(map[string]any)
		hostname := f[TF_FIELD_CLUSTER_HOSTNAME].(string)
		hostname_data := f[TF_FIELD_CLUSTER_HOSTNAME_DATA].(string)
		hostname_internal := f[TF_FIELD_CLUSTER_HOSTNAME_INT].(string)
		role := f[TF_FIELD_CLUSTER_HOST_ROLE].(string)
		datadir := f[TF_FIELD_CLUSTER_HOST_DD].(string)

		if hostname == "" {
			return errors.New("Hostname cannot be empty")
		}
		var node = openapi.JobsJobJobSpecJobDataNodesInner{
			Hostname: &hostname,
		}

		node.SetClassName(CMON_CLASS_NAME_REDIS_SHARDED_HOST)

		if role == "" {
			return errors.New("Role cannot be empty")
		}
		if !strings.EqualFold(role, CMON_DB_HOST_ROLE_PRIMARY) &&
			!strings.EqualFold(role, CMON_DB_HOST_ROLE_REPLICA) {
			return errors.New("Unsupported role for memmber host.")
		}
		node.SetRole(role)

		if hostname_data != "" {
			node.SetHostnameData(hostname_data)
		}
		if hostname_internal != "" {
			node.SetHostnameInternal(hostname_internal)
		}
		node.SetPort(strconv.Itoa(int(topLevelPort)))
		//if port == "" {
		//	node.SetPort(strconv.Itoa(int(topLevelPort)))
		//} else {
		//	node.SetPort(strconv.Itoa(int(convertPortToInt(port, topLevelPort))))
		//}
		if datadir != "" {
			node.SetDatadir(datadir)
		}

		nodes = append(nodes, node)
	}
	jobData.SetNodes(nodes)

	return nil
}

func (c *RedisSharded) HandleRead(ctx context.Context, d *schema.ResourceData, apiClient *openapi.APIClient, clusterInfo *openapi.ClusterResponse) error {

	if err := c.RedisStuff.Common.HandleRead(ctx, d, apiClient, clusterInfo); err != nil {
		return err
	}

	return nil
}

func (c *RedisSharded) IsUpdateBatchAllowed(d *schema.ResourceData) error {
	return (c.RedisStuff.IsUpdateBatchAllowed(d))
}

func (c *RedisSharded) HandleUpdate(ctx context.Context, d *schema.ResourceData, apiClient *openapi.APIClient, clusterInfo *openapi.ClusterResponse) error {
	funcName := "RedisSharded::HandleUpdate"
	slog.Debug(funcName)

	var err error

	// handle things like cluster-name, tags, and toggling cluster-auto-covery in base ...
	if err := c.RedisStuff.Common.HandleUpdate(ctx, d, apiClient, clusterInfo); err != nil {
		return err
	}

	tmpJobData := openapi.NewJobsJobJobSpecJobData()
	if err = c.GetInputs(d, tmpJobData); err != nil {
		return err
	}

	strPort := strconv.Itoa(int(tmpJobData.GetPort()))

	if d.HasChange(TF_FIELD_CLUSTER_HOST) {
		var nodesToAdd []openapi.JobsJobJobSpecJobDataNodesInner
		var nodesToRemove []openapi.JobsJobJobSpecJobDataNodesInner

		hostClassName := CMON_CLASS_NAME_REDIS_SHARDED_HOST
		command := CMON_JOB_ADD_NODE_COMMAND

		// Compare Terraform and CMON to determine whether adding node, remove node or promoting standby/slave
		nodes, _ := c.RedisStuff.Common.getHosts(d)
		if nodesToAdd, nodesToRemove, err = c.RedisStuff.Common.determineNodesDelta(nodes, clusterInfo, hostClassName); err != nil {
			return err
		}

		isAddNode := len(nodesToAdd) > 0
		isRemoveNode := len(nodesToRemove) > 0

		if isAddNode && len(nodesToAdd) > 1 {
			for i := 0; i < len(nodesToAdd); i++ {
				slog.Info(funcName, "node", nodesToAdd[i].GetHostname())
			}
			return errors.New("Can't add more than one node at a time")
		}

		if isRemoveNode && len(nodesToRemove) > 1 {
			for i := 0; i < len(nodesToAdd); i++ {
				slog.Info(funcName, "node", nodesToAdd[i].GetHostname())
			}
			return errors.New("Can't remove more than one node at a time")
		}

		if isAddNode && isRemoveNode {
			return errors.New("Sorry, can't add and remove a node in a single operation.")
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
			return errors.New("Standby promotion is is not supported for Redis and Valkey.")
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

		//var primaryInCmon *openapi.ClusterResponseHostsInner
		//// Find the Primary/Master node in CMON
		//if primaryInCmon, err = c.RedisStuff.Common.findMasterNode(clusterInfo, hostClassName, CMON_DB_HOST_ROLE_MASTER); err != nil {
		//	return err
		//}
		//slog.Debug(funcName, "Master:", primaryInCmon.GetHostname())

		jobData.SetInstallSoftware(tmpJobData.GetInstallSoftware())
		jobData.SetDisableSelinux(tmpJobData.GetDisableSelinux())
		jobData.SetEnableUninstall(true /*tmpJobData.GetEnableUninstall()*/)
		jobData.SetDisableFirewall(tmpJobData.GetDisableFirewall())

		if isAddNode {
			var node openapi.JobsJobJobSpecJobDataNodesInner
			var nodes []openapi.JobsJobJobSpecJobDataNodesInner
			node.SetClassName(hostClassName)
			node.SetHostname(nodeFromTf.GetHostname())

			hostTfRec := c.RedisStuff.Common.findHostEntry(nodeFromTf.GetHostname(), d.Get(TF_FIELD_CLUSTER_HOST))
			hostname_data := hostTfRec[TF_FIELD_CLUSTER_HOSTNAME_DATA].(string)
			hostname_internal := hostTfRec[TF_FIELD_CLUSTER_HOSTNAME_INT].(string)
			role := hostTfRec[TF_FIELD_CLUSTER_HOST_ROLE].(string)
			//port := hostTfRec[TF_FIELD_CLUSTER_HOST_PORT].(string)
			datadir := hostTfRec[TF_FIELD_CLUSTER_HOST_DD].(string)

			if hostname_data != "" {
				node.SetHostnameData(hostname_data)
			} else {
				node.SetHostnameData(node.GetHostname())
			}

			if hostname_internal != "" {
				node.SetHostnameInternal(hostname_internal)
			}

			if role == "" {
				return errors.New("Role cannot be empty")
			}
			if !strings.EqualFold(role, CMON_DB_HOST_ROLE_PRIMARY) &&
				!strings.EqualFold(role, CMON_DB_HOST_ROLE_REPLICA) {
				return errors.New("Unsupported role for memmber host.")
			}
			node.SetRole(role)

			if strings.EqualFold(role, CMON_DB_HOST_ROLE_REPLICA) {
				primaryHostTfRec := c.RedisStuff.Common.findPrimaryOfReplica(nodeFromTf.GetHostname(), d.Get(TF_FIELD_CLUSTER_HOST))
				if primaryHostTfRec == nil {
					return errors.New("Cannot add Replica without a corresponding primary host.")
				}
				primaryHost := primaryHostTfRec[TF_FIELD_CLUSTER_HOSTNAME].(string)

				jobData.SetMasterAddress(primaryHost + ":" + strPort)
			}

			if strings.EqualFold(tmpJobData.GetClusterType(), CLUSTER_TYPE_REDIS_SHARDED) {
				node.SetProtocol(CLUSTER_TYPE_REDIS_SHARDED)
			} else if strings.EqualFold(tmpJobData.GetClusterType(), CLUSTER_TYPE_VALKEY_SHARDED) {
				node.SetProtocol(CLUSTER_TYPE_VALKEY_SHARDED)
			} else {
				return errors.New("Unsupported cluster-type:" + tmpJobData.GetClusterType())
			}

			if datadir != "" {
				node.SetDatadir(datadir)
			}

			node.SetPort(strPort)
			//if port != "" {
			//	node.SetPort(strconv.Itoa(int(convertPortToInt(port, tmpJobData.GetPort()))))
			//} else {
			//	node.SetPort(strconv.Itoa(int(tmpJobData.GetPort())))
			//}

			nodes = append(nodes, node)
			jobData.SetNodes(nodes)
		} else if isRemoveNode {
			var node openapi.JobsJobJobSpecJobDataNodesInner
			var nodes []openapi.JobsJobJobSpecJobDataNodesInner

			node.SetHostname(nodeToAddOrRemove.GetHostname())
			node.SetPort(strPort)

			jobData.SetEnableUninstall(true)
			jobData.SetUnregisterOnly(false)

			nodes = append(nodes, node)

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

func (c *RedisSharded) GetBackupInputs(d *schema.ResourceData, jobData *openapi.JobsJobJobSpecJobData) error {
	funcName := "RedisSharded::GetBackupInputs"
	slog.Debug(funcName)

	var err error

	// parent/super - get common attributes
	if err = c.RedisStuff.GetBackupInputs(d, jobData); err != nil {
		return err
	}

	return err
}

func (c *RedisSharded) IsValidBackupOptions(vendor string, clusterType string, jobData *openapi.JobsJobJobSpecJobData) error {
	return c.RedisStuff.IsValidBackupOptions(vendor, clusterType, jobData)
}

func (c *RedisSharded) SetBackupJobData(jobData *openapi.JobsJobJobSpecJobData) error {
	return c.RedisStuff.SetBackupJobData(jobData)
}

func (c *RedisSharded) IsBackupRemovable(clusterInfo *openapi.ClusterResponse, jobData *openapi.JobsJobJobSpecJobData) bool {
	return c.RedisStuff.IsBackupRemovable(clusterInfo, jobData)
}

func NewRedisSharded() *RedisSharded {
	return &RedisSharded{}
}
