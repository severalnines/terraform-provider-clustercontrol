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

type MySQLMaria struct {
	Common DbCommon
	Backup DbBackupCommon
}

func (m *MySQLMaria) GetInputs(d *schema.ResourceData, jobData *openapi.JobsJobJobSpecJobData) error {
	funcName := "MySQLMaria::GetInputs"
	slog.Info(funcName)

	var err error

	// parent/super - get common attributes
	if err = m.Common.GetInputs(d, jobData); err != nil {
		return err
	}

	clusterType := jobData.GetClusterType()
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

	dataDirectory := d.Get(TF_FIELD_CLUSTER_DATA_DIR).(string)
	if err = CheckForEmptyAndSetDefault(&dataDirectory, gDefultDataDir, clusterType); err != nil {
		return err
	}
	jobData.SetDatadir(dataDirectory)

	semiSyncReplication := d.Get(TF_FIELD_CLUSTER_SEMISYNC_REP).(bool)
	jobData.SetMysqlSemiSync(semiSyncReplication)

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

func (c *MySQLMaria) IsUpdateBatchAllowed(d *schema.ResourceData) error {
	var err error

	if err = c.Common.IsUpdateBatchAllowed(d); err != nil {
		return err
	}

	updateClassA := d.HasChange(TF_FIELD_CLUSTER_HOST)
	updateClassAprime := d.HasChangeExcept(TF_FIELD_CLUSTER_HOST)
	if updateClassA && updateClassAprime {
		err = errors.New(fmt.Sprintf("You are not allowed to update %s along with any other fields.", TF_FIELD_CLUSTER_HOST))
		return err
	}

	updateClassA = d.HasChange(TF_FIELD_CLUSTER_LOAD_BALANCER)
	updateClassAprime = d.HasChangeExcept(TF_FIELD_CLUSTER_LOAD_BALANCER)
	if updateClassA && updateClassAprime {
		err = errors.New(fmt.Sprintf("You are not allowed to update %s along with any other fields.", TF_FIELD_CLUSTER_HOST))
		return err
	}

	return nil
}

func (c *MySQLMaria) HandleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}, clusterInfo *openapi.ClusterResponse) error {
	funcName := "MySQLMaria::HandleUpdate"
	slog.Info(funcName)

	var err error

	// handle things like cluster-name, tags, and toggling cluster-auto-covery in base ...
	if err = c.Common.HandleUpdate(ctx, d, m, clusterInfo); err != nil {
		return err
	}

	clusterType := clusterInfo.GetClusterType()
	clusterType = strings.ToLower(clusterType)

	if d.HasChange(TF_FIELD_CLUSTER_HOST) {
		// handle other big updates such as add-replication-slave, remove-node, etc
		tmpJobData := openapi.NewJobsJobJobSpecJobData()
		if err = c.GetInputs(d, tmpJobData); err != nil {
			return err
		}

		isReplicationType := true
		hostClassName := CMON_CLASS_NAME_MYSQL_HOST
		command := CMON_JOB_ADD_REPLICATION_SLAVE_COMMAND
		if clusterType == CLUSTER_TYPE_GALERA {
			isReplicationType = false
			hostClassName = CMON_CLASS_NAME_GALERA_HOST
			command = CMON_JOB_ADD_NODE_COMMAND
		}

		var nodesToAdd []openapi.JobsJobJobSpecJobDataNodesInner
		var nodesToRemove []openapi.JobsJobJobSpecJobDataNodesInner

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
		//jobData.SetVersion("")

		var node openapi.JobsJobJobSpecJobDataNode

		if isAddNode {

			if isReplicationType {
				jobData.SetMysqlSemiSync(true)
			} else {
				jobData.SetGaleraSegment("0")
			}

			node.SetHostname(nodeFromTf.GetHostname())

			// NOTE: host is guaranteed to be non-nil.
			h := c.Common.findHost(nodeFromTf.GetHostname(), d.Get(TF_FIELD_CLUSTER_HOST))
			hostname_data := h[TF_FIELD_CLUSTER_HOSTNAME_DATA].(string)
			hostname_internal := h[TF_FIELD_CLUSTER_HOSTNAME_INT].(string)
			port := h[TF_FIELD_CLUSTER_HOST_PORT].(string)
			datadir := h[TF_FIELD_CLUSTER_HOST_DD].(string)

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

	if d.HasChange(TF_FIELD_CLUSTER_LOAD_BALANCER) {
		apiClient := m.(*openapi.APIClient)
		addOrRemoveNodeJob := NewCCJob(CMON_JOB_CREATE_JOB)
		addOrRemoveNodeJob.SetClusterId(clusterInfo.GetClusterId())
		job := addOrRemoveNodeJob.GetJob()
		jobSpec := job.GetJobSpec()
		jobData := jobSpec.GetJobData()

		var nodesToAdd []openapi.JobsJobJobSpecJobDataNodesInner
		var nodesToRemove []openapi.JobsJobJobSpecJobDataNodesInner

		// Compare Terraform and CMON to determine whether adding node, remove node or promoting standby/slave
		if nodesToAdd, nodesToRemove, err = c.Common.determineProxyDelta(d, clusterInfo, CMON_CLASS_NAME_PROXYSQL_HOST); err != nil {
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
			jobSpec.SetCommand(CMON_JOB_CREATE_PROXYSQL_COMMAND)

			loadBalancers := d.Get(TF_FIELD_CLUSTER_LOAD_BALANCER)
			var theTfRecord = map[string]any{}
			isFound := false
			for _, ff := range loadBalancers.([]any) {
				if isFound {
					break
				}
				f := ff.(map[string]any)
				myhost := f[TF_FIELD_LB_MY_HOST]
				for _, tt := range myhost.([]any) {
					if isFound {
						break
					}
					t := tt.(map[string]any)
					hostname := t[TF_FIELD_CLUSTER_HOSTNAME].(string)
					if strings.EqualFold(hostname, nodeToAddOrRemove.GetHostname()) {
						// found it
						theTfRecord = f
						isFound = true
					}
				}
			}
			var getInputs DbLoadBalancerInterface
			getInputs = NewProxySql()
			err = getInputs.GetInputs(theTfRecord, &jobData)

		} else if isRemoveNode {
			jobSpec.SetCommand(CMON_JOB_REMOVE_NODE_COMMAND)
			nodeToAddOrRemove = &nodesToRemove[0]
			var node openapi.JobsJobJobSpecJobDataNode
			node.SetHostname(nodeToAddOrRemove.GetHostname())
			jobData.SetNode(node)
			//node.SetPort("6033")
			jobData.SetEnableUninstall(true)
			jobData.SetUnregisterOnly(false)
			slog.Info(funcName, "Removing hostname", node.GetHostname())
		} else {
			return nil
		}

		jobSpec.SetJobData(jobData)
		job.SetJobSpec(jobSpec)
		addOrRemoveNodeJob.SetJob(job)

		if err = SendAndWaitForJobCompletion(ctx, apiClient, addOrRemoveNodeJob); err != nil {
			slog.Error(err.Error())
		}

	} // d.HasChange(TF_FIELD_CLUSTER_LOAD_BALANCER)

	// ************************
	// NOTE: don't remove this block of code
	// ************************

	//if d.HasChange(TF_FIELD_CLUSTER_TOPOLOGY) {
	//	// Could be one of a few scenarios...
	//
	//	// Scenario #2 - Started with 1 node; Then, added a replica and specified the topology for future purposes (Scenario #1)
	//	if isHostChanged {
	//		// Noop: Just ignore. Nothing to do. This could possibly be on e of the following:
	//		// A. The addition of a new replica and the user has specified proper Master=>Slave links
	//		// B. The removal of a replica and the user has updated (or, if no replicas are present any longer, then completely removed) the topology def
	//	} else {
	//		// Scenario #1 - Role change: Slave promoted to Master
	//	}
	//}

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

	backupFailChoice := d.Get(TF_FIELD_BACKUP_FAILOVER).(bool)
	jobData.SetBackupFailover(backupFailChoice)

	failoverHost := d.Get(TF_FIELD_BACKUP_FAILOVER_HOST).(string)
	if failoverHost == "" {
		failoverHost = STINRG_AUTO
	}
	jobData.SetBackupFailoverHost(failoverHost)

	backupStorageHost := d.Get(TF_FIELD_BACKUP_STORAGE_HOST).(string)
	if backupStorageHost != "" {
		jobData.SetStorageHost(backupStorageHost)
	}

	return err
}

func (c *MySQLMaria) IsValidBackupOptions(vendor string, clusterType string, jobData *openapi.JobsJobJobSpecJobData) error {
	return c.Backup.IsValidBackupOptions(vendor, clusterType, jobData)
}

func (c *MySQLMaria) SetBackupJobData(jobData *openapi.JobsJobJobSpecJobData) error {
	return c.Backup.SetBackupJobData(jobData)
}

func NewMySQLMaria() *MySQLMaria {
	return &MySQLMaria{}
}
