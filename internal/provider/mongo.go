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

type MongoDb struct {
	Common DbCommon
	Backup DbBackupCommon
}

func (c *MongoDb) HostClassRoleAndReplicasetCompare(one *openapi.JobsJobJobSpecJobDataNodesInner, two *openapi.ClusterResponseHostsInner, options ...string) bool {
	// TODO

	//ii := 0
	//hostClass := ""
	//hostRole := ""
	//replicasetName := ""
	//if len(options) > ii {
	//	hostClass = options[ii]
	//	ii++
	//}
	//if len(options) > ii {
	//	hostRole = options[ii]
	//	ii++
	//}
	//if len(options) > ii {
	//	replicasetName = options[ii]
	//	ii++
	//}

	return strings.EqualFold(one.GetHostname(), two.GetHostname())
}

func (m *MongoDb) GetInputs(d *schema.ResourceData, jobData *openapi.JobsJobJobSpecJobData) error {
	funcName := "Mongo::GetInputs"
	slog.Debug(funcName)

	var err error

	// parent/super - get common attributes
	if err = m.Common.GetInputs(d, jobData); err != nil {
		return err
	}

	clusterType := jobData.GetClusterType()

	var iPort int
	port := d.Get(TF_FIELD_CLUSTER_MONGODB_PORT).(string)
	if err = CheckForEmptyAndSetDefault(&port, gDefultDbPortMap, clusterType); err != nil {
		return err
	}
	if iPort, err = strconv.Atoi(port); err != nil {
		slog.Error(funcName, "ERROR", "Non-numeric database port")
		return err
	}
	jobData.SetPort(int32(iPort))
	topLevelPort := jobData.GetPort()

	configServerPort := d.Get(TF_FIELD_CLUSTER_MONGODB_CFG_SRVR_PORT).(string)
	if configServerPort == "" {
		configServerPort = DEFAULT_MONGO_CONFIG_SRVR_PORT
	}
	//jobData.SetConfigServerPort(configServerPort)

	dbVendor := jobData.GetVendor()
	dbVersion := jobData.GetVersion()

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

	authDb := d.Get(TF_FIELD_CLUSTER_MONGO_AUTH_DB).(string)
	if authDb == "" {
		authDb = MONGO_DEFAULT_AUTH_DB
	}
	jobData.SetMongodbAuthdb(authDb)

	mongosTemplate, ok := gDbMongosConfigTemplate[dbVendor]
	if ok && mongosTemplate != "" {
		jobData.SetMongosConfTemplate(mongosTemplate)
	}

	//********************************************
	// Get Mongo - Replicasets
	//********************************************
	//m.getReplicasets(d, jobData)
	replicaSetsFromTF := d.Get(TF_FIELD_CLUSTER_REPLICA_SET)
	replicaSets := []openapi.JobsJobJobSpecJobDataReplicaSetsInner{}
	for _, replicaSetFromTf := range replicaSetsFromTF.([]any) {
		rsFromTF := replicaSetFromTf.(map[string]any)
		rs := rsFromTF[TF_FIELD_CLUSTER_REPLICA_SET_RS].(string)
		membersFromTF := rsFromTF[TF_FIELD_CLUSTER_REPLICA_MEMBER]
		members := []openapi.JobsJobJobSpecJobDataReplicaSetsInnerMembersInner{}
		memberNum := 0
		for _, memberFromTF := range membersFromTF.([]any) {
			memFromTF := memberFromTF.(map[string]any)

			hostname := memFromTF[TF_FIELD_CLUSTER_HOSTNAME].(string)
			hostname_data := memFromTF[TF_FIELD_CLUSTER_HOSTNAME_DATA].(string)
			hostname_internal := memFromTF[TF_FIELD_CLUSTER_HOSTNAME_INT].(string)
			//port := memFromTF[TF_FIELD_CLUSTER_HOST_PORT].(string)

			if hostname == "" {
				// Not specifying hostname is disallowed !!!
				return errors.New("Hostname attrbute must be set in replicaset member.")
			}
			var mem = openapi.JobsJobJobSpecJobDataReplicaSetsInnerMembersInner{
				Hostname: &hostname,
			}
			if hostname_data != "" {
				mem.SetHostnameData(hostname_data)
			}
			if hostname_internal != "" {
				mem.SetHostnameInternal(hostname_internal)
			}
			mem.SetPort(strconv.Itoa(int(topLevelPort)))
			//if port == "" {
			//	mem.SetPort(strconv.Itoa(int(topLevelPort)))
			//} else {
			//	mem.SetPort(strconv.Itoa(int(convertPortToInt(port, topLevelPort))))
			//}
			arbiter_only := memFromTF[TF_FIELD_CLUSTER_HOST_ARBITER_ONLY].(bool)
			mem.SetArbiterOnly(arbiter_only)

			if memberNum > 0 {
				priority := memFromTF[TF_FIELD_CLUSTER_HOST_PRIORITY].(int)
				if priority == 0 {
					priority = 1
				}
				mem.SetPriority(int32(priority))
				//if memFromTF[TF_FIELD_CLUSTER_HOST_PRIORITY] != nil {
				//	priority := memFromTF[TF_FIELD_CLUSTER_HOST_PRIORITY].(int32)
				//	mem.SetPriority(priority)
				//}
				slave_delay := memFromTF[TF_FIELD_CLUSTER_HOST_SLAVE_DELAY].(string)
				if slave_delay != "" {
					mem.SetSlaveDelay(slave_delay)
				}
				hidden := memFromTF[TF_FIELD_CLUSTER_HOST_HIDDEN].(bool)
				mem.SetHidden(hidden)
			}

			members = append(members, mem)

			memberNum++
		}
		var node = openapi.JobsJobJobSpecJobDataReplicaSetsInner{
			Rs:      &rs,
			Members: members,
		}

		replicaSets = append(replicaSets, node)
	}
	jobData.SetReplicaSets(replicaSets)

	//********************************************
	// Get Mongo - Config Servers
	//********************************************
	configServerFromTF := d.Get(TF_FIELD_CLUSTER_MONGO_CONFIG_SERVER)
	for _, cfgServerFromTF := range configServerFromTF.([]any) {
		cfgFromTF := cfgServerFromTF.(map[string]any)
		rs := cfgFromTF[TF_FIELD_CLUSTER_REPLICA_SET_RS].(string)
		membersFromTF := cfgFromTF[TF_FIELD_CLUSTER_REPLICA_MEMBER]
		members := []openapi.JobsJobJobSpecJobDataConfigServersMembersInner{}
		for _, memberFromTF := range membersFromTF.([]any) {
			memFromTF := memberFromTF.(map[string]any)

			hostname := memFromTF[TF_FIELD_CLUSTER_HOSTNAME].(string)
			hostname_data := memFromTF[TF_FIELD_CLUSTER_HOSTNAME_DATA].(string)
			hostname_internal := memFromTF[TF_FIELD_CLUSTER_HOSTNAME_INT].(string)
			//port := memFromTF[TF_FIELD_CLUSTER_HOST_PORT].(string)

			if hostname == "" {
				return errors.New("Mongo config server hostname cannot be empty.")
			}
			var mem = openapi.JobsJobJobSpecJobDataConfigServersMembersInner{
				Hostname: &hostname,
			}
			if hostname_data != "" {
				mem.SetHostnameData(hostname_data)
			}
			if hostname_internal != "" {
				mem.SetHostnameInternal(hostname_internal)
			}
			//mem.SetPort(strconv.Itoa(iPort))
			mem.SetPort(configServerPort)
			//if port == "" {
			//	mem.SetPort(strconv.Itoa(iPort))
			//} else {
			//	mem.SetPort(strconv.Itoa(int(convertPortToInt(port, int32(iPort)))))
			//}
			members = append(members, mem)
		}
		var cfgSrvr = openapi.JobsJobJobSpecJobDataConfigServers{
			Rs:      &rs,
			Members: members,
		}
		jobData.SetConfigServers(cfgSrvr)

		// There should only be one entry here. Therefore, get out of here!
		break
	}

	//********************************************
	// Get Mongo - Mongos servers
	//********************************************
	//m.getMongosServers(d, jobData)
	mongosServersFromTF := d.Get(TF_FIELD_CLUSTER_MONGOS_SERVER)
	mongos := []openapi.JobsJobJobSpecJobDataConfigServersMembersInner{}
	for _, ff := range mongosServersFromTF.([]any) {
		f := ff.(map[string]any)
		hostname := f[TF_FIELD_CLUSTER_HOSTNAME].(string)
		hostname_data := f[TF_FIELD_CLUSTER_HOSTNAME_DATA].(string)
		hostname_internal := f[TF_FIELD_CLUSTER_HOSTNAME_INT].(string)
		//port := f[TF_FIELD_CLUSTER_HOST_PORT].(string)

		if hostname == "" {
			return errors.New("Mongo config server hostname cannot be empty.")
		}
		var mem = openapi.JobsJobJobSpecJobDataConfigServersMembersInner{
			Hostname: &hostname,
		}
		if hostname_data != "" {
			mem.SetHostnameData(hostname_data)
		}
		if hostname_internal != "" {
			mem.SetHostnameInternal(hostname_internal)
		}
		mem.SetPort(strconv.Itoa(int(topLevelPort)))
		//if port == "" {
		//	mem.SetPort(strconv.Itoa(int(topLevelPort)))
		//} else {
		//	mem.SetPort(strconv.Itoa(int(convertPortToInt(port, topLevelPort))))
		//}
		mongos = append(mongos, mem)
	}
	jobData.SetMongosServers(mongos)

	return nil
}

