package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/severalnines/clustercontrol-client-sdk/go/pkg/openapi"
	"log/slog"
	"strconv"
)

type MongoDb struct {
	common DbCommon
}

func (m *MongoDb) GetInputs(d *schema.ResourceData, jobData *openapi.JobsJobJobSpecJobData) error {
	funcName := "Mongo::GetInputs"
	slog.Debug(funcName)

	var err error

	// parent/super - get common attributes
	if err = m.common.GetInputs(d, jobData); err != nil {
		return err
	}

	getReplicasets(d, jobData)

	getConfigServers(d, jobData)

	getMongosServers(d, jobData)

	return nil
}

func getConfigServers(d *schema.ResourceData, jobData *openapi.JobsJobJobSpecJobData) error {
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
			getCommonHostAttributes(memFromTF, iPort, clusterType, memHost)
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

func getMongosServers(d *schema.ResourceData, jobData *openapi.JobsJobJobSpecJobData) error {
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
		getCommonHostAttributes(f, iPort, clusterType, memHost)

		mongos = append(mongos, mem)
	}
	//fmt.Fprintf(os.Stderr, "getMongosServers: %v\n", mongos)
	jobData.SetMongosServers(mongos)

	return nil
}

func getReplicasets(d *schema.ResourceData, jobData *openapi.JobsJobJobSpecJobData) error {
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
			getCommonHostAttributes(memFromTF, iPort, clusterType, memHost)

			priority := memFromTF[TF_FIELD_CLUSTER_HOST_PRIORITY].(string)
			mem.SetPriority(priority)
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

func NewMongo() *MongoDb {
	return &MongoDb{}
}
