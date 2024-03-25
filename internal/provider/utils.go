package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/severalnines/clustercontrol-client-sdk/go/pkg/openapi"
	"io"
	"log/slog"
	"net/http"
	"strconv"
)

func GetClusterIdByClusterName(ctx context.Context, apiClient *openapi.APIClient, clusterName string) (int32, error) {
	funcName := "GetClusterIdByName"
	slog.Debug(funcName)

	var err error
	var resp *http.Response
	var clusterId int32 = -1

	/*
	 * Get `cluster_id` to return back to Terraform
	 */
	clusterInfoReq := *openapi.NewClusters(CMON_CLUSTERS_OPERATION_GET_CLUSTERS)
	clusterInfoReq.SetClusterName(clusterName)
	if resp, err = apiClient.ClustersAPI.ClustersPost(ctx).Clusters(clusterInfoReq).Execute(); err != nil {
		PrintError(err, resp)
		return clusterId, err
	}
	slog.Debug(funcName, "Resp `ClustersPost.getclusterinfo`", resp)

	var respBytes []byte
	if respBytes, err = io.ReadAll(resp.Body); err != nil {
		PrintError(err, nil)
		return clusterId, err
	}

	var clusterInfoResp ClusterResponseFields
	if err = json.Unmarshal(respBytes, &clusterInfoResp); err != nil {
		PrintError(err, nil)
		return clusterId, err
	}
	slog.Debug(funcName, "Resp `Job`", clusterInfoResp)

	clusterId = clusterInfoResp.Cluster.GetClusterId()
	return clusterId, err
}

func GetClusterByClusterStrId(ctx context.Context, apiClient *openapi.APIClient, clusterIdStr string) (*openapi.ClusterResponse, error) {
	funcName := "GetClusterByClusterStrId"
	slog.Debug(funcName)

	var clusterId int
	var err error

	if clusterId, err = strconv.Atoi(clusterIdStr); err != nil {
		return nil, err
	}

	return (GetClusterByClusterIntId(ctx, apiClient, int32(clusterId)))
}

func GetClusterByClusterIntId(ctx context.Context, apiClient *openapi.APIClient, clusterId int32) (*openapi.ClusterResponse, error) {
	funcName := "GetClusterByClusterIntId"
	slog.Debug(funcName)

	var err error
	var resp *http.Response
	var clusterInfoResp ClusterResponseFields

	/*
	 * Get `cluster_id` to return back to Terraform
	 */
	clusterInfoReq := *openapi.NewClusters(CMON_CLUSTERS_OPERATION_GET_CLUSTERS)
	clusterInfoReq.SetClusterId(clusterId)
	clusterInfoReq.SetWithHosts(true)
	clusterInfoReq.SetOffset(0)
	clusterInfoReq.SetLimit(100)
	if resp, err = apiClient.ClustersAPI.ClustersPost(ctx).Clusters(clusterInfoReq).Execute(); err != nil {
		PrintError(err, resp)
		return &clusterInfoResp.Cluster, err
	}
	slog.Debug(funcName, "Resp `ClustersPost.getclusterinfo`", resp, "clusterId", clusterId)

	var respBytes []byte
	if respBytes, err = io.ReadAll(resp.Body); err != nil {
		PrintError(err, nil)
		return &clusterInfoResp.Cluster, err
	}

	if err = json.Unmarshal(respBytes, &clusterInfoResp); err != nil {
		PrintError(err, nil)
		return &clusterInfoResp.Cluster, err
	}
	slog.Debug(funcName, "Resp `Job`", clusterInfoResp)

	//clusterId = clusterInfoResp.Cluster.Cluster_Id
	return &clusterInfoResp.Cluster, nil
}