func (c *MongoDb) HandleRead(ctx context.Context, d *schema.ResourceData, m interface{}, clusterInfo *openapi.ClusterResponse) error {

	if err := c.Common.HandleRead(ctx, d, m, clusterInfo); err != nil {
		return err
	}

	return nil
}

func (c *MongoDb) IsUpdateBatchAllowed(d *schema.ResourceData) error {
	var err error

	if err = c.Common.IsUpdateBatchAllowed(d); err != nil {
		return err
	}

	//updateClassA := d.HasChange(TF_FIELD_CLUSTER_HOST)
	//updateClassAprime := d.HasChangeExcept(TF_FIELD_CLUSTER_HOST)
	//if updateClassA && updateClassAprime {
	//	err = errors.New(fmt.Sprintf("You are not allowed to update %s along with any other fields.", TF_FIELD_CLUSTER_HOST))
	//	return err
	//}

	updateClassA := d.HasChange(TF_FIELD_CLUSTER_ENABLE_PGM_AGENT)
	updateClassAprime := d.HasChangeExcept(TF_FIELD_CLUSTER_ENABLE_PGM_AGENT)
	if updateClassA && updateClassAprime {
		err = errors.New(fmt.Sprintf("You are not allowed to update %s along with any other fields.", TF_FIELD_CLUSTER_ENABLE_PGM_AGENT))
		return err
	}

	updateClassA = d.HasChange(TF_FIELD_CLUSTER_REPLICA_SET)
	updateClassAprime = d.HasChangeExcept(TF_FIELD_CLUSTER_REPLICA_SET)
	if updateClassA && updateClassAprime {
		err = errors.New(fmt.Sprintf("You are not allowed to update %s along with any other fields.", TF_FIELD_CLUSTER_REPLICA_SET))
		return err
	}

	return nil
}

