package provider

// *******************************************************************************
// ClickHouse engine support.
//
// STATUS / CAVEATS:
//   - ClickHouse is entirely absent from the current clustercontrol-client-sdk
//     OpenAPI spec (no cluster_type, vendor, or host class_name entries exist
//     for it). This file is written against the SDK version available at the
//     time of writing and will not compile until the SDK is regenerated with
//     ClickHouse support added to clustercontrol-v2.yaml.
//   - CMON_CLASS_CLICKHOUSE_HOST (constants.go) is a proposed value, not
//     confirmed against the real CMON backend.
//   - There is no dedicated "keeper port" field on JobsJobJobSpecJobData today.
//     The call to jobData.SetKeeperPort(...) below is commented out and marked
//     TODO; it needs a corresponding job_data field added to the OpenAPI spec
//     (see clustercontrol-client-sdk/clustercontrol-v2.yaml) before it can be
//     wired in for real. Everything else in this file uses fields that already
//     exist on JobsJobJobSpecJobData today (Port, Nodes, ClassName, Hostname,
//     etc.) and should compile as-is.
//   - Topology: ClusterControl 2.5 supports standalone (1 node) or replicated
//     (3+ nodes) ClickHouse with embedded Keeper - no separate Keeper host
//     tier, no sharding yet. This mirrors Elastic's single-host-class model
//     (see elastic.go) rather than Mongo's tiered config-server/shard model:
//     every ClickHouse node is a symmetric peer, so no "roles" attribute is
//     needed on db_host the way Elasticsearch needs master/data roles.
// *******************************************************************************

import (
	"context"
	"errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/severalnines/clustercontrol-client-sdk/go/pkg/openapi"
	"log/slog"
	"strconv"
	"strings"
)

type ClickHouse struct {
	Common DbCommon
	Backup DbBackupCommon
}

func (c *ClickHouse) GetInputs(d *schema.ResourceData, jobData *openapi.JobsJobJobSpecJobData) error {
	funcName := "ClickHouse::GetInputs"
	slog.Debug(funcName)

	var err error

	// parent/super - get common attributes (also sets db_enable_ssl from
	// TF_FIELD_CLUSTER_SSL; the ClickHouse examples hardcode that variable
	// to true in main.tf since SSL is mandatory for this engine)
	if err = c.Common.GetInputs(d, jobData); err != nil {
		return err
	}

	clusterType := jobData.GetClusterType()

	// Native TCP port (secure/TLS variant by default - see DEFAULT_CLICKHOUSE_NATIVE_PORT)
	var iPort int
	port := d.Get(TF_FIELD_CLUSTER_CLICKHOUSE_NATIVE_PORT).(string)
	if err = CheckForEmptyAndSetDefault(&port, gDefultDbPortMap, clusterType); err != nil {
		return err
	}
	if iPort, err = strconv.Atoi(port); err != nil {
		slog.Error(funcName, "ERROR", "Non-numeric database port")
		return err
	}
	jobData.SetPort(int32(iPort))

	// Embedded ClickHouse Keeper client port (replicated topology only;
	// harmless to send for standalone too, CMON should ignore it there).
	keeperPort := d.Get(TF_FIELD_CLUSTER_CLICKHOUSE_KEEPER_PORT).(string)
	if keeperPort == "" {
		keeperPort = DEFAULT_CLICKHOUSE_KEEPER_PORT
	}
	// TODO: jobData.SetKeeperPort(keeperPort) - field does not exist yet on
	// JobsJobJobSpecJobData. Needs to be added to clustercontrol-v2.yaml and
	// the SDK regenerated (see clustercontrol-client-sdk/generate_go.sh)
	// before this can be sent to CMON.
	_ = keeperPort

	hosts := d.Get(TF_FIELD_CLUSTER_HOST)
	nodes := []openapi.JobsJobJobSpecJobDataNodesInner{}
	for _, ff := range hosts.([]any) {
		f := ff.(map[string]any)
		hostname := f[TF_FIELD_CLUSTER_HOSTNAME].(string)
		hostname_data := f[TF_FIELD_CLUSTER_HOSTNAME_DATA].(string)
		hostname_internal := f[TF_FIELD_CLUSTER_HOSTNAME_INT].(string)

		if hostname == "" {
			return errors.New("Hostname cannot be empty")
		}
		var node = openapi.JobsJobJobSpecJobDataNodesInner{
			Hostname: &hostname,
		}
		node.SetClassName(CMON_CLASS_CLICKHOUSE_HOST)
		if hostname_data != "" {
			node.SetHostnameData(hostname_data)
		}
		if hostname_internal != "" {
			node.SetHostnameInternal(hostname_internal)
		}
		// No roles/protocol needed - ClickHouse nodes are symmetric peers,
		// unlike Elasticsearch's master/data roles.
		nodes = append(nodes, node)
	}
	jobData.SetNodes(nodes)

	return nil
}

func (c *ClickHouse) HandleRead(ctx context.Context, d *schema.ResourceData, apiClient *openapi.APIClient, clusterInfo *openapi.ClusterResponse) error {

	if err := c.Common.HandleRead(ctx, d, apiClient, clusterInfo); err != nil {
		return err
	}

	return nil
}

func (c *ClickHouse) IsUpdateBatchAllowed(d *schema.ResourceData) error {
	var err error

	if err = c.Common.IsUpdateBatchAllowed(d); err != nil {
		return err
	}

	return nil
}

