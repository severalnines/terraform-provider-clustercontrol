package provider

import (
	"context"
	"encoding/json"
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
		// TODO
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
