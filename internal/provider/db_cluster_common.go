package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

	var err error

	clusterName := d.Get(TF_FIELD_CLUSTER_NAME).(string)
	jobData.SetClusterName(clusterName)

	extClusterType := d.Get(TF_FIELD_CLUSTER_TYPE).(string)
	clusterType, ok := gExtClusterTypeToIntClusterTypeMap[extClusterType]
	if !ok {
		return errors.New(fmt.Sprintf("Unsupported cluster-type: %s", extClusterType))
	}
	jobData.SetClusterType(clusterType)

	extVendor := d.Get(TF_FIELD_CLUSTER_VENDOR).(string)
	dbVendor, ok := gExtVendorIntVendorMap[extVendor]
	if !ok {
		return errors.New(fmt.Sprintf("Unsupported vendor: %s", extVendor))
	}
	jobData.SetVendor(dbVendor)

	dbVersion := d.Get(TF_FIELD_CLUSTER_VERSION).(string)
	jobData.SetVersion(dbVersion)

	dbAdminUsername := d.Get(TF_FIELD_CLUSTER_ADMIN_USER).(string)
	if err = CheckForEmptyAndSetDefault(&dbAdminUsername, gDefultDbAdminUser, clusterType); err != nil {
		return err
	}
	if dbAdminUsername != "" {
		jobData.SetDbUser(dbAdminUsername)
	}

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

	disableSelinux := d.Get(TF_FIELD_CLUSTER_DISABLE_SELINUX).(bool)
	jobData.SetDisableSelinux(disableSelinux)

	installSoftware := d.Get(TF_FIELD_CLUSTER_INSTALL_SW).(bool)
	jobData.SetInstallSoftware(installSoftware)

	uninstallSoftware := d.Get(TF_FIELD_CLUSTER_ENABLE_UNINSTALL).(bool)
	jobData.SetEnableUninstall(uninstallSoftware)

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

	enableSsl := d.Get(TF_FIELD_CLUSTER_SSL).(bool)
	jobData.SetEnableSsl(enableSsl)

	return nil
}

func (c *DbCommon) HandleRead(ctx context.Context, d *schema.ResourceData, m interface{}, clusterInfo *openapi.ClusterResponse) error {
	funcName := "DbCommon::HandleRead"
	slog.Info(funcName)

	//var err error

	return nil
}

// *******************************************************************************
// Method: IsUpdateBatchAllowed()
//
// NOTE - check if a given update batch is allowed or not. For, e.g. updates to
// Name and/or Tags should not be compbined with updates to other fields such as -
// cluster-auto-recovery, add-node, remove-node, etc. This needs to be kept in mind!
// The check is handled in this method.
// *******************************************************************************
func (c *DbCommon) IsUpdateBatchAllowed(d *schema.ResourceData) error {
	var err error

	updateClassA := d.HasChange(TF_FIELD_CLUSTER_NAME) || d.HasChange(TF_FIELD_CLUSTER_TAGS)
	updateClassAprime := d.HasChangesExcept(TF_FIELD_CLUSTER_NAME, TF_FIELD_CLUSTER_TAGS)
	if updateClassA && updateClassAprime {
		err = errors.New("You are not allowed to update Cluster (Name or Tags) along with any other fields." +
			"Update Name/Tags in one batch and any other allowed fileds in a separate batch.")
		return err
	}

	updateClassB := d.HasChange(TF_FIELD_CLUSTER_AUTO_RECOVERY)
	updateClassBprime := d.HasChangeExcept(TF_FIELD_CLUSTER_AUTO_RECOVERY)
	if updateClassB && updateClassBprime {
		err = errors.New("You are not allowed to update cluster auto-recovery along with any other fields." +
			"Update cluster auto-recovery in one batch and any other allowed fileds in a separate batch.")
		return err
	}

	return nil
}

// *******************************************************************************
// Method: HandleUpdate()
//
// Prerequisite: The caller of this method has already called IsUpdateBatchAllowed()
// to check whether the allowed batch of updates is allowed or now. If disallowed,
// the caller should NOT call this method.
// *******************************************************************************
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