func (c *ClickHouse) HandleUpdate(ctx context.Context, d *schema.ResourceData, apiClient *openapi.APIClient, clusterInfo *openapi.ClusterResponse) error {
	funcName := "ClickHouse::HandleUpdate"
	slog.Debug(funcName)

	var err error

	if err := c.Common.HandleUpdate(ctx, d, apiClient, clusterInfo); err != nil {
		return err
	}

	tmpJobData := openapi.NewJobsJobJobSpecJobData()
	if err = c.GetInputs(d, tmpJobData); err != nil {
		return err
	}

	if d.HasChange(TF_FIELD_CLUSTER_HOST) {
		var nodesToAdd []openapi.JobsJobJobSpecJobDataNodesInner
		var nodesToRemove []openapi.JobsJobJobSpecJobDataNodesInner

		hostClassName := CMON_CLASS_CLICKHOUSE_HOST
		command := CMON_JOB_ADD_NODE_COMMAND

		// Compare Terraform and CMON to determine whether adding or removing a node
		nodes, _ := c.Common.getHosts(d)
		if nodesToAdd, nodesToRemove, err = c.Common.determineNodesDelta(nodes, clusterInfo, hostClassName); err != nil {
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
			return errors.New("Unsupported ClickHouse operation. Neither Add nor Remove node.")
		}

		// From Terraform
		tmpJobDataNodes := tmpJobData.GetNodes()
		var nodeFromTf *openapi.JobsJobJobSpecJobDataNodesInner
		for i := 0; i < len(tmpJobDataNodes) && nodeToAddOrRemove != nil; i++ {
			tmpJobDataNode := tmpJobDataNodes[i]
			if strings.EqualFold(tmpJobDataNode.GetHostname(), nodeToAddOrRemove.GetHostname()) {
				nodeFromTf = &tmpJobDataNode
				break
			}
		}
		// No need to error check as the node must be in the list

		addOrRemoveNodeJob := NewCCJob(CMON_JOB_CREATE_JOB)
		addOrRemoveNodeJob.SetClusterId(clusterInfo.GetClusterId())
		job := addOrRemoveNodeJob.GetJob()
		jobSpec := job.GetJobSpec()
		jobData := jobSpec.GetJobData()
		jobSpec.SetCommand(command)

		// No concept of Primary/master node in ClickHouse's symmetric-peer
		// replicated topology, so no promotion step is needed here.

		jobData.SetInstallSoftware(tmpJobData.GetInstallSoftware())
		jobData.SetDisableSelinux(tmpJobData.GetDisableSelinux())
		jobData.SetDisableFirewall(tmpJobData.GetDisableFirewall())

		if isAddNode {

			var nodes []openapi.JobsJobJobSpecJobDataNodesInner
			var node openapi.JobsJobJobSpecJobDataNodesInner

			node.SetClassName(hostClassName)
			node.SetHostname(nodeFromTf.GetHostname())

			// NOTE: host is guaranteed to be non-nil.
			hostTfRec := c.Common.findHostEntry(nodeFromTf.GetHostname(), d.Get(TF_FIELD_CLUSTER_HOST))
			hostname_data := hostTfRec[TF_FIELD_CLUSTER_HOSTNAME_DATA].(string)
			hostname_internal := hostTfRec[TF_FIELD_CLUSTER_HOSTNAME_INT].(string)

			if hostname_data != "" {
				node.SetHostnameData(hostname_data)
			} else {
				node.SetHostnameData(node.GetHostname())
			}
			if hostname_internal != "" {
				node.SetHostnameInternal(hostname_internal)
			}
			node.SetPort(strconv.Itoa(int(tmpJobData.GetPort())))

			nodes = append(nodes, node)
			jobData.SetNodes(nodes)
		} else if isRemoveNode {
			var node openapi.JobsJobJobSpecJobDataNode
			node.SetHostname(nodeToAddOrRemove.GetHostname())
			node.SetPort(tmpJobData.GetPort())
			jobData.SetEnableUninstall(true)
			jobData.SetUnregisterOnly(false)
			jobData.SetNode(node)
		} else {
			return errors.New("Unsupported ClickHouse operation.")
		}

		jobSpec.SetJobData(jobData)
		job.SetJobSpec(jobSpec)
		addOrRemoveNodeJob.SetJob(job)

		if err = SendAndWaitForJobCompletion(ctx, apiClient, addOrRemoveNodeJob); err != nil {
			slog.Error(err.Error())
			return err
		}

	} // d.HasChange(TF_FIELD_CLUSTER_HOST)

	return nil
}

func (c *ClickHouse) GetBackupInputs(d *schema.ResourceData, jobData *openapi.JobsJobJobSpecJobData) error {
	funcName := "ClickHouse::GetBackupInputs"
	slog.Debug(funcName)

	var err error

	// parent/super - get common attributes (db_backup_method is expected to
	// be "clickhouse-native" or "clickhouse-native-incr" - see IsValidBackupOptions)
	if err = c.Backup.GetBackupInputs(d, jobData); err != nil {
		return err
	}

	return err
}

func (c *ClickHouse) IsValidBackupOptions(vendor string, clusterType string, jobData *openapi.JobsJobJobSpecJobData) error {
	return c.Backup.IsValidBackupOptions(vendor, clusterType, jobData)
}

func (c *ClickHouse) SetBackupJobData(jobData *openapi.JobsJobJobSpecJobData) error {
	return c.Backup.SetBackupJobData(jobData)
}

func (c *ClickHouse) IsBackupRemovable(clusterInfo *openapi.ClusterResponse, jobData *openapi.JobsJobJobSpecJobData) bool {
	return true
}

func NewClickHouse() *ClickHouse {
	return &ClickHouse{}
}