func GetClusterIdFromSchema(d *schema.ResourceData) (int32, diag.Diagnostics) {
	funcName := "GetClusterId"
	slog.Debug(funcName)

	var diags diag.Diagnostics
	var err error
	var iCid int

	clusterId := d.Get(TF_FIELD_CLUSTER_ID).(string)
	if clusterId == "" {
		strErr := fmt.Sprintf("%s: %s - not declared", funcName, TF_FIELD_CLUSTER_ID)
		slog.Error(strErr)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  strErr,
		})
		return int32(iCid), diags
	}

	if iCid, err = strconv.Atoi(clusterId); err != nil {
		strErr := fmt.Sprintf("%s: %s - non-numeric cluster-id", funcName, clusterId)
		slog.Error(strErr)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  strErr,
		})
		return int32(iCid), diags
	}

	if iCid == 0 {
		strErr := fmt.Sprintf("%s: - invalid cluster-id 0", funcName)
		slog.Error(strErr)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  strErr,
		})
		return int32(iCid), diags
	}

	return int32(iCid), nil
}

func GetBackupIdForCluster(ctx context.Context, apiClient *openapi.APIClient, clusterId int32, jobId int32) (int32, error) {
	funcName := "GetBackupIdForCluster"
	slog.Debug(funcName)

	var err error
	var backupId int32
	var resp *http.Response
	//var backupResp BackupGetOperationResponse

	slog.Info(funcName, "Job Id", jobId)

	getBackupsOp := *openapi.NewBackup(CMON_BACKUP_OPERATION_GET)
	getBackupsOp.SetClusterId(clusterId)
	getBackupsOp.SetBackupRecordVersion(BACKUP_RECORD_VERSION_2)
	getBackupsOp.SetOrder(BACKUP_ORDER_CREATED_DESC)
	getBackupsOp.SetAscending(false)
	getBackupsOp.SetLimit(500)
	getBackupsOp.SetOffset(0)
	if resp, err = apiClient.BackupAPI.BackupPost(ctx).Backup(getBackupsOp).Execute(); err != nil {
		PrintError(err, resp)
		return backupId, err
	}
	slog.Debug(funcName, "Resp `BackupPost.getBackups`", resp, "clusterId", clusterId)

	var respBytes []byte
	if respBytes, err = io.ReadAll(resp.Body); err != nil {
		PrintError(err, nil)
		return backupId, err
	}

	var topLevelResp BackupGetOperationResponseTop
	if err = json.Unmarshal(respBytes, &topLevelResp); err != nil {
		PrintError(err, nil)
		return backupId, err
	}
	slog.Info(funcName, "Resp `GetBackups` req-status", topLevelResp.Request_Status, "Total recs", topLevelResp.Total)
	//slog.Info(funcName, "Resp `GetBackups` total recs", topLevelResp)

	// Find the backup-id that matches the job-id
	// BackupResponseBackupRecordsInner
	if topLevelResp.Total == 0 {
		err = errors.New("ERROR: Zero backup records.")
		return backupId, err
	}

	var fullResp BackupGetOperationResponse
	fullResp.Backup_Records = make([]openapi.BackupResponseBackupRecordsInner, topLevelResp.Total)
	//var meta openapi.BackupResponseBackupRecordsInnerMetadata
	//for _, br := range fullResp.Backup_Records {
	//	br.SetMetadata(meta)
	//}
	if err = json.Unmarshal(respBytes, &fullResp); err != nil {
		PrintError(err, nil)
		return backupId, err
	}
	for _, backupRec := range fullResp.Backup_Records {
		metaData := backupRec.GetMetadata()
		if metaData.GetJobId() == jobId {
			backupId = metaData.GetId()
			break
		}
	}

	if backupId == 0 {
		err = errors.New("Invalid backup ID: 0")
	}

	return backupId, err
}

func convertPortToInt(strPort string, useAsDefaultPort int32) int32 {
	var iP int
	var err error
	if iP, err = strconv.Atoi(strPort); err != nil {
		return useAsDefaultPort
	}
	return int32(iP)
}

//func ConvertTimeToZulu(in string, tmFmt string) (time.Time, error) {
//	funcName := "GetClusterId"
//	slog.Debug(funcName)
//
//	var err error
//	var convertedTime time.Time
//
//	if convertedTime, err = time.Parse(tmFmt, in); err != nil {
//		return convertedTime, err
//	}
//
//	return convertedTime, nil
//}
