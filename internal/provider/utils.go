package provider

import (
	"context"
	"encoding/json"
	"github.com/severalnines/clustercontrol-client-sdk/go/pkg/openapi"
	"io"
	"log/slog"
	"net/http"
)

func GetClusterIdByName(ctx context.Context, apiClient *openapi.APIClient, clusterName string) (int32, error) {
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
	if resp, err = apiClient.ClustersApi.ClustersPost(ctx).Clusters(clusterInfoReq).Execute(); err != nil {
		PrintError(err, resp)
		return clusterId, err
	}
	slog.Debug(funcName, "Resp `ClustersPost.getallclusterinfo`", resp)

	var respBytes []byte
	if respBytes, err = io.ReadAll(resp.Body); err != nil {
		PrintError(err, nil)
		return clusterId, err
	}

	var clusterInfoResp ClusterInfoRespJson
	if err = json.Unmarshal(respBytes, &clusterInfoResp); err != nil {
		PrintError(err, nil)
		return clusterId, err
	}
	slog.Debug(funcName, "Resp `Job`", clusterInfoResp)

	clusterId = clusterInfoResp.Cluster.Cluster_Id
	return clusterId, err
}