func (c *DbCommon) findMasterNode(clusterInfo *openapi.ClusterResponse, hostClass string, masterRole string) (*openapi.ClusterResponseHostsInner, error) {
	var node *openapi.ClusterResponseHostsInner
	var err error

	isFound := false
	hosts := clusterInfo.GetHosts()
	for i := 0; i < len(hosts) && !isFound; i++ {
		node = &hosts[i]
		if strings.EqualFold(node.GetClassName(), hostClass) &&
			strings.EqualFold(node.GetRole(), masterRole) {
			isFound = true
		}
	}

	if !isFound {
		err = errors.New("Master/Primary not found in CMON")
		node = nil
	}

	return node, err
}

func (c *DbCommon) findHostEntry(hostname string, hosts interface{}) map[string]any {
	for _, ff := range hosts.([]any) {
		f := ff.(map[string]any)
		hn := f[TF_FIELD_CLUSTER_HOSTNAME].(string)
		if strings.EqualFold(hostname, hn) {
			return f
		}
	}
	return nil
}

func (c *DbCommon) getHosts(d *schema.ResourceData) ([]openapi.JobsJobJobSpecJobDataNodesInner, error) {
	var nodes []openapi.JobsJobJobSpecJobDataNodesInner

	hosts := d.Get(TF_FIELD_CLUSTER_HOST)
	for _, ff := range hosts.([]any) {
		f := ff.(map[string]any)
		// Capturing hostname only as this is only used for comparison purposes.
		hostname := f[TF_FIELD_CLUSTER_HOSTNAME].(string)
		if hostname == "" {
			continue
		}
		var node = openapi.JobsJobJobSpecJobDataNodesInner{
			Hostname: &hostname,
		}
		nodes = append(nodes, node)
	}

	return nodes, nil
}

func (c *DbCommon) determineNodesDelta(nodes []openapi.JobsJobJobSpecJobDataNodesInner, clusterInfo *openapi.ClusterResponse, hostClass string) ([]openapi.JobsJobJobSpecJobDataNodesInner, []openapi.JobsJobJobSpecJobDataNodesInner, error) {
	funcName := "DbCommon::determineNodesDelta"
	slog.Info(funcName)

	var nodesToAdd []openapi.JobsJobJobSpecJobDataNodesInner
	var nodesToRemove []openapi.JobsJobJobSpecJobDataNodesInner

	// Locate the node that is in TF but not in CMON; That node needs to be added
	for i := 0; i < len(nodes); i++ {
		node := nodes[i]
		isFound := false
		hh := clusterInfo.GetHosts()
		for j := 0; j < len(hh) && !isFound; j++ {
			if !strings.EqualFold(hh[j].GetClassName(), hostClass) {
				continue
			}
			if strings.EqualFold(node.GetHostname(), hh[j].GetHostname()) {
				slog.Info(funcName, "Found node from TF in CMON", node.GetHostname())
				isFound = true
			}
		}
		if !isFound {
			slog.Info(funcName, "Node not in CMON. Adding to CMON add-node list", node.GetHostname())
			// Need to add this node to the cluster
			nodesToAdd = append(nodesToAdd, node)
		}
	}

	// Locate the node that is in CMON but not in TF; That node needs to be removed
	hh := clusterInfo.GetHosts()
	for i := 0; i < len(hh); i++ {
		h := hh[i]
		if !strings.EqualFold(h.GetClassName(), hostClass) {
			continue
		}
		isFound := false
		for j := 0; j < len(nodes) && !isFound; j++ {
			if strings.EqualFold(nodes[j].GetHostname(), h.GetHostname()) {
				slog.Info(funcName, "Found node from CMON in TF", h.GetHostname())
				isFound = true
			}
		}
		if !isFound {
			slog.Info(funcName, "Node not in TF. Adding to CMON remove-node list", h.GetHostname())
			// Need to remove this node from the cluster
			var n = openapi.JobsJobJobSpecJobDataNodesInner{}
			n.SetHostname(h.GetHostname())
			n.SetHostnameInternal(h.GetHostnameInternal())
			n.SetHostnameData(h.GetHostnameData())
			nodesToRemove = append(nodesToRemove, n)
		}
	}

	return nodesToAdd, nodesToRemove, nil
}
