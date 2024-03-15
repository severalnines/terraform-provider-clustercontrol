package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/severalnines/clustercontrol-client-sdk/go/pkg/openapi"
	"log/slog"
	"strconv"
)

type DbCommon struct{}

func (c *DbCommon) GetInputs(d *schema.ResourceData, jobData *openapi.JobsJobJobSpecJobData) error {
	funcName := "DbCommon::GetInputs"
	slog.Info(funcName)
	//fmt.Fprintf(os.Stderr, "%s", funcName)

	var err error

	clusterName := d.Get(TF_FIELD_CLUSTER_NAME).(string)
	jobData.SetClusterName(clusterName)

	clusterType := d.Get(TF_FIELD_CLUSTER_TYPE).(string)
	jobData.SetClusterType(clusterType)

	dbVendor := d.Get(TF_FIELD_CLUSTER_VENDOR).(string)
	jobData.SetVendor(dbVendor)

	dbVersion := d.Get(TF_FIELD_CLUSTER_VERSION).(string)
	jobData.SetVersion(dbVersion)

	dbAdminUsername := d.Get(TF_FIELD_CLUSTER_ADMIN_USER).(string)
	if err = CheckForEmptyAndSetDefault(&dbAdminUsername, gDefultDbAdminUser, clusterType); err != nil {
		return err
	}
	jobData.SetDbUser(dbAdminUsername)

	dbAdminUserPassword := d.Get(TF_FIELD_CLUSTER_ADMIN_PW).(string)
	jobData.SetAdminPassword(dbAdminUserPassword)

	var iPort int
	port := d.Get(TF_FIELD_CLUSTER_PORT).(string)
	if err = CheckForEmptyAndSetDefault(&port, gDefultDbPortMap, clusterType); err != nil {
		return err
	}
	if iPort, err = strconv.Atoi(port); err != nil {
		slog.Error(funcName, "ERROR", "Non-numeric database port")
		return err
	}
	jobData.SetPort(int32(iPort))

	disableFirewall := d.Get(TF_FIELD_CLUSTER_DISABLE_FW).(bool)
	jobData.SetDisableFirewall(disableFirewall)

	installSoftware := d.Get(TF_FIELD_CLUSTER_INSTALL_SW).(bool)
	jobData.SetInstallSoftware(installSoftware)

	jobData.SetDisableSelinux(true)
	jobData.SetGenerateToken(true)

	sshUser := d.Get(TF_FIELD_CLUSTER_SSH_USER).(string)
	jobData.SetSshUser(sshUser)

	// TODO: need to provide support for it in api definition (yaml)
	//sshUserPassword := d.Get(TF_FIELD_CLUSTER_SSH_PW).(string)
	//jobData.SetSshUserPassword(sshUserPassword)

	sshKeyFile := d.Get(TF_FIELD_CLUSTER_SSH_KEY_FILE).(string)
	jobData.SetSshKeyfile(sshKeyFile)

	sshPort := d.Get(TF_FIELD_CLUSTER_SSH_PORT).(string)
	jobData.SetSshPort(sshPort)

	tfTags := d.Get(TF_FIELD_CLUSTER_TAGS).(*schema.Set).List()
	tags := make([]string, len(tfTags))
	for i, tfTag := range tfTags {
		tags[i] = tfTag.(string)
	}
	jobData.SetWithTags(tags)

	// TODO: provide support for timeout configuration - galera deployment ....
	//timeouts := d.Get("timeouts").(types.Map)

	return nil
}

type memberHosts struct {
	vanillaNode  *openapi.JobsJobJobSpecJobDataNodesInner
	mongoCfgNode *openapi.JobsJobJobSpecJobDataConfigServersMembersInner
	mongoRsNode  *openapi.JobsJobJobSpecJobDataReplicaSetsInnerMembersInner
}

func getCommonHostAttributes(f map[string]any, iPort int, clusterType string, node memberHosts) {
	//func getCommonHostAttributes(f map[string]any, iPort int, clusterType string, node *openapi.JobsJobJobSpecJobDataNodesInner) {
	funcName := "getCommonHostAttributes"

	hostname := f[TF_FIELD_CLUSTER_HOSTNAME].(string)
	hostname_data := f[TF_FIELD_CLUSTER_HOSTNAME_DATA].(string)
	hostname_internal := f[TF_FIELD_CLUSTER_HOSTNAME_INT].(string)
	port := f[TF_FIELD_CLUSTER_HOST_PORT].(string)

	if len(hostname_data) == 0 {
		hostname_data = hostname
	}

	if port == "" {
		if iPort > 1024 {
			port = strconv.Itoa(iPort)
		} else {
			port = gDefultDbPortMap[clusterType]
		}
	}

	slog.Debug(funcName, TF_FIELD_CLUSTER_HOSTNAME, hostname,
		TF_FIELD_CLUSTER_HOSTNAME_DATA, hostname_data,
		TF_FIELD_CLUSTER_HOSTNAME_INT, hostname_internal,
		TF_FIELD_CLUSTER_HOST_PORT, port)

	if node.vanillaNode != nil {
		node.vanillaNode.SetHostname(hostname)
		node.vanillaNode.SetHostnameData(hostname_data)
		node.vanillaNode.SetHostnameInternal(hostname_internal)
		node.vanillaNode.SetPort(port)
	} else if node.mongoCfgNode != nil {
		node.mongoCfgNode.SetHostname(hostname)
		node.mongoCfgNode.SetHostnameData(hostname_data)
		node.mongoCfgNode.SetHostnameInternal(hostname_internal)
		node.mongoCfgNode.SetPort(port)
	} else if node.mongoRsNode != nil {
		node.mongoRsNode.SetHostname(hostname)
		node.mongoRsNode.SetHostnameData(hostname_data)
		node.mongoRsNode.SetHostnameInternal(hostname_internal)
		node.mongoRsNode.SetPort(port)
	} else {
		slog.Warn(funcName, "Unknown node", "")
	}
}
