package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/severalnines/clustercontrol-client-sdk/go/pkg/openapi"
	"log/slog"
	"strconv"
	"strings"
	"time"
)

type DbClusterInterface interface {
	GetInputs(d *schema.ResourceData, jobData *openapi.JobsJobJobSpecJobData) error
	HandleRead(ctx context.Context, d *schema.ResourceData, m interface{}, clusterInfo *openapi.ClusterResponse) error
	IsUpdateBatchAllowed(d *schema.ResourceData) error
	HandleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}, clusterInfo *openapi.ClusterResponse) error
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
			TF_FIELD_CLUSTER_ID: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "TODO",
			},
			TF_FIELD_LAST_UPDATED: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "TODO",
			},
			TF_FIELD_CLUSTER_CREATE: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to create this resource or not?",
			},
			TF_FIELD_CLUSTER_IMPORT: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to import this resource or not?",
			},
			TF_FIELD_CLUSTER_NAME: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the database cluster.",
			},
			TF_FIELD_CLUSTER_TYPE: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Type of cluster - replication, galera, postgresql_single (single is a misnomer), etc",
			},
			TF_FIELD_CLUSTER_VENDOR: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Database vendor - oracle, percona, mariadb, 10gen, microsoft, redis, elasticsearch, for postgresql it is `default` etc",
			},
			TF_FIELD_CLUSTER_VERSION: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The database version",
			},
			TF_FIELD_CLUSTER_ADMIN_USER: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name for the admin/root user for the database",
			},
			TF_FIELD_CLUSTER_ADMIN_PW: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Password for the admin/root user for the database. Note that this may show up in logs, and it will be stored in the state file",
				Sensitive:   true,
			},
			TF_FIELD_CLUSTER_PORT: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The port on which the DB will accepts connections",
			},
			TF_FIELD_CLUSTER_SENTINEL_PORT: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The port Redis Sentinel uses to communicate",
			},
			TF_FIELD_CLUSTER_DATA_DIR: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The data directory for the database data files. If not set explicily, the default for the respective DB vendor will be chosen",
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
			TF_FIELD_CLUSTER_INSTALL_SW: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Install DB packages from respective repos",
			},
			TF_FIELD_CLUSTER_ENABLE_UNINSTALL: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "When removing DB cluster from ClusterControl, enable uinstalling DB packages.",
			},
			TF_FIELD_CLUSTER_SEMISYNC_REP: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Semi-synchronous replication for MySQL and MariaDB non-galera clusters",
			},
			TF_FIELD_CLUSTER_PG_TIMESALE_EXT: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to setup TimescaleDB extension or not",
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
				Required:    true,
				Description: "SSH Key file. The path to the private key file for the Sudo user on the ClusterControl host",
			},
			TF_FIELD_CLUSTER_SSH_PORT: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ssh port.",
			},
			TF_FIELD_CLUSTER_SNAPSHOT_LOC: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Elasticsearch snapshot location",
			},
			TF_FIELD_CLUSTER_SNAPSHOT_REPO: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Elasticsearch snapshot repository",
			},
			TF_FIELD_CLUSTER_SNAPSHOT_HOST: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Elasticsearch snapshot host",
			},
			TF_FIELD_CLUSTER_MONGO_AUTH_DB: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The mongodb database to use for authentication purposes",
			},
			// ****************************
			// Database host attributes
			// ****************************
			TF_FIELD_CLUSTER_HOST: {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The list of nodes/hosts that make up the cluster",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						TF_FIELD_CLUSTER_HOSTNAME: {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Hostname of the DB host. Can be IP address as well",
						},
						TF_FIELD_CLUSTER_HOSTNAME_DATA: {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Hostname/IP used for data comms (may be legacy ClusterControl).",
						},
						TF_FIELD_CLUSTER_HOSTNAME_INT: {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "If there's a private net that all DB hosts can communicate, use it here.",
						},
						TF_FIELD_CLUSTER_HOST_PORT: {
							Type:     schema.TypeString,
							Optional: true,
							Description: "The port on which the DB server will listen for connections. If one is not provided, " +
								"default for the DB type will be used, or inherited from earlier/top-level specification.",
						},
						TF_FIELD_CLUSTER_SYNC_REP: {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Applicable to PostgreSQL hot-standby nodes only. Use synchronous replication (or  not)",
						},
						TF_FIELD_CLUSTER_HOST_DD: {
							Type:     schema.TypeString,
							Optional: true,
							Description: "The data directory for the database data files. If not set explicily, default " +
								"for the DB type will be used, or inherited from earlier/top-level specification.",
						},
						TF_FIELD_CLUSTER_HOST_PROTO: {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "TODO.",
						},
						TF_FIELD_CLUSTER_HOST_ROLES: {
							Type:     schema.TypeString,
							Optional: true,
							Description: "Applicable to Elasticsearch - the role of this host (master-data: host will " +
								"be designated as the master node and a data node, etc)",
						},
					},
				},
			},
			// ****************************
			// MongoDB Replicaset specific attributes
			// ****************************
			TF_FIELD_CLUSTER_REPLICA_SET: {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The hosts that make up the cluster.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						TF_FIELD_CLUSTER_REPLICA_SET_RS: {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The replicaset's name.",
						},
						TF_FIELD_CLUSTER_REPLICA_MEMBER: {
							Type:        schema.TypeList,
							Required:    true,
							Description: "The hosts that make up the replicaset HA nodes.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									TF_FIELD_CLUSTER_HOSTNAME: {
										Type:     schema.TypeString,
										Required: true,
										//Optional:    true,
										Description: "Hostname of the DB host. Can be IP address as well",
									},
									TF_FIELD_CLUSTER_HOSTNAME_DATA: {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Hostname/IP used for data comms (may be legacy ClusterControl).",
									},
									TF_FIELD_CLUSTER_HOSTNAME_INT: {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "If there's a private net that all DB hosts can communicate, use it here.",
									},
									TF_FIELD_CLUSTER_HOST_PORT: {
										Type:     schema.TypeString,
										Optional: true,
										Description: "The port on which the DB server will listen for connections. If " +
											"one is not provided, the default for the DB type will be used.",
									},
									TF_FIELD_CLUSTER_HOST_ARBITER_ONLY: {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "The host is acting as an arbiter only.",
									},
									TF_FIELD_CLUSTER_HOST_PRIORITY: {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Priority of the host in the mongo replication setup.",
									},
									TF_FIELD_CLUSTER_HOST_HIDDEN: {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "TODO.",
									},
									TF_FIELD_CLUSTER_HOST_SLAVE_DELAY: {
										Type:     schema.TypeString,
										Optional: true,
										Description: "Used in non-galera MySQL/MariaDB standby setup. Specifies the lag " +
											"for the slave.",
									},
								},
							},
						},
					},
				},
			},
			// ****************************
			// MongoDB Config Server specific attributes
			// ****************************
			TF_FIELD_CLUSTER_MONGO_CONFIG_SERVER: {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specification for the MongoDB Configuration Server.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						TF_FIELD_CLUSTER_REPLICA_SET_RS: {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The replicaset's name.",
						},
						TF_FIELD_CLUSTER_REPLICA_MEMBER: {
							Type:        schema.TypeList,
							Required:    true,
							Description: "The host that make up the replicaset member.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									TF_FIELD_CLUSTER_HOSTNAME: {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Hostname of the DB host. Can be IP address as well.",
									},
									TF_FIELD_CLUSTER_HOSTNAME_DATA: {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Hostname/IP used for data comms (may be legacy ClusterControl).",
									},
									TF_FIELD_CLUSTER_HOSTNAME_INT: {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "If there's a private net that all DB hosts can communicate, use it here.",
									},
									TF_FIELD_CLUSTER_HOST_PORT: {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Port for the config server. If one is not provided, the default will be used.",
									},
								},
							},
						},
					},
				},
			},
			// ****************************
			// MongoDB Mongos Server specific attributes
			// ****************************
			TF_FIELD_CLUSTER_MONGOS_SERVER: {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specification for the MongoDB mongos Server.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						TF_FIELD_CLUSTER_HOSTNAME: {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Hostname of the DB host. Can be IP address as well.",
						},
						TF_FIELD_CLUSTER_HOSTNAME_DATA: {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Hostname/IP used for data comms (may be legacy ClusterControl).",
						},
						TF_FIELD_CLUSTER_HOSTNAME_INT: {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "If there's a private net that all DB hosts can communicate, use it here.",
						},
						TF_FIELD_CLUSTER_HOST_PORT: {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Port for the config server. If one is not provided, the default will be used.",
						},
					},
				},
			},
			// ****************************
			// MySQL Master/Slave topology speification attributes
			// ****************************
			TF_FIELD_CLUSTER_TOPOLOGY: {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Only applicable to MySQL/MariaDB non-galera clusters. A way to specify Master and Slave(s).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						TF_FIELD_CLUSTER_PRIMARY: {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The Master host",
						},
						TF_FIELD_CLUSTER_REPLICA: {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The Slave host.",
						},
					},
				},
			},
			TF_FIELD_CLUSTER_TAGS: {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Tags to associate with a DB cluster. The tags are only relevant in the ClusterControl domain.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			TF_FIELD_CLUSTER_DEPLOY_AGENTS: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Automatically deploy prometheus and other relevant agents after setting up the intial DB cluster.",
			},
			TF_FIELD_CLUSTER_AUTO_RECOVERY: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Have cluster auto-recovery on (or off)",
			},
			TF_FIELD_CLUSTER_SSL: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable SSL based comms between the cluster nodes and client access to node.",
			},
			TF_FIELD_CLUSTER_ENABLE_PGM_AGENT: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable percona backup for mongodb.",
			},
			TF_FIELD_CLUSTER_PBM_BACKUP_DIR: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Backup dir, nfs mounted directory / path for PBM backup.",
			},
			TF_FIELD_CLUSTER_ENABLE_PGBACKREST_AGENT: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable PgBackRest for Postgres.",
			},
			// ****************************
			// Database load balancer attributes
			// ****************************
			TF_FIELD_CLUSTER_LOAD_BALANCER: {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The list of nodes/hosts that make up the cluster",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						//TF_FIELD_CLUSTER_HOST: {
						//	Type:        schema.TypeList,
						//	Optional:    true,
						//	Description: "The Database hosts that make up the cluster.",
						//	Elem: &schema.Resource{
						//		Schema: map[string]*schema.Schema{
						//			TF_FIELD_CLUSTER_HOSTNAME: {
						//				Type:        schema.TypeString,
						//				Required:    true,
						//				Description: "Hostname/IP of the DB host behind this load balancer. Can be IP address as well.",
						//			},
						//			TF_FIELD_CLUSTER_HOST_PORT: {
						//				Type:        schema.TypeString,
						//				Optional:    true,
						//				Description: "The port the DB host behind this load balancer.",
						//			},
						//		},
						//	},
						//},
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
				},
			},
		},
	}
}

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
		str := fmt.Sprintf("%s: No work to be done. Create and Import are disabled.", funcName)
		slog.Info(str)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  str,
		})
		return diags
	}

	if isImport && !isCreate {
		str := "Importing a cluster into ClusterControl is not supported at this time."
		slog.Info(str)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  str,
		})
		return diags
	}

	createCluster := NewCCJob(CMON_JOB_CREATE_JOB)
	job := createCluster.GetJob()
	jobSpec := job.GetJobSpec()
	jobSpec.SetCommand(CMON_JOB_CREATE_CLUSTER_COMMAND)
	jobData := jobSpec.GetJobData()

	extClusterType := d.Get(TF_FIELD_CLUSTER_TYPE).(string)
	clusterType, ok := gExtClusterTypeToIntClusterTypeMap[extClusterType]
	if !ok {
		str := fmt.Sprintf("Unsupported cluster-type: %s", extClusterType)
		slog.Info(str)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  str,
		})
		return diags
	}
	slog.Info(clusterType)

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
		str := fmt.Sprintf("%s - Unknown cluster type: %s", funcName, clusterType)
		slog.Warn(str)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  str,
		})
		return diags
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
	id := strconv.Itoa(int(clusterId))
	d.SetId(id)
	d.Set(TF_FIELD_CLUSTER_ID, id)
	d.Set(TF_FIELD_LAST_UPDATED, time.Now().Format(time.RFC822))
	//d.SetId(clusterName)

	// *************************
	// Create the load balancer
	// *************************

	loadBalancers := d.Get(TF_FIELD_CLUSTER_LOAD_BALANCER)
	for _, ff := range loadBalancers.([]any) {
		f := ff.(map[string]any)

		createLb := NewCCJob(CMON_JOB_CREATE_JOB)
		lbJob := createLb.GetJob()
		lbJobSpec := lbJob.GetJobSpec()
		lbJobData := lbJobSpec.GetJobData()
		createLb.SetClusterId(clusterId)

		lbType := f[TF_FIELD_LB_TYPE].(string)
		var getLbInputs DbLoadBalancerInterface
		switch lbType {
		case LOAD_BLANCER_TYPE_PROXYSQL:
			jobSpec.SetCommand(CMON_JOB_CREATE_PROXYSQL_COMMAND)
			getLbInputs = NewProxySql()
		case LOAD_BLANCER_TYPE_HAPROXY:
			jobSpec.SetCommand(CMON_JOB_CREATE_HAPROXY_COMMAND)
			getLbInputs = NewHAProxy()
		default:
			str := fmt.Sprintf("%s - Unknown load balancer type: %s", funcName, lbType)
			slog.Warn(str)
			// Just because LB creation failed, we will not fail a successful DB cluster deployment !!!
			return diags
		}

		if getInputs != nil {
			if err = getLbInputs.GetInputs(f, &lbJobData); err != nil {
				slog.Error(err.Error())
				// Just because LB creation failed, we will not fail a successful DB cluster deployment !!!
				return diags
			}
		}

		lbJobSpec.SetJobData(lbJobData)
		lbJob.SetJobSpec(lbJobSpec)
		createLb.SetJob(lbJob)

		if err = SendAndWaitForJobCompletion(newCtx, apiClient, createLb); err != nil {
			slog.Error(err.Error())
			// Just because LB creation failed, we will not fail a successful DB cluster deployment !!!
			return diags
		}

	} // For each Load balancer

	return diags
}