func (c *MongoDb) HandleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}, clusterInfo *openapi.ClusterResponse) error {
	funcName := "MongoDb::HandleUpdate"
	slog.Debug(funcName)

	var err error

	if err := c.Common.HandleUpdate(ctx, d, m, clusterInfo); err != nil {
		return err
	}

	// handle other big updates such as add-replication-slave, remove-node, etc
	tmpJobData := openapi.NewJobsJobJobSpecJobData()
	if err = c.GetInputs(d, tmpJobData); err != nil {
		return err
	}

	if d.HasChange(TF_FIELD_CLUSTER_ENABLE_PGM_AGENT) {
		enablePbm := d.Get(TF_FIELD_CLUSTER_ENABLE_PGM_AGENT).(bool)
		backupDir := d.Get(TF_FIELD_CLUSTER_PBM_BACKUP_DIR).(string)
		if !enablePbm {
			return nil
		}

		if backupDir == "" {
			return errors.New("db_backup_dir must be set. It is an nfs mounted fs on all mongo hosts.")
		}

		apiClient := m.(*openapi.APIClient)

		enablePbmJob := NewCCJob(CMON_JOB_CREATE_JOB)
		job := enablePbmJob.GetJob()
		jobSpec := job.GetJobSpec()
		jobData := jobSpec.GetJobData()
		enablePbmJob.SetClusterId(clusterInfo.GetClusterId())
		jobSpec.SetCommand(CMON_JOB_PBM_AGENT_COMMAND)
		jobData.SetAction(JOB_ACTION_SETUP)

		var nodes = []openapi.JobsJobJobSpecJobDataNodesInner{}
		var node = openapi.JobsJobJobSpecJobDataNodesInner{}
		node.SetClassName(CMON_CLASS_NAME_PBM_AGENT_HOST)
		node.SetHostname("*")
		node.SetBackupDir(backupDir)
		nodes = append(nodes, node)
		jobData.SetNodes(nodes)

		jobSpec.SetJobData(jobData)
		job.SetJobSpec(jobSpec)
		enablePbmJob.SetJob(job)

		if err = SendAndWaitForJobCompletion(ctx, apiClient, enablePbmJob); err != nil {
			slog.Error(err.Error())
		}

	}

	if d.HasChange(TF_FIELD_CLUSTER_REPLICA_SET) {
		apiClient := m.(*openapi.APIClient)
		addOrRemoveNodeJob := NewCCJob(CMON_JOB_CREATE_JOB)
		addOrRemoveNodeJob.SetClusterId(clusterInfo.GetClusterId())
		job := addOrRemoveNodeJob.GetJob()
		jobSpec := job.GetJobSpec()
		jobData := jobSpec.GetJobData()

		var rsMembersToAdd []openapi.JobsJobJobSpecJobDataReplicaSetsInner
		var rsMembersToRemove []openapi.JobsJobJobSpecJobDataReplicaSetsInner

		// Compare Terraform and CMON to determine whether adding node, remove node or promoting standby/slave
		// The logic here basically is...
		// 1. Get the list of replicasets from TF decleration
		// 2. Get a list of the member hosts for each of the replicaset

		// Get the list of replicasets from TF decleration
		replicaSets, _ := c.getReplicasetHosts(d)

		// Get a list of the member hosts for each of the replicasets
		additions := 0
		removals := 0
		for _, replicaSet := range replicaSets {
			// accumulate the member hosts in a form that can be compared with what's currently available in CMON
			var nodesToAdd []openapi.JobsJobJobSpecJobDataNodesInner
			var nodesToRemove []openapi.JobsJobJobSpecJobDataNodesInner
			rsMembers := replicaSet.GetMembers()
			var nodes []openapi.JobsJobJobSpecJobDataNodesInner
			for _, rsMember := range rsMembers {
				var node openapi.JobsJobJobSpecJobDataNodesInner
				node.SetHostname(rsMember.GetHostname())
				// accumulating...
				nodes = append(nodes, node)
			}
			// For each replicaset, compare with what's currently available in CMON.
			// Result: either nodes need to be added to CMON or removed from CMON
			if nodesToAdd, nodesToRemove, err = c.Common.determineNodesDelta(nodes, clusterInfo,
				CMON_CLASS_NAME_MONGO_HOST, CMON_DB_HOST_ROLE_MONGO_SHARD_SERVER, replicaSet.GetRs()); err != nil {
				return err
			}
			// convert back to replicaset format
			var rsMemberToAdd = openapi.JobsJobJobSpecJobDataReplicaSetsInner{}
			rsMemberToAdd.SetRs(replicaSet.GetRs())
			var addMembers []openapi.JobsJobJobSpecJobDataReplicaSetsInnerMembersInner
			for _, nodeToAdd := range nodesToAdd {
				var mem = openapi.JobsJobJobSpecJobDataReplicaSetsInnerMembersInner{}
				mem.SetHostname(nodeToAdd.GetHostname())
				addMembers = append(addMembers, mem)
				additions++
			}
			rsMemberToAdd.SetMembers(addMembers)
			rsMembersToAdd = append(rsMembersToAdd, rsMemberToAdd)

			var rsMemberToRemove = openapi.JobsJobJobSpecJobDataReplicaSetsInner{}
			rsMemberToRemove.SetRs(replicaSet.GetRs())
			var removeMembers []openapi.JobsJobJobSpecJobDataReplicaSetsInnerMembersInner
			for _, nodeToRemove := range nodesToRemove {
				var mem = openapi.JobsJobJobSpecJobDataReplicaSetsInnerMembersInner{}
				mem.SetHostname(nodeToRemove.GetHostname())
				removeMembers = append(removeMembers, mem)
				removals++
			}
			rsMemberToRemove.SetMembers(removeMembers)
			rsMembersToRemove = append(rsMembersToRemove, rsMemberToRemove)
		}

		isAdd := additions > 0
		isRemove := removals > 0

		if additions > 1 || removals > 1 {
			return errors.New("Can only Add/Remove one node at-a-time.")
		}

		var rsNodeToAddOrRemove *openapi.JobsJobJobSpecJobDataReplicaSetsInner
		if isAdd {
			jobSpec.SetCommand(CMON_JOB_ADD_NODE_COMMAND)
			rsNodeToAddOrRemove = &rsMembersToAdd[0]
			jobData = *tmpJobData
			var rsMemNode = openapi.JobsJobJobSpecJobDataNode{}
			t := rsNodeToAddOrRemove.GetMembers()
			rsMemNode.SetHostname(t[0].GetHostname())
			rsMemNode.SetPort(tmpJobData.GetPort())
			//if t[0].GetPort() == "" {
			//	rsMemNode.SetPort(tmpJobData.GetPort())
			//} else {
			//	iP, _ := strconv.Atoi(t[0].GetPort())
			//	rsMemNode.SetPort(int32(iP))
			//}
			jobData.SetReplicaset(rsNodeToAddOrRemove.GetRs())
			jobData.SetNodeType(0)
			jobData.SetNode(rsMemNode)
			var cfgServers = openapi.JobsJobJobSpecJobDataConfigServers{}
			jobData.SetConfigServers(cfgServers)
			var mongosServers = []openapi.JobsJobJobSpecJobDataConfigServersMembersInner{}
			jobData.SetMongosServers(mongosServers)
			var emptyReplicasets []openapi.JobsJobJobSpecJobDataReplicaSetsInner
			jobData.SetReplicaSets(emptyReplicasets)
			slog.Info(funcName, "Adding hostname", rsMemNode.GetHostname())
		} else if isRemove {
			jobSpec.SetCommand(CMON_JOB_REMOVE_NODE_COMMAND)
			rsNodeToAddOrRemove = &rsMembersToRemove[0]
			t := rsNodeToAddOrRemove.GetMembers()
			var node openapi.JobsJobJobSpecJobDataNode
			node.SetHostname(t[0].GetHostname())
			node.SetPort(tmpJobData.GetPort())
			jobData.SetNode(node)
			jobData.SetEnableUninstall(true)
			jobData.SetUnregisterOnly(false)
			slog.Debug(funcName, "Removing hostname", node.GetHostname())
		} else {
			return nil
		}

		jobSpec.SetJobData(jobData)
		job.SetJobSpec(jobSpec)
		addOrRemoveNodeJob.SetJob(job)

		if err = SendAndWaitForJobCompletion(ctx, apiClient, addOrRemoveNodeJob); err != nil {
			slog.Error(err.Error())
		}

	}

	return nil
}

