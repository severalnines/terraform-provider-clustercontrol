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
//   - Host class_name (CMON_CLASS_CLICKHOUSE_HOST) and the per-node "nodetype"
//     value for dedicated Keeper hosts (CLICKHOUSE_NODETYPE_KEEPER) are BOTH
//     CONFIRMED against real CMON job_data (not guesses). There is only ONE
//     host class for ClickHouse - every node uses it, regardless of role.
//   - job_data cannot distinguish a keeper-less replica from a replica that
//     also happens to run embedded Keeper - both produce an identical node
//     payload (no "nodetype" field). So only two roles are exposed here:
//     "replica" (nodetype omitted - CMON manages embedded Keeper placement
//     among replica hosts on its own) and "keeper" (dedicated Keeper-only
//     host: nodetype = "clickhouse_keeper", and its "port" is the Keeper
//     client port rather than the ClickHouse native port).
//   - Sharded ClickHouse is on ClusterControl's roadmap, not shipped yet.
//     The per-host "shard" attribute (TF_FIELD_CLUSTER_HOST_SHARD) is
//     accepted and validated here, but deliberately NOT sent to CMON -
//     there is nothing on the backend to receive it yet.
// *******************************************************************************

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

type ClickHouse struct {
	Common DbCommon
	Backup DbBackupCommon
}

// resolveClickHouseRole normalizes and validates a db_host's "roles" value
// for ClickHouse, defaulting to CLICKHOUSE_ROLE_REPLICA when unset (mirrors
// Elastic's default-when-empty pattern).
func resolveClickHouseRole(role string) (string, error) {
	if role == "" {
		return CLICKHOUSE_ROLE_REPLICA, nil
	}
	switch strings.ToLower(role) {
	case CLICKHOUSE_ROLE_REPLICA, CLICKHOUSE_ROLE_KEEPER:
		return strings.ToLower(role), nil
	default:
		return "", fmt.Errorf(
			"invalid ClickHouse host role %q - must be one of %q, %q (or unset, defaults to %q)",
			role, CLICKHOUSE_ROLE_REPLICA, CLICKHOUSE_ROLE_KEEPER, CLICKHOUSE_ROLE_REPLICA)
	}
}

// applyClickHouseNodeAttrs sets the class_name/nodetype/port on a node per
// its resolved role, and defaults hostname_data to hostname when not
// explicitly given (matches the shape of real CMON job_data, which always
// carries a populated hostname_data).
func applyClickHouseNodeAttrs(node *openapi.JobsJobJobSpecJobDataNodesInner, role string, nativePort string, keeperPort string) {
	node.SetClassName(CMON_CLASS_CLICKHOUSE_HOST)
	if role == CLICKHOUSE_ROLE_KEEPER {
		node.SetNodetype(CLICKHOUSE_NODETYPE_KEEPER)
		node.SetPort(keeperPort)
	} else {
		node.SetPort(nativePort)
	}
	if node.GetHostnameData() == "" {
		node.SetHostnameData(node.GetHostname())
	}
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
	nativePort := d.Get(TF_FIELD_CLUSTER_CLICKHOUSE_NATIVE_PORT).(string)
	if err = CheckForEmptyAndSetDefault(&nativePort, gDefultDbPortMap, clusterType); err != nil {
		return err
	}
	if iPort, err = strconv.Atoi(nativePort); err != nil {
		slog.Error(funcName, "ERROR", "Non-numeric database port")
		return err
	}
	jobData.SetPort(int32(iPort))

	// Keeper client port - used per-node for dedicated "keeper" hosts (see
	// applyClickHouseNodeAttrs). Confirmed via real CMON job_data that this
	// is a per-node Port value, not a separate cluster-wide job_data field.
	keeperPort := d.Get(TF_FIELD_CLUSTER_CLICKHOUSE_KEEPER_PORT).(string)
	if keeperPort == "" {
		keeperPort = DEFAULT_CLICKHOUSE_KEEPER_PORT
	}

	hosts := d.Get(TF_FIELD_CLUSTER_HOST)
	nodes := []openapi.JobsJobJobSpecJobDataNodesInner{}
	for _, ff := range hosts.([]any) {
		f := ff.(map[string]any)
		hostname := f[TF_FIELD_CLUSTER_HOSTNAME].(string)
		hostname_data := f[TF_FIELD_CLUSTER_HOSTNAME_DATA].(string)
		hostname_internal := f[TF_FIELD_CLUSTER_HOSTNAME_INT].(string)
		rawRole := f[TF_FIELD_CLUSTER_HOST_ROLES].(string)
		shard := f[TF_FIELD_CLUSTER_HOST_SHARD].(string)

		if hostname == "" {
			return errors.New("Hostname cannot be empty")
		}

		role, err := resolveClickHouseRole(rawRole)
		if err != nil {
			return err
		}

		// Shard is accepted/validated but not yet actionable (see
		// TF_FIELD_CLUSTER_HOST_SHARD) - sharded ClickHouse isn't shipped
		// by ClusterControl yet. Dedicated Keeper-only hosts never belong
		// to a shard, so catch that combination as a config error now
		// rather than silently ignoring it.
		if role == CLICKHOUSE_ROLE_KEEPER && shard != "" {
			return fmt.Errorf("host %q: shard cannot be set on a dedicated \"keeper\" host", hostname)
		}
		// TODO: once ClusterControl supports sharded ClickHouse, thread
		// `shard` through into job_data here (field TBD - needs SDK support).
		_ = shard

		var node = openapi.JobsJobJobSpecJobDataNodesInner{
			Hostname: &hostname,
		}
		if hostname_data != "" {
			node.SetHostnameData(hostname_data)
		}
		if hostname_internal != "" {
			node.SetHostnameInternal(hostname_internal)
		}
		applyClickHouseNodeAttrs(&node, role, nativePort, keeperPort)
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

		// Only one CMON host class exists for ClickHouse (confirmed), so a
		// single delta pass covers all hosts regardless of role.
		hostClassName := CMON_CLASS_CLICKHOUSE_HOST
		command := CMON_JOB_ADD_NODE_COMMAND

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

			node.SetHostname(nodeFromTf.GetHostname())

			// NOTE: host is guaranteed to be non-nil.
			hostTfRec := c.Common.findHostEntry(nodeFromTf.GetHostname(), d.Get(TF_FIELD_CLUSTER_HOST))
			hostname_data := hostTfRec[TF_FIELD_CLUSTER_HOSTNAME_DATA].(string)
			hostname_internal := hostTfRec[TF_FIELD_CLUSTER_HOSTNAME_INT].(string)
			rawRole := hostTfRec[TF_FIELD_CLUSTER_HOST_ROLES].(string)

			role, err := resolveClickHouseRole(rawRole)
			if err != nil {
				return err
			}

			if hostname_data != "" {
				node.SetHostnameData(hostname_data)
			}
			if hostname_internal != "" {
				node.SetHostnameInternal(hostname_internal)
			}

			nativePort := strconv.Itoa(int(tmpJobData.GetPort()))
			keeperPort := d.Get(TF_FIELD_CLUSTER_CLICKHOUSE_KEEPER_PORT).(string)
			if keeperPort == "" {
				keeperPort = DEFAULT_CLICKHOUSE_KEEPER_PORT
			}
			applyClickHouseNodeAttrs(&node, role, nativePort, keeperPort)

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
