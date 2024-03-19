package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/severalnines/clustercontrol-client-sdk/go/pkg/openapi"
	"log/slog"
	"strconv"
	"time"
)

type DbClusterInterface interface {
	GetInputs(d *schema.ResourceData, jobData *openapi.JobsJobJobSpecJobData) error
}

func resourceDbCluster() *schema.Resource {
	funcName := "resourceDbCluster"
	slog.Debug(funcName)

	return &schema.Resource{
		CreateContext: resourceCreateDbCluster,
		ReadContext:   resourceReadDbCluster,
		UpdateContext: resourceUpdateDbCluster,
		DeleteContext: resourceDeleteDbCluster,
		Importer:      &schema.ResourceImporter{},
		Schema: map[string]*schema.Schema{
			TF_FIELD_CLUSTER_CREATE: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "TODO",
			},
			TF_FIELD_CLUSTER_IMPORT: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "TODO",
			},
			TF_FIELD_CLUSTER_NAME: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "TODO: The name of the resource, also acts as it's unique ID",
			},
			TF_FIELD_CLUSTER_TYPE: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "TODO",
			},
			TF_FIELD_CLUSTER_VENDOR: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "TODO",
			},
			TF_FIELD_CLUSTER_VERSION: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "TODO",
			},
			TF_FIELD_CLUSTER_ADMIN_USER: {
				Type: schema.TypeString,
				//Required:    true,
				Optional:    true,
				Description: "TODO",
			},
			TF_FIELD_CLUSTER_ADMIN_PW: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "TODO",
				Sensitive:   true,
			},
			TF_FIELD_CLUSTER_PORT: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "TODO",
			},
			TF_FIELD_CLUSTER_DATA_DIR: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "TODO",
			},
			// TODO: perhaps can be removed later
			//TF_FIELD_CLUSTER_CFG_TEMPLATE: {
			//	Type:        schema.TypeString,
			//	Optional:    true,
			//	Description: "TODO",
			//},
			TF_FIELD_CLUSTER_DISABLE_FW: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "TODO",
			},
			TF_FIELD_CLUSTER_INSTALL_SW: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "TODO",
			},
			// TODO: need to remove from here. Moved inside `host` sub-section
			//TF_FIELD_CLUSTER_SYNC_REP: {
			//	Type:        schema.TypeBool,
			//	Optional:    true,
			//	Description: "TODO",
			//},
			TF_FIELD_CLUSTER_SEMISYNC_REP: {
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
			TF_FIELD_CLUSTER_SNAPSHOT_LOC: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "TODO",
			},
			TF_FIELD_CLUSTER_SNAPSHOT_REPO: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "TODO",
			},
			TF_FIELD_CLUSTER_SNAPSHOT_HOST: {
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
						TF_FIELD_CLUSTER_HOSTNAME_DATA: {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "TODO.",
						},
						TF_FIELD_CLUSTER_HOSTNAME_INT: {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "TODO.",
						},
						TF_FIELD_CLUSTER_HOST_PORT: {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "TODO.",
						},
						TF_FIELD_CLUSTER_SYNC_REP: {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "TODO",
						},
						TF_FIELD_CLUSTER_HOSTNAME_DD: {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "TODO.",
						},
						TF_FIELD_CLUSTER_HOST_PROTO: {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "TODO.",
						},
						TF_FIELD_CLUSTER_HOST_ROLES: {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "TODO.",
						},
					},
				},
			},
			TF_FIELD_CLUSTER_REPLICA_SET: {
				Type: schema.TypeList,
				//Required: true,
				Optional:    true,
				Description: "The hosts that make up the cluster.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						TF_FIELD_CLUSTER_REPLICA_SET_RS: {
							Type:     schema.TypeString,
							Required: true,
							//Optional:    true,
							Description: "TODO.",
						},
						TF_FIELD_CLUSTER_REPLICA_MEMBER: {
							Type:     schema.TypeList,
							Required: true,
							//Optional:    true,
							Description: "The hosts that make up the cluster.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									TF_FIELD_CLUSTER_HOSTNAME: {
										Type:     schema.TypeString,
										Required: true,
										//Optional:    true,
										Description: "TODO.",
									},
									TF_FIELD_CLUSTER_HOSTNAME_DATA: {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "TODO.",
									},
									TF_FIELD_CLUSTER_HOSTNAME_INT: {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "TODO.",
									},
									TF_FIELD_CLUSTER_HOST_PORT: {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "TODO.",
									},
									TF_FIELD_CLUSTER_HOST_SLAVE_DELAY: {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "TODO.",
									},
									TF_FIELD_CLUSTER_HOST_ARBITER_ONLY: {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "TODO.",
									},
									TF_FIELD_CLUSTER_HOST_HIDDEN: {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "TODO.",
									},
									TF_FIELD_CLUSTER_HOST_PRIORITY: {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "TODO.",
									},
								},
							},
						},
					},
				},
			},
			TF_FIELD_CLUSTER_MONGO_CONFIG_SERVER: {
				Type: schema.TypeList,
				//Required: true,
				Optional:    true,
				Description: "The hosts that make up the cluster.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						TF_FIELD_CLUSTER_REPLICA_SET_RS: {
							Type:     schema.TypeString,
							Required: true,
							//Optional:    true,
							Description: "TODO.",
						},
						TF_FIELD_CLUSTER_REPLICA_MEMBER: {
							Type:     schema.TypeList,
							Required: true,
							//Optional:    true,
							Description: "The hosts that make up the cluster.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									TF_FIELD_CLUSTER_HOSTNAME: {
										Type:     schema.TypeString,
										Required: true,
										//Optional:    true,
										Description: "TODO.",
									},
									TF_FIELD_CLUSTER_HOSTNAME_DATA: {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "TODO.",
									},
									TF_FIELD_CLUSTER_HOSTNAME_INT: {
										Type:        schema.TypeString,
										Optional:    true,
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
				},
			},
			TF_FIELD_CLUSTER_MONGOS_SERVER: {
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
						TF_FIELD_CLUSTER_HOSTNAME_DATA: {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "TODO.",
						},
						TF_FIELD_CLUSTER_HOSTNAME_INT: {
							Type:        schema.TypeString,
							Optional:    true,
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
			TF_FIELD_CLUSTER_TOPOLOGY: {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "TODO",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						TF_FIELD_CLUSTER_PRIMARY: {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "TODO.",
						},
						TF_FIELD_CLUSTER_REPLICA: {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "TODO.",
						},
					},
				},
			},
			TF_FIELD_CLUSTER_TAGS: {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "An optional list of tags, represented as a key, value pair",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			TF_FIELD_CLUSTER_TIMEOUTS: {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "TODO",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

// TODO: remove after adding support timeout in the provider plugin - galera deployment takes for ever ... !
//
//	//create_timeout := d.Timeout("create").(string)
//	//import_timeout := d.Timeout("import").(string)
//	//slog.Debug("Create timeout: %s", create_timeout)
//	//slog.Debug("Import timeout: %s", import_timeout)
//
// Prem
func resourceCreateDbCluster(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	funcName := "resourceCreateDbCluster"
	slog.Debug(funcName)

	var diags diag.Diagnostics
	var err error

	newCtx := context.WithValue(ctx, "cookie", gNewCtx.Value("cookie"))

	apiClient := m.(*openapi.APIClient)

	isCreate := d.Get(TF_FIELD_CLUSTER_CREATE).(bool)
	isImport := d.Get(TF_FIELD_CLUSTER_IMPORT).(bool)
	if !isCreate && !isImport {
		slog.Info(fmt.Sprintf("%s: No work to be done. Create and Import are disabled.", funcName))
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Neither ClueterCreate nor ClusterImport!",
		})
		return diags
	}

	createCluster := NewCCJob(CMON_JOB_CREATE_JOB)
	job := createCluster.GetJob()
	jobSpec := job.GetJobSpec()
	jobSpec.SetCommand(CMON_JOB_CREATE_CLUSTER_COMMAND)
	jobData := jobSpec.GetJobData()

	clusterType := d.Get(TF_FIELD_CLUSTER_TYPE).(string)
	slog.Debug(clusterType)

	var getInputs DbClusterInterface = nil

	switch clusterType {
	case CLUSTER_TYPE_REPLICATION:
		getInputs = NewMySQLMaria()
	case CLUSTER_TYPE_GALERA:
		getInputs = NewMySQLMaria()
	case CLUSTER_TYPE_PG_SINGLE:
		getInputs = NewPostgres()
	case CLUSTER_TYPE_MOGNODB:
		getInputs = NewMongo()
	case CLUSTER_TYPE_REDIS:
		getInputs = NewRedis()
	case CLUSTER_TYPE_MSSQL_SINGLE:
		getInputs = NewMsSql()
	case CLUSTER_TYPE_MSSQL_AO_ASYNC:
		getInputs = NewMsSql()
	case CLUSTER_TYPE_ELASTIC:
		getInputs = NewElastic()
	default:
		slog.Warn(funcName, "Unknown cluster type", clusterType)
	}

	if getInputs != nil {
		if err = getInputs.GetInputs(d, &jobData); err != nil {
			slog.Error(err.Error())
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error getting inputs for ClusterCreate",
			})
			return diags
		}
	}

	jobSpec.SetJobData(jobData)
	job.SetJobSpec(jobSpec)
	createCluster.SetJob(job)

	if err = SendAndWaitForJobCompletion(newCtx, apiClient, createCluster); err != nil {
		slog.Error(err.Error())
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Job Failed for ClusterCreate",
		})
		return diags
	}

	clusterName := jobData.GetClusterName()

	var clusterId int32 = -1
	if clusterId, err = GetClusterIdByClusterName(newCtx, apiClient, clusterName); err != nil {
		slog.Error(err.Error())
		return diag.FromErr(err)
	}

	// update resource can ask to change name, which is a valid ask.
	d.SetId(strconv.Itoa(int(clusterId)))
	//d.SetId(clusterName)

	return diags
}