func resourceReadDbCluster(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	funcName := "resourceReadDbCluster"
	slog.Debug(funcName)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	var clusterInfo *openapi.ClusterResponse
	var err error

	newCtx := context.WithValue(ctx, "cookie", gNewCtx.Value("cookie"))

	apiClient := m.(*openapi.APIClient)

	clusterId := d.Id()

	if clusterInfo, err = GetClusterByClusterStrId(newCtx, apiClient, clusterId); err != nil {
		slog.Error(err.Error())
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  err.Error(),
		})
		return diags
	}

	//clusterType := d.Get(TF_FIELD_CLUSTER_TYPE).(string)
	clusterType := clusterInfo.GetClusterType()
	clusterType = strings.ToLower(clusterType)
	slog.Debug(clusterType)

	var readHandler DbClusterInterface = nil

	switch clusterType {
	case CLUSTER_TYPE_REPLICATION:
		readHandler = NewMySQLMaria()
	case CLUSTER_TYPE_GALERA:
		readHandler = NewMySQLMaria()
	case CLUSTER_TYPE_PG_SINGLE:
		readHandler = NewPostgres()
	case CLUSTER_TYPE_MOGNODB:
		readHandler = NewMongo()
	case CLUSTER_TYPE_REDIS:
		readHandler = NewRedis()
	case CLUSTER_TYPE_MSSQL_SINGLE:
		readHandler = NewMsSql()
	case CLUSTER_TYPE_MSSQL_AO_ASYNC:
		readHandler = NewMsSql()
	case CLUSTER_TYPE_ELASTIC:
		readHandler = NewElastic()
	default:
		str := fmt.Sprintf("%s - Unknown cluster type: %s", funcName, clusterType)
		slog.Warn(str)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  str,
		})
		return diags
	}

	if readHandler != nil {
		if err = readHandler.HandleRead(newCtx, d, m, clusterInfo); err != nil {
			slog.Error(err.Error())
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error in DB cluster Read handler",
			})
			return diags
		}
	}

	d.Set(TF_FIELD_LAST_UPDATED, time.Now().Format(time.RFC822))

	return diags
}

