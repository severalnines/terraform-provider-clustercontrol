package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/severalnines/clustercontrol-client-sdk/go/pkg/openapi"
	"log/slog"
	"time"
)

type DbLoadBalancerInterface interface {
	GetInputs(d *schema.ResourceData, jobData *openapi.JobsJobJobSpecJobData) error
}

func resourceDbLoadBalancer() *schema.Resource {
	funcName := "resourceDbLoadBalancer"
	slog.Debug(funcName)

	return &schema.Resource{
		CreateContext: resourceCreateDbLoadBalancer,
		ReadContext:   resourceReadDbLoadBalancer,
		UpdateContext: resourceUpdateDbLoadBalancer,
		DeleteContext: resourceDeleteDbLoadBalancer,
		Importer:      &schema.ResourceImporter{},
		Schema: map[string]*schema.Schema{
			TF_FIELD_RESOURCE_ID: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "TODO",
			},
			TF_FIELD_LAST_UPDATED: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "TODO",
			},
			TF_FIELD_LB_CREATE: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "TODO",
			},
			TF_FIELD_LB_IMPORT: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "TODO",
			},
			TF_FIELD_CLUSTER_ID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The database cluster ID for which this LB is being deployed to.",
			},
			TF_FIELD_LB_TYPE: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "TODO",
			},
			TF_FIELD_LB_VERSION: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "TODO",
			},
			TF_FIELD_LB_ADMIN_USER: {
				Type: schema.TypeString,
				//Required:    true,
				Optional:    true,
				Description: "TODO",
			},
			TF_FIELD_LB_ADMIN_USER_PW: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "TODO",
				Sensitive:   true,
			},
			TF_FIELD_LB_MONITOR_USER: {
				Type: schema.TypeString,
				//Required:    true,
				Optional:    true,
				Description: "TODO",
			},
			TF_FIELD_LB_MONITOR_USER_PW: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "TODO",
				Sensitive:   true,
			},
			TF_FIELD_LB_PORT: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "TODO",
			},
			TF_FIELD_LB_USE_CLUSTERING: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "TODO",
			},
			TF_FIELD_LB_USE_RW_SPLITTING: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "TODO",
			},
			TF_FIELD_CLUSTER_DISABLE_FW: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "TODO",
			},
			TF_FIELD_LB_INSTALL_SW: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "TODO",
			},
			TF_FIELD_CLUSTER_SSH_USER: {
				Type: schema.TypeString,
				//Optional:    true,
				Required:    true,
				Description: "TODO",
			},
			TF_FIELD_CLUSTER_SSH_PW: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "TODO",
				Sensitive:   true,
			},
			TF_FIELD_CLUSTER_SSH_KEY_FILE: {
				Type: schema.TypeString,
				//Required:    true,
				Optional:    true,
				Description: "SSH Key file.",
			},
			TF_FIELD_CLUSTER_SSH_PORT: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "TODO",
			},
			TF_FIELD_CLUSTER_HOST: {
				Type: schema.TypeList,
				//Required: true,
				Optional:    true,
				Description: "The hosts that make up the cluster.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						TF_FIELD_CLUSTER_HOSTNAME: {
							Type:     schema.TypeString,
							Required: true,
							//Optional:    true,
							Description: "TODO.",
						},
						TF_FIELD_CLUSTER_HOST_PORT: {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "TODO.",
						},
					},
				},
			},
			TF_FIELD_LB_MY_HOST: {
				Type: schema.TypeList,
				//Required: true,
				Optional:    true,
				Description: "The hosts that make up the cluster.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						TF_FIELD_CLUSTER_HOSTNAME: {
							Type:     schema.TypeString,
							Required: true,
							//Optional:    true,
							Description: "TODO.",
						},
						TF_FIELD_CLUSTER_HOST_PORT: {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "TODO.",
						},
					},
				},
			},
		},
	}
}

