package provider

import (
	"errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/severalnines/clustercontrol-client-sdk/go/pkg/openapi"
	"log/slog"
	"strings"
)

type LBCommon struct {
}

func (m *LBCommon) GetInputs(d map[string]any, jobData *openapi.JobsJobJobSpecJobData) error {
	funcName := "ProxySql::GetInputs"
	slog.Debug(funcName)

	disableFirewall := d[TF_FIELD_CLUSTER_DISABLE_FW].(bool)
	jobData.SetDisableFirewall(disableFirewall)

	disableSelinux := d[TF_FIELD_CLUSTER_DISABLE_SELINUX].(bool)
	jobData.SetDisableSelinux(disableSelinux)

	installSoftware := d[TF_FIELD_LB_INSTALL_SW].(bool)
	jobData.SetInstallSoftware(installSoftware)

	uninstallSoftware := d[TF_FIELD_LB_ENABLE_UNINSTALL].(bool)
	jobData.SetEnableUninstall(uninstallSoftware)

	sshUser := d[TF_FIELD_CLUSTER_SSH_USER].(string)
	jobData.SetSshUser(sshUser)

	sshUserPassword := d[TF_FIELD_CLUSTER_SSH_PW].(string)
	jobData.SetSudoPassword(sshUserPassword)

	sshKeyFile := d[TF_FIELD_CLUSTER_SSH_KEY_FILE].(string)
	jobData.SetSshKeyfile(sshKeyFile)

	sshPort := d[TF_FIELD_CLUSTER_SSH_PORT].(string)
	jobData.SetSshPort(sshPort)

	return nil
}

func (c *LBCommon) getHosts(d *schema.ResourceData) ([]openapi.JobsJobJobSpecJobDataNodesInner, error) {
	var nodes []openapi.JobsJobJobSpecJobDataNodesInner

	hosts := d.Get(TF_FIELD_CLUSTER_LOAD_BALANCER)
	for _, ff := range hosts.([]any) {
		f := ff.(map[string]any)
		// Capturing hostname only as this is only used for comparison purposes.
		myHost := f[TF_FIELD_LB_MY_HOST]
		for _, tt := range myHost.([]any) {
			t := tt.(map[string]any)
			hostname := t[TF_FIELD_CLUSTER_HOSTNAME].(string)
			if hostname == "" {
				continue
			}
			var node = openapi.JobsJobJobSpecJobDataNodesInner{
				Hostname: &hostname,
			}
			nodes = append(nodes, node)

			// Should only have one db_my_host entry
			break
		}

		// Suppose we can allow for multiple LB entries
	}

	return nodes, nil
}

func (c *LBCommon) findLoadbalancerEntry(d *schema.ResourceData, hostname string) (map[string]any, error) {
	var err error
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
			hn := t[TF_FIELD_CLUSTER_HOSTNAME].(string)
			if strings.EqualFold(hn, hostname) {
				// found it
				theTfRecord = f
				isFound = true
			}
		}
	}
	return theTfRecord, err
}

func (c *LBCommon) determineProxyDelta(d *schema.ResourceData, clusterInfo *openapi.ClusterResponse, hostClass string) ([]openapi.JobsJobJobSpecJobDataNodesInner, []openapi.JobsJobJobSpecJobDataNodesInner, error) {
	funcName := "determineProxyDelta::determineNodesDelta"
	slog.Info(funcName)

	var nodesToAdd []openapi.JobsJobJobSpecJobDataNodesInner
	var nodesToRemove []openapi.JobsJobJobSpecJobDataNodesInner
	hosts := d.Get(TF_FIELD_CLUSTER_LOAD_BALANCER)
	nodes := []openapi.JobsJobJobSpecJobDataNodesInner{}
	for _, ff := range hosts.([]any) {
		f := ff.(map[string]any)
		// Capturing hostname only as this is only used for comparison purposes.
		myHost := f[TF_FIELD_LB_MY_HOST]
		for _, tt := range myHost.([]any) {
			t := tt.(map[string]any)
			hostname := t[TF_FIELD_CLUSTER_HOSTNAME].(string)
			if hostname == "" {
				return nil, nil, errors.New("Hostname cannot be empty")
			}
			var node = openapi.JobsJobJobSpecJobDataNodesInner{
				Hostname: &hostname,
			}
			nodes = append(nodes, node)
		}
	}

	// Locate the node that is in TF but not in CMON; That node needs to be added
	for i := 0; i < len(nodes); i++ {
		node := nodes[i]
		isFound := false
		hs := clusterInfo.GetHosts()
		for j := 0; j < len(hs); j++ {
			if !strings.EqualFold(hs[j].GetClassName(), hostClass) {
				continue
			}
			if strings.EqualFold(node.GetHostname(), hs[j].GetHostname()) {
				slog.Info(funcName, "Found node from TF in CMON", node.GetHostname())
				isFound = true
				break
			}
		}
		if !isFound {
			slog.Info(funcName, "Node not in CMON. Adding to CMON add-node list", node.GetHostname())
			// Need to add this node to the cluster
			nodesToAdd = append(nodesToAdd, node)
		}
	}

	// Locate the node that is in CMON but not in TF; That node needs to be removed
	h := clusterInfo.GetHosts()
	for i := 0; i < len(h); i++ {
		host := h[i]
		if !strings.EqualFold(host.GetClassName(), hostClass) {
			continue
		}
		isFound := false
		for j := 0; j < len(nodes); j++ {
			if strings.EqualFold(nodes[j].GetHostname(), host.GetHostname()) {
				slog.Info(funcName, "Found node from CMON in TF", host.GetHostname())
				isFound = true
				break
			}
		}
		if !isFound {
			slog.Info(funcName, "Node not in TF. Adding to CMON remove-node list", host.GetHostname())
			// Need to remove this node from the cluster
			var n = openapi.JobsJobJobSpecJobDataNodesInner{}
			n.SetHostname(host.GetHostname())
			n.SetHostnameInternal(host.GetHostnameInternal())
			n.SetHostnameData(host.GetHostnameData())
			nodesToRemove = append(nodesToRemove, n)
		}
	}

	return nodesToAdd, nodesToRemove, nil
}
