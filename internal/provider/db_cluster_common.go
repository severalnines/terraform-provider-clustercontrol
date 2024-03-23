package provider

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/severalnines/clustercontrol-client-sdk/go/pkg/openapi"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
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

	sshUserPassword := d.Get(TF_FIELD_CLUSTER_SSH_PW).(string)
	jobData.SetSudoPassword(sshUserPassword)

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

	deployAgents := d.Get(TF_FIELD_CLUSTER_DEPLOY_AGENTS).(bool)
	jobData.SetDeployAgents(deployAgents)

	// TODO: provide support for timeout configuration - galera deployment ....
	//timeouts := d.Get("timeouts").(types.Map)

	return nil
}

func (c *DbCommon) HandleRead(ctx context.Context, d *schema.ResourceData, m interface{}, clusterInfo *openapi.ClusterResponse) error {
	funcName := "DbCommon::HandleRead"
	slog.Info(funcName)

	//var err error

	return nil
}

func (c *DbCommon) HandleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}, clusterInfo *openapi.ClusterResponse) error {
	funcName := "DbCommon::HandleUpdate"
	slog.Info(funcName)

	var configChanges []openapi.ClustersConfigurationInner

	isConfigUpdated := false
	if d.HasChange(TF_FIELD_CLUSTER_NAME) {
		var clusterNameTf string
		clusterNameTf = d.Get(TF_FIELD_CLUSTER_NAME).(string)
		var tags = openapi.ClustersConfigurationInner{
			Name:  &CMON_CLUSTERS_OPERATION_SET_NAME,
			Value: &clusterNameTf,
		}
		configChanges = append(configChanges, tags)
		isConfigUpdated = true
	}

	if d.HasChange(TF_FIELD_CLUSTER_TAGS) {
		tfTags := d.Get(TF_FIELD_CLUSTER_TAGS).(*schema.Set).List()
		tags := make([]string, len(tfTags))
		for i, tfTag := range tfTags {
			tags[i] = tfTag.(string)
		}
		newTags := strings.Join(tags[:], ";")
		var cfgTags = openapi.ClustersConfigurationInner{
			Name:  &CMON_CLUSTERS_OPERATION_SET_CLUSTER_TAG,
			Value: &newTags,
		}
		configChanges = append(configChanges, cfgTags)
		isConfigUpdated = true
	}

	var err error
	var resp *http.Response
	var clusterSetResp MinResponseFields

	if isConfigUpdated {
		apiClient := m.(*openapi.APIClient)

		clusterInfoReq := *openapi.NewClusters(CMON_CLUSTERS_OPERATION_SET_CONFIG)
		clusterInfoReq.SetClusterId(clusterInfo.GetClusterId())

		// Finally set the config changes
		clusterInfoReq.SetConfiguration(configChanges)

		if resp, err = apiClient.ClustersAPI.ClustersPost(ctx).Clusters(clusterInfoReq).Execute(); err != nil {
			PrintError(err, resp)
			return err
		}
		slog.Info(funcName, "Resp `ClustersPost.setConfig`", resp, "clusterId", clusterInfo.GetClusterId())

		var respBytes []byte
		if respBytes, err = io.ReadAll(resp.Body); err != nil {
			PrintError(err, nil)
			return err
		}

		if err = json.Unmarshal(respBytes, &clusterSetResp); err != nil {
			PrintError(err, nil)
			return err
		}
		slog.Debug(funcName, "Resp `setConfig`", clusterSetResp)

	}

	if d.HasChange(TF_FIELD_CLUSTER_AUTO_RECOVERY) {
		apiClient := m.(*openapi.APIClient)

		toggleAutoRecovery := NewCCJob(CMON_JOB_CREATE_JOB)
		toggleAutoRecovery.SetClusterId(clusterInfo.GetClusterId())
		job := toggleAutoRecovery.GetJob()
		jobSpec := job.GetJobSpec()
		jobData := jobSpec.GetJobData()

		isEnableClusterAutoRecovery := d.Get(TF_FIELD_CLUSTER_AUTO_RECOVERY).(bool)
		if isEnableClusterAutoRecovery {
			jobSpec.SetCommand(CMON_JOB_ENABLE_CLUSTER_RECOVERY_COMMAND)
		} else {
			jobSpec.SetCommand(CMON_JOB_DISABLE_CLUSTER_RECOVERY_COMMAND)
		}

		jobSpec.SetJobData(jobData)
		job.SetJobSpec(jobSpec)
		toggleAutoRecovery.SetJob(job)

		if err = SendAndWaitForJobCompletion(ctx, apiClient, toggleAutoRecovery); err != nil {
			slog.Error(err.Error())
		}
	}

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