// Prem
func resourceCreateDbLoadBalancer(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	funcName := "resourceCreateDbLoadBalancer"
	slog.Debug(funcName)

	var diags diag.Diagnostics
	var err error

	newCtx := context.WithValue(ctx, "cookie", gNewCtx.Value("cookie"))

	apiClient := m.(*openapi.APIClient)

	isCreate := d.Get(TF_FIELD_LB_CREATE).(bool)
	isImport := d.Get(TF_FIELD_LB_IMPORT).(bool)
	if !isCreate && !isImport {
		slog.Error(fmt.Sprintf("%s: No work to be done. Create and Import are disabled.", funcName))
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Neither ClueterCreate nor ClusterImport!",
		})
		return diags
	}

	lbType := d.Get(TF_FIELD_LB_TYPE).(string)
	slog.Debug(lbType)

	createLb := NewCCJob(CMON_JOB_CREATE_JOB)
	job := createLb.GetJob()
	jobSpec := job.GetJobSpec()
	jobData := jobSpec.GetJobData()

	var clusterId int32
	if clusterId, diags = GetClusterIdFromSchema(d); diags != nil {
		return diags
	}

	//clusterId := d.Get(TF_FIELD_CLUSTER_ID).(string)
	//if clusterId == "" {
	//	strErr := fmt.Sprintf("%s: %s - not declared", funcName, TF_FIELD_CLUSTER_ID)
	//	slog.Error(strErr)
	//	diags = append(diags, diag.Diagnostic{
	//		Severity: diag.Error,
	//		Summary:  strErr,
	//	})
	//	return diags
	//}
	//
	//if iCid, err = strconv.Atoi(clusterId); err != nil {
	//	strErr := fmt.Sprintf("%s: %s - non-numeric cluster-id", funcName, clusterId)
	//	slog.Error(strErr)
	//	diags = append(diags, diag.Diagnostic{
	//		Severity: diag.Error,
	//		Summary:  strErr,
	//	})
	//	return diags
	//}

	//if clusterId == 0 {
	//	strErr := fmt.Sprintf("%s: - invalid cluster-id 0", funcName)
	//	slog.Error(strErr)
	//	diags = append(diags, diag.Diagnostic{
	//		Severity: diag.Error,
	//		Summary:  strErr,
	//	})
	//	return diags
	//}

	createLb.SetClusterId(clusterId)

	var getInputs DbLoadBalancerInterface
	switch lbType {
	case LOAD_BLANCER_TYPE_PROXYSQL:
		jobSpec.SetCommand(CMON_JOB_CREATE_PROXYSQL_COMMAND)
		getInputs = NewProxySql()
	case LOAD_BLANCER_TYPE_HAPROXY:
		jobSpec.SetCommand(CMON_JOB_CREATE_HAPROXY_COMMAND)
		getInputs = NewHAProxy()
	default:
		slog.Warn(funcName, "Unsupported LB type", lbType)
	}

	if getInputs != nil {
		if err = getInputs.GetInputs(d, &jobData); err != nil {
			slog.Error(err.Error())
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error getting inputs for LoadBalancerCreate",
			})
			return diags
		}
	}

	jobSpec.SetJobData(jobData)
	job.SetJobSpec(jobSpec)
	createLb.SetJob(job)

	if err = SendAndWaitForJobCompletion(newCtx, apiClient, createLb); err != nil {
		slog.Error(err.Error())
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Job Failed for ClusterCreate",
		})
		return diags
	}

	lbHostname := jobData.GetHostname()
	resourceId := fmt.Sprintf("%s;%s;%s", clusterId, lbType, lbHostname)
	d.SetId(resourceId)
	d.Set(TF_FIELD_RESOURCE_ID, resourceId)
	d.Set(TF_FIELD_LAST_UPDATED, time.Now().Format(time.RFC850))

	return diags
}

// Prem
func resourceReadDbLoadBalancer(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	funcName := "resourceReadDbCluster"
	slog.Debug(funcName)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	//var cluster *openapi.ClusterResponse
	//var err error
	//
	//newCtx := context.WithValue(ctx, "cookie", gNewCtx.Value("cookie"))
	//
	//apiClient := m.(*openapi.APIClient)

	//d.Set("last_updated", time.Now().Format(time.RFC850))

	return diags
}

// Prem
func resourceUpdateDbLoadBalancer(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	funcName := "resourceReadDbCluster"
	slog.Debug(funcName)

	// Warning or errors can be collected in a slice type
	//var cluster *openapi.ClusterResponse
	//var err error
	//
	//newCtx := context.WithValue(ctx, "cookie", gNewCtx.Value("cookie"))
	//
	//apiClient := m.(*openapi.APIClient)

	//d.Set("last_updated", time.Now().Format(time.RFC850))

	return resourceReadDbCluster(ctx, d, m)
}

// Prem
func resourceDeleteDbLoadBalancer(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	funcName := "resourceDeleteDbCluster"
	slog.Debug(funcName)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	//newCtx := context.WithValue(ctx, "cookie", gNewCtx.Value("cookie"))
	//
	//apiClient := m.(*openapi.APIClient)

	d.SetId("")

	return diags
}