func resourceUpdateDbCluster(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	funcName := "resourceUpdateDbCluster"
	slog.Debug(funcName)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	var clusterInfo *openapi.ClusterResponse
	var err error

	newCtx := context.WithValue(ctx, "cookie", gNewCtx.Value("cookie"))

	apiClient := m.(*openapi.APIClient)

	clusterId := d.Id()

	if clusterInfo, err = GetClusterByClusterStrId(newCtx, apiClient, clusterId); err != nil {
		slog.Error(err.Error())
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  err.Error(),
		})
		return diags
	}

	//clusterType := d.Get(TF_FIELD_CLUSTER_TYPE).(string)
	clusterType := clusterInfo.GetClusterType()
	clusterType = strings.ToLower(clusterType)
	slog.Debug(clusterType)

	var updateHandler DbClusterInterface = nil

	switch clusterType {
	case CLUSTER_TYPE_REPLICATION:
		updateHandler = NewMySQLMaria()
	case CLUSTER_TYPE_GALERA:
		updateHandler = NewMySQLMaria()
	case CLUSTER_TYPE_PG_SINGLE:
		updateHandler = NewPostgres()
	case CLUSTER_TYPE_MOGNODB:
		updateHandler = NewMongo()
	case CLUSTER_TYPE_REDIS:
		updateHandler = NewRedis()
	case CLUSTER_TYPE_MSSQL_SINGLE:
		updateHandler = NewMsSql()
	case CLUSTER_TYPE_MSSQL_AO_ASYNC:
		updateHandler = NewMsSql()
	case CLUSTER_TYPE_ELASTIC:
		updateHandler = NewElastic()
	default:
		str := fmt.Sprintf("%s - Unknown cluster type: %s", funcName, clusterType)
		slog.Warn(str)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  str,
		})
		return diags
	}

	if updateHandler != nil {
		// NOTE: Must always check of the allowed batch of updates is allowed or not.
		if err = updateHandler.IsUpdateBatchAllowed(d); err != nil {
			slog.Error(err.Error())
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error in DB cluster update handler",
			})
			return diags
		}

		// The allowed batch of updates is Good. Therefore, it is a GO for update. Do it...
		if err = updateHandler.HandleUpdate(newCtx, d, m, clusterInfo); err != nil {
			slog.Error(err.Error())
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error in DB cluster update handler",
			})
			return diags
		}
	}

	d.Set(TF_FIELD_LAST_UPDATED, time.Now().Format(time.RFC822))

	return resourceReadDbCluster(ctx, d, m)
}

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
	jobData.SetRemoveBackups(true)

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

	d.Set(TF_FIELD_LAST_UPDATED, time.Now().Format(time.RFC822))
	d.SetId("")

	return diags
}
