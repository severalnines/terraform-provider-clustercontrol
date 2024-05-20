package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/severalnines/clustercontrol-client-sdk/go/pkg/openapi"
	"io"
	"log/slog"
	"net/http"
	"time"
)

func resourceDbClusterMaintenance() *schema.Resource {
	funcName := "resourceDbLoadBalancer"
	slog.Debug(funcName)

	return &schema.Resource{
		CreateContext: resourceCreateDbMaintenance,
		ReadContext:   resourceReadDbMaintenance,
		UpdateContext: resourceUpdateDbMaintenance,
		DeleteContext: resourceDeleteDbMaintenance,
		Importer:      &schema.ResourceImporter{},
		Schema: map[string]*schema.Schema{
			TF_FIELD_RESOURCE_ID: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the resource allocated by ClusterControl.",
			},
			TF_FIELD_LAST_UPDATED: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last updated timestamp for the resource in question.",
			},
			TF_FIELD_CLUSTER_ID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The database cluster ID for which this LB is being deployed to.",
			},
			TF_FIELD_MAINT_START_TIME: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Format: `Jan-02-2006T15:04`",
			},
			TF_FIELD_MAINT_STOP_TIME: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Format: `Jan-02-2006T15:04`",
			},
			TF_FIELD_MAINT_REASON: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The reason for the maintenance window. Something meaningful and short.",
			},
		},
	}
}

func resourceCreateDbMaintenance(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	funcName := "resourceCreateDbMaintenance"
	slog.Debug(funcName)

	var diags diag.Diagnostics
	var err error

	providerDetails := m.(*ProviderDetails)

	//newCtx := context.WithValue(ctx, "cookie", gNewCtx.Value("cookie"))
	newCtx := context.WithValue(ctx, "cookie", providerDetails.SessionIdCtx.Value("cookie"))

	//apiClient := m.(*openapi.APIClient)
	apiClient := providerDetails.ApiClient

	maintStartTmStr := d.Get(TF_FIELD_MAINT_START_TIME).(string)
	maintStopTmStr := d.Get(TF_FIELD_MAINT_STOP_TIME).(string)
	maintReasonStr := d.Get(TF_FIELD_MAINT_REASON).(string)

	var clusterId int32
	if clusterId, diags = GetClusterIdFromSchema(d); diags != nil {
		return diags
	}

	var maintStartTm time.Time
	if maintStartTm, err = time.Parse(TIME_FORMAT, maintStartTmStr); err != nil {
		strErr := fmt.Sprintf("%s: %s - Time error", funcName, maintStartTmStr)
		slog.Error(strErr)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  strErr,
		})
		return diags
	}

	var maintStopTm time.Time
	if maintStopTm, err = time.Parse(TIME_FORMAT, maintStopTmStr); err != nil {
		strErr := fmt.Sprintf("%s: %s - Time error", funcName, maintStopTmStr)
		slog.Error(strErr)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  strErr,
		})
		return diags
	}

	maintOperation := *openapi.NewMaintenance(CMON_MAINTENANCE_OPERATION_ADD_MAINT)
	maintOperation.SetClusterId(clusterId)
	maintOperation.SetInitiate(maintStartTm.UTC().String())
	maintOperation.SetDeadline(maintStopTm.UTC().String())
	maintOperation.SetReason(maintReasonStr)

	var resp *http.Response
	if resp, err = apiClient.MaintenanceAPI.MaintenancePost(newCtx).Maintenance(maintOperation).Execute(); err != nil {
		strErr := fmt.Sprintf("%s: Error in Maintenance API call - %v", funcName, resp)
		slog.Error(strErr)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  strErr,
		})
		return diags
	}
	slog.Debug(funcName, "Resp `MaintenancePost.addMaintenance`", resp, "clusterId", clusterId)

	var respBytes []byte
	if respBytes, err = io.ReadAll(resp.Body); err != nil {
		strErr := fmt.Sprintf("%s: Error in io.ReadAll - %s", funcName, err.Error())
		slog.Error(strErr)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  strErr,
		})
		return diags
	}

	var maintOperResp MaintenanceOperationResponse
	if err = json.Unmarshal(respBytes, &maintOperResp); err != nil {
		strErr := fmt.Sprintf("%s: Error in json.Unmarshal - %s", funcName, err.Error())
		slog.Error(strErr)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  strErr,
		})
		return diags
	}
	slog.Debug(funcName, "Resp `MaintenanceOperation`", maintOperResp)

	d.SetId(maintOperResp.UUID)
	d.Set(TF_FIELD_RESOURCE_ID, maintOperResp.UUID)
	d.Set(TF_FIELD_LAST_UPDATED, time.Now().Format(time.RFC822))

	return diags
}

func resourceDeleteDbMaintenance(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	funcName := "resourceDeleteDbMaintenance"
	slog.Debug(funcName)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	var err error

	providerDetails := m.(*ProviderDetails)

	//newCtx := context.WithValue(ctx, "cookie", gNewCtx.Value("cookie"))
	newCtx := context.WithValue(ctx, "cookie", providerDetails.SessionIdCtx.Value("cookie"))

	//apiClient := m.(*openapi.APIClient)
	apiClient := providerDetails.ApiClient

	maintOperation := *openapi.NewMaintenance(CMON_MAINTENANCE_OPERATION_REMOVE_MAINT)
	maintOperation.SetUUID(d.Id())

	var resp *http.Response
	if resp, err = apiClient.MaintenanceAPI.MaintenancePost(newCtx).Maintenance(maintOperation).Execute(); err != nil {
		strErr := fmt.Sprintf("%s: Error in Maintenance API call - %v", funcName, resp)
		slog.Error(strErr)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  strErr,
		})
		return diags
	}
	slog.Debug(funcName, "Resp `MaintenancePost.removeMaintenance`", resp, "maint UUID", d.Id())

	d.Set(TF_FIELD_LAST_UPDATED, time.Now().Format(time.RFC822))
	d.SetId("")

	return diags
}

func resourceReadDbMaintenance(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	funcName := "resourceReadDbMaintenance"
	slog.Debug(funcName)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}

func resourceUpdateDbMaintenance(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	funcName := "resourceUpdateDbMaintenance"
	slog.Debug(funcName)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}
