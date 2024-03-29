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
	GetInputs(d map[string]any, jobData *openapi.JobsJobJobSpecJobData) error
}

// ********************************************
// NOTE: the rest of the code is no longer used !!!
// ********************************************
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
				Description: "Whether to create this resource or not?",
			},
			TF_FIELD_LB_IMPORT: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to import this resource or not?",
			},
			TF_FIELD_CLUSTER_ID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The database cluster ID for which this LB is being deployed to.",
			},
			TF_FIELD_LB_TYPE: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The load balancer type (e.g., proxysql, haproxy, etc)",
			},
			TF_FIELD_LB_VERSION: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Software version",
			},
			TF_FIELD_LB_ADMIN_USER: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The load balancer admin user",
			},
			TF_FIELD_LB_ADMIN_USER_PW: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The load balancer admin user's password",
				Sensitive:   true,
			},
			TF_FIELD_LB_MONITOR_USER: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The load balancer monitor user (only applicable to proxysql)",
			},
			TF_FIELD_LB_MONITOR_USER_PW: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The load balancer monitor user's password",
				Sensitive:   true,
			},
			TF_FIELD_LB_PORT: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The load balancer port that it will accept connections on behalf of the database it is front-ending.",
			},
			TF_FIELD_LB_USE_CLUSTERING: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to use ProxySQL clustering or not. Only applicable to ProxySQL at this time",
			},
			TF_FIELD_LB_USE_RW_SPLITTING: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to Read/Write splitting for queries or not?",
			},
			TF_FIELD_CLUSTER_DISABLE_FW: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Disable firewall on the host OS when installing DB packages.",
			},
			TF_FIELD_CLUSTER_DISABLE_SELINUX: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Disable SELinux on the host OS when installing DB packages.",
			},
			TF_FIELD_LB_INSTALL_SW: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Install DB packages from respective repos",
			},
			TF_FIELD_LB_ENABLE_UNINSTALL: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "When removing DB cluster from ClusterControl, enable uinstalling DB packages.",
			},
			TF_FIELD_CLUSTER_SSH_USER: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The SSH user ClusterControl will use to SSH to the DB host from the ClusterControl host",
			},
			TF_FIELD_CLUSTER_SSH_PW: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sudo user's password. If sudo user doesn't have a password, leave this field blank",
				Sensitive:   true,
			},
			TF_FIELD_CLUSTER_SSH_KEY_FILE: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "SSH Key file. The path to the private key file for the Sudo user on the ClusterControl host.",
			},
			TF_FIELD_CLUSTER_SSH_PORT: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ssh port.",
			},
			TF_FIELD_CLUSTER_HOST: {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Database hosts that make up the cluster.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						TF_FIELD_CLUSTER_HOSTNAME: {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Hostname/IP of the DB host behind this load balancer. Can be IP address as well.",
						},
						TF_FIELD_CLUSTER_HOST_PORT: {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The port the DB host behind this load balancer.",
						},
					},
				},
			},
			TF_FIELD_LB_MY_HOST: {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The load balancer host in question (i.e, self)",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						TF_FIELD_CLUSTER_HOSTNAME: {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Hostname/IP of this load balancer. Can be IP address as well.",
						},
						TF_FIELD_CLUSTER_HOST_PORT: {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The port of this load balancer.",
						},
					},
				},
			},
		},
	}
}

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
		str := fmt.Sprintf("%s: No work to be done. Create and Import are disabled.", funcName)
		slog.Info(str)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  str,
		})
		return diags
	}

	if isImport && !isCreate {
		str := "Importing a load balancer into ClusterControl is not supported at this time."
		slog.Info(str)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  str,
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
		str := fmt.Sprintf("%s - Unknown load balancer type: %s", funcName, lbType)
		slog.Warn(str)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  str,
		})
		return diags
	}

	if getInputs != nil {
		//if err = getInputs.GetInputs(d, &jobData); err != nil {
		//	slog.Error(err.Error())
		//	diags = append(diags, diag.Diagnostic{
		//		Severity: diag.Error,
		//		Summary:  "Error getting inputs for LoadBalancerCreate",
		//	})
		//	return diags
		//}
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

func resourceReadDbLoadBalancer(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	funcName := "resourceReadDbLoadBalancer"
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

func resourceUpdateDbLoadBalancer(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	funcName := "resourceUpdateDbLoadBalancer"
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

func resourceDeleteDbLoadBalancer(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	funcName := "resourceDeleteDbLoadBalancer"
	slog.Debug(funcName)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	//newCtx := context.WithValue(ctx, "cookie", gNewCtx.Value("cookie"))
	//
	//apiClient := m.(*openapi.APIClient)

	d.SetId("")

	return diags
}