func resourceReadDbCluster(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	funcName := "resourceReadDbCluster"
	slog.Debug(funcName)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	var cluster *openapi.ClusterResponse
	var err error

	newCtx := context.WithValue(ctx, "cookie", gNewCtx.Value("cookie"))

	apiClient := m.(*openapi.APIClient)

	clusterId := d.Id()

	if cluster, err = GetClusterByClusterStrId(newCtx, apiClient, clusterId); err != nil {
		// TODO
	}

	cluster.GetClusterType()
	cluster.GetClusterId()
	cluster.GetVendor()
	cluster.GetVersion()
	cluster.GetTags()

	//foo := d.Get("foo")
	//
	//for _, f := range foo.([]any) {
	//	f := f.(map[string]any)
	//
	//	for _, b := range f["bar"].([]any) {
	//		b := b.(map[string]any)
	//		b["version"] = uuid.New().String()
	//	}
	//}
	//
	//if err := d.Set("foo", foo); err != nil {
	//	diags = append(diags, diag.Diagnostic{
	//		Severity: diag.Error,
	//		Summary:  err.Error(),
	//	})
	//	return diags
	//}

	d.Set("last_updated", time.Now().Format(time.RFC850))

	return diags
}

func resourceUpdateDbCluster(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	funcName := "resourceReadDbCluster"
	slog.Debug(funcName)

	// Warning or errors can be collected in a slice type
	var cluster *openapi.ClusterResponse
	var err error

	newCtx := context.WithValue(ctx, "cookie", gNewCtx.Value("cookie"))

	apiClient := m.(*openapi.APIClient)

	clusterId := d.Id()

	if cluster, err = GetClusterByClusterStrId(newCtx, apiClient, clusterId); err != nil {
		// TODO
	}

	cluster.GetClusterType()
	cluster.GetClusterId()
	cluster.GetVendor()
	cluster.GetVersion()
	cluster.GetTags()

	// Warning or errors can be collected in a slice type
	//var diags diag.Diagnostics
	//
	//resourceID := d.Id()
	//
	//if d.HasChange("foo") {
	//	foo := d.Get("foo").([]any)
	//	d.Set("last_updated", time.Now().Format(time.RFC850))
	//}

	d.Set("last_updated", time.Now().Format(time.RFC850))

	return resourceReadDbCluster(ctx, d, m)
}

// Prem
func resourceDeleteDbCluster(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	funcName := "resourceDeleteDbCluster"
	slog.Debug(funcName)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	newCtx := context.WithValue(ctx, "cookie", gNewCtx.Value("cookie"))

	apiClient := m.(*openapi.APIClient)

	delCluster := NewCCJob(CMON_JOB_CREATE_JOB)
	job := delCluster.GetJob()
	jobSpec := job.GetJobSpec()
	jobSpec.SetCommand(CMON_JOB_REMOVE_CLUSTER_COMMAND)
	jobData := jobSpec.GetJobData()

	var err error
	var clusterId int = -1
	clusterId, err = strconv.Atoi(d.Id())

	jobData.SetClusterid(int32(clusterId))
	jobSpec.SetJobData(jobData)
	job.SetJobSpec(jobSpec)
	delCluster.SetJob(job)

	if err = SendAndWaitForJobCompletion(newCtx, apiClient, delCluster); err != nil {
		slog.Error(err.Error())
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Job Failed for ClusterDelete",
		})
		return diags
	}

	d.SetId("")

	return diags
}
