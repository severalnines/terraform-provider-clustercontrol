package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/severalnines/clustercontrol-client-sdk/go/pkg/openapi"
	"log/slog"
	"strconv"
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

	m.getReplicasets(d, jobData)

	m.getConfigServers(d, jobData)

	m.getMongosServers(d, jobData)

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

func (c *MongoDb) getConfigServers(d *schema.ResourceData, jobData *openapi.JobsJobJobSpecJobData) error {
	clusterType := jobData.GetClusterType()
	iPort, _ := strconv.Atoi(DEFAULT_MONGO_CONFIG_SRVR_PORT)

	configServerFromTF := d.Get(TF_FIELD_CLUSTER_MONGO_CONFIG_SERVER)
	for _, cfgServerFromTF := range configServerFromTF.([]any) {
		cfgFromTF := cfgServerFromTF.(map[string]any)
		rs := cfgFromTF[TF_FIELD_CLUSTER_REPLICA_SET_RS].(string)
		membersFromTF := cfgFromTF[TF_FIELD_CLUSTER_REPLICA_MEMBER]
		members := []openapi.JobsJobJobSpecJobDataConfigServersMembersInner{}
		for _, memberFromTF := range membersFromTF.([]any) {
			memFromTF := memberFromTF.(map[string]any)

			var mem = openapi.JobsJobJobSpecJobDataConfigServersMembersInner{}
			var memHost = memberHosts{
				mongoCfgNode: &mem,
			}
			c.Common.getCommonHostAttributes(memFromTF, iPort, clusterType, memHost)
			members = append(members, mem)
		}
		var cfgSrvr = openapi.JobsJobJobSpecJobDataConfigServers{
			Rs:      &rs,
			Members: members,
		}
		//fmt.Fprintf(os.Stderr, "getConfigServers: %v\n", cfgSrvr)
		jobData.SetConfigServers(cfgSrvr)

		// There should only be one entry here. Therefore, get out of the outer loop
		break
	}

	return nil
}

func (c *MongoDb) getMongosServers(d *schema.ResourceData, jobData *openapi.JobsJobJobSpecJobData) error {
	iPort := int(jobData.GetPort())
	clusterType := jobData.GetClusterType()

	mongosServersFromTF := d.Get(TF_FIELD_CLUSTER_MONGOS_SERVER)
	mongos := []openapi.JobsJobJobSpecJobDataConfigServersMembersInner{}
	for _, ff := range mongosServersFromTF.([]any) {
		f := ff.(map[string]any)

		var mem = openapi.JobsJobJobSpecJobDataConfigServersMembersInner{}
		var memHost = memberHosts{
			mongoCfgNode: &mem,
		}
		c.Common.getCommonHostAttributes(f, iPort, clusterType, memHost)

		mongos = append(mongos, mem)
	}
	//fmt.Fprintf(os.Stderr, "getMongosServers: %v\n", mongos)
	jobData.SetMongosServers(mongos)

	return nil
}

func (c *MongoDb) getReplicasets(d *schema.ResourceData, jobData *openapi.JobsJobJobSpecJobData) error {
	iPort := int(jobData.GetPort())
	clusterType := jobData.GetClusterType()

	replicaSetsFromTF := d.Get(TF_FIELD_CLUSTER_REPLICA_SET)
	replicaSets := []openapi.JobsJobJobSpecJobDataReplicaSetsInner{}
	for _, replicaSetFromTf := range replicaSetsFromTF.([]any) {
		rsFromTF := replicaSetFromTf.(map[string]any)
		rs := rsFromTF[TF_FIELD_CLUSTER_REPLICA_SET_RS].(string)
		membersFromTF := rsFromTF[TF_FIELD_CLUSTER_REPLICA_MEMBER]
		members := []openapi.JobsJobJobSpecJobDataReplicaSetsInnerMembersInner{}
		for _, memberFromTF := range membersFromTF.([]any) {
			memFromTF := memberFromTF.(map[string]any)

			var mem = openapi.JobsJobJobSpecJobDataReplicaSetsInnerMembersInner{}
			var memHost = memberHosts{
				mongoRsNode: &mem,
			}
			c.Common.getCommonHostAttributes(memFromTF, iPort, clusterType, memHost)

			if memFromTF[TF_FIELD_CLUSTER_HOST_PRIORITY] != nil {
				priority := memFromTF[TF_FIELD_CLUSTER_HOST_PRIORITY].(int32)
				mem.SetPriority(priority)
			}
			slave_delay := memFromTF[TF_FIELD_CLUSTER_HOST_SLAVE_DELAY].(string)
			mem.SetSlaveDelay(slave_delay)
			arbiter_only := memFromTF[TF_FIELD_CLUSTER_HOST_ARBITER_ONLY].(bool)
			mem.SetArbiterOnly(arbiter_only)
			hidden := memFromTF[TF_FIELD_CLUSTER_HOST_HIDDEN].(bool)
			mem.SetHidden(hidden)

			//fmt.Fprintf(os.Stderr, "getReplicasets: %s - %s\n", hostname, port)
			members = append(members, mem)
		}
		var node = openapi.JobsJobJobSpecJobDataReplicaSetsInner{
			Rs:      &rs,
			Members: members,
		}

		replicaSets = append(replicaSets, node)
	}
	//fmt.Fprintf(os.Stderr, "getReplicasets: %v\n", replicaSets)
	jobData.SetReplicaSets(replicaSets)

	return nil
}

func (m *MongoDb) GetBackupInputs(d *schema.ResourceData, jobData *openapi.JobsJobJobSpecJobData) error {
	funcName := "MongoDb::GetBackupInputs"
	slog.Info(funcName)

	var err error

	// parent/super - get common attributes
	if err = m.Backup.GetBackupInputs(d, jobData); err != nil {
		return err
	}

	return err
}

func (c *MongoDb) IsValidBackupOptions(vendor string, clusterType string, jobData *openapi.JobsJobJobSpecJobData) error {
	return c.Backup.IsValidBackupOptions(vendor, clusterType, jobData)
}

func NewMongo() *MongoDb {
	return &MongoDb{}
}