func (m *MongoDb) GetBackupInputs(d *schema.ResourceData, jobData *openapi.JobsJobJobSpecJobData) error {
	funcName := "MongoDb::GetBackupInputs"
	slog.Debug(funcName)

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
	if backupStorageHost == "" {
		if strings.EqualFold(jobData.GetBackupMethod(), BACKUP_METHOD_MONGODUMP) {
			return errors.New("db_backup_storage_host must be set for mongodump")
		}
	} else {
		jobData.SetStorageHost(backupStorageHost)
	}

	return err
}

func (c *MongoDb) IsValidBackupOptions(vendor string, clusterType string, jobData *openapi.JobsJobJobSpecJobData) error {
	return c.Backup.IsValidBackupOptions(vendor, clusterType, jobData)
}

func (c *MongoDb) SetBackupJobData(jobData *openapi.JobsJobJobSpecJobData) error {
	return c.Backup.SetBackupJobData(jobData)
}

func (c *MongoDb) IsBackupRemovable(clusterInfo *openapi.ClusterResponse, jobData *openapi.JobsJobJobSpecJobData) bool {
	return true
}

func (c *MongoDb) getReplicasetHosts(d *schema.ResourceData) ([]openapi.JobsJobJobSpecJobDataReplicaSetsInner, error) {
	var nodes []openapi.JobsJobJobSpecJobDataReplicaSetsInner

	hosts := d.Get(TF_FIELD_CLUSTER_REPLICA_SET)
	for _, ff := range hosts.([]any) {
		f := ff.(map[string]any)
		rsName := f[TF_FIELD_CLUSTER_REPLICA_SET_RS].(string)
		var node = openapi.JobsJobJobSpecJobDataReplicaSetsInner{
			Rs: &rsName,
		}
		var rsMembers []openapi.JobsJobJobSpecJobDataReplicaSetsInnerMembersInner
		// Capturing hostname only as this is only used for comparison purposes.
		myHost := f[TF_FIELD_CLUSTER_REPLICA_MEMBER]
		for _, tt := range myHost.([]any) {
			t := tt.(map[string]any)
			hostname := t[TF_FIELD_CLUSTER_HOSTNAME].(string)
			if hostname == "" {
				continue
			}
			var rsMem = openapi.JobsJobJobSpecJobDataReplicaSetsInnerMembersInner{
				Hostname: &hostname,
			}
			rsMembers = append(rsMembers, rsMem)
		}
		node.SetMembers(rsMembers)
		nodes = append(nodes, node)
	}

	return nodes, nil
}

func NewMongo() *MongoDb {
	return &MongoDb{}
}
