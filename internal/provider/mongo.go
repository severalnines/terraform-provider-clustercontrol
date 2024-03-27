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

func (m *MongoDb) GetInputs(d *schema.ResourceData, jobData *openapi.JobsJobJobSpecJobData) error {
	funcName := "Mongo::GetInputs"
	slog.Debug(funcName)

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
			port := memFromTF[TF_FIELD_CLUSTER_HOST_PORT].(string)

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
			if port == "" {
				mem.SetPort(strconv.Itoa(int(topLevelPort)))
			} else {
				mem.SetPort(strconv.Itoa(int(convertPortToInt(port, topLevelPort))))
			}
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
	//m.getConfigServers(d, jobData)
	iPort, _ := strconv.Atoi(DEFAULT_MONGO_CONFIG_SRVR_PORT)
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
			port := memFromTF[TF_FIELD_CLUSTER_HOST_PORT].(string)

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
			if port == "" {
				mem.SetPort(strconv.Itoa(int(topLevelPort)))
			} else {
				mem.SetPort(strconv.Itoa(int(convertPortToInt(port, int32(iPort)))))
			}
			//var memHost = memberHosts{
			//	mongoCfgNode: &mem,
			//}
			//c.Common.getCommonHostAttributes(memFromTF, iPort, clusterType, memHost)
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
		port := f[TF_FIELD_CLUSTER_HOST_PORT].(string)

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
		if port == "" {
			mem.SetPort(strconv.Itoa(int(topLevelPort)))
		} else {
			mem.SetPort(strconv.Itoa(int(convertPortToInt(port, topLevelPort))))
		}
		//var memHost = memberHosts{
		//	mongoCfgNode: &mem,
		//}
		//c.Common.getCommonHostAttributes(f, iPort, clusterType, memHost)

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

	return nil
}

func (c *MongoDb) HandleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}, clusterInfo *openapi.ClusterResponse) error {

	if err := c.Common.HandleUpdate(ctx, d, m, clusterInfo); err != nil {
		return err
	}

	return nil
}

//func (c *MongoDb) getConfigServers(d *schema.ResourceData, jobData *openapi.JobsJobJobSpecJobData) error {
//	clusterType := jobData.GetClusterType()
//	iPort, _ := strconv.Atoi(DEFAULT_MONGO_CONFIG_SRVR_PORT)
//
//	configServerFromTF := d.Get(TF_FIELD_CLUSTER_MONGO_CONFIG_SERVER)
//	for _, cfgServerFromTF := range configServerFromTF.([]any) {
//		cfgFromTF := cfgServerFromTF.(map[string]any)
//		rs := cfgFromTF[TF_FIELD_CLUSTER_REPLICA_SET_RS].(string)
//		membersFromTF := cfgFromTF[TF_FIELD_CLUSTER_REPLICA_MEMBER]
//		members := []openapi.JobsJobJobSpecJobDataConfigServersMembersInner{}
//		for _, memberFromTF := range membersFromTF.([]any) {
//			memFromTF := memberFromTF.(map[string]any)
//
//			var mem = openapi.JobsJobJobSpecJobDataConfigServersMembersInner{}
//			var memHost = memberHosts{
//				mongoCfgNode: &mem,
//			}
//			c.Common.getCommonHostAttributes(memFromTF, iPort, clusterType, memHost)
//			members = append(members, mem)
//		}
//		var cfgSrvr = openapi.JobsJobJobSpecJobDataConfigServers{
//			Rs:      &rs,
//			Members: members,
//		}
//		//fmt.Fprintf(os.Stderr, "getConfigServers: %v\n", cfgSrvr)
//		jobData.SetConfigServers(cfgSrvr)
//
//		// There should only be one entry here. Therefore, get out of the outer loop
//		break
//	}
//
//	return nil
//}
//
//func (c *MongoDb) getMongosServers(d *schema.ResourceData, jobData *openapi.JobsJobJobSpecJobData) error {
//	iPort := int(jobData.GetPort())
//	clusterType := jobData.GetClusterType()
//
//	mongosServersFromTF := d.Get(TF_FIELD_CLUSTER_MONGOS_SERVER)
//	mongos := []openapi.JobsJobJobSpecJobDataConfigServersMembersInner{}
//	for _, ff := range mongosServersFromTF.([]any) {
//		f := ff.(map[string]any)
//
//		var mem = openapi.JobsJobJobSpecJobDataConfigServersMembersInner{}
//		var memHost = memberHosts{
//			mongoCfgNode: &mem,
//		}
//		c.Common.getCommonHostAttributes(f, iPort, clusterType, memHost)
//
//		mongos = append(mongos, mem)
//	}
//	//fmt.Fprintf(os.Stderr, "getMongosServers: %v\n", mongos)
//	jobData.SetMongosServers(mongos)
//
//	return nil
//}
//
//func (c *MongoDb) getReplicasets(d *schema.ResourceData, jobData *openapi.JobsJobJobSpecJobData) error {
//	iPort := int(jobData.GetPort())
//	clusterType := jobData.GetClusterType()
//
//	replicaSetsFromTF := d.Get(TF_FIELD_CLUSTER_REPLICA_SET)
//	replicaSets := []openapi.JobsJobJobSpecJobDataReplicaSetsInner{}
//	for _, replicaSetFromTf := range replicaSetsFromTF.([]any) {
//		rsFromTF := replicaSetFromTf.(map[string]any)
//		rs := rsFromTF[TF_FIELD_CLUSTER_REPLICA_SET_RS].(string)
//		membersFromTF := rsFromTF[TF_FIELD_CLUSTER_REPLICA_MEMBER]
//		members := []openapi.JobsJobJobSpecJobDataReplicaSetsInnerMembersInner{}
//		for _, memberFromTF := range membersFromTF.([]any) {
//			memFromTF := memberFromTF.(map[string]any)
//
//			var mem = openapi.JobsJobJobSpecJobDataReplicaSetsInnerMembersInner{}
//			var memHost = memberHosts{
//				mongoRsNode: &mem,
//			}
//			c.Common.getCommonHostAttributes(memFromTF, iPort, clusterType, memHost)
//
//			if memFromTF[TF_FIELD_CLUSTER_HOST_PRIORITY] != nil {
//				priority := memFromTF[TF_FIELD_CLUSTER_HOST_PRIORITY].(int32)
//				mem.SetPriority(priority)
//			}
//			slave_delay := memFromTF[TF_FIELD_CLUSTER_HOST_SLAVE_DELAY].(string)
//			mem.SetSlaveDelay(slave_delay)
//			arbiter_only := memFromTF[TF_FIELD_CLUSTER_HOST_ARBITER_ONLY].(bool)
//			mem.SetArbiterOnly(arbiter_only)
//			hidden := memFromTF[TF_FIELD_CLUSTER_HOST_HIDDEN].(bool)
//			mem.SetHidden(hidden)
//
//			//fmt.Fprintf(os.Stderr, "getReplicasets: %s - %s\n", hostname, port)
//			members = append(members, mem)
//		}
//		var node = openapi.JobsJobJobSpecJobDataReplicaSetsInner{
//			Rs:      &rs,
//			Members: members,
//		}
//
//		replicaSets = append(replicaSets, node)
//	}
//	//fmt.Fprintf(os.Stderr, "getReplicasets: %v\n", replicaSets)
//	jobData.SetReplicaSets(replicaSets)
//
//	return nil
//}

func (m *MongoDb) GetBackupInputs(d *schema.ResourceData, jobData *openapi.JobsJobJobSpecJobData) error {
	funcName := "MongoDb::GetBackupInputs"
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

func NewMongo() *MongoDb {
	return &MongoDb{}
}
