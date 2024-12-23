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
	"strconv"
	"strings"
	"time"
)

func resourceDbClusterBackupSchedule() *schema.Resource {
	funcName := "resourceDbClusterBackupSchedule"
	slog.Debug(funcName)

	return &schema.Resource{
		CreateContext: resourceCreateDbClusterBackupSched,
		ReadContext:   resourceReadDbClusterBackupSched,
		UpdateContext: resourceUpdateDbClusterBackupSched,
		DeleteContext: resourceDeleteDbClusterBackupSched,
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
			TF_FIELD_BACKUP_SCHED_TITLE: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A title for the backup schedule (e.g., Daily full, Hourly incremental, etc)",
			},
			TF_FIELD_BACKUP_SCHED_TIME: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The time to kick off a backup (e.g. 'TZ=UTC 0 0 * * *')",
			},
			TF_FIELD_BACKUP_METHOD: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "mariabackup, xtrabackup, ...",
			},
			TF_FIELD_BACKUP_DIR: {
				Type:        schema.TypeString,
				Optional:    true, // PBM for mongo doesn't require a backupdir ...
				Description: "Base direcory where backups is to be stored",
			},
			TF_FIELD_BACKUP_SUBDIR: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sub-dir for backups - default: \"BACKUP-%I\" ",
			},
			TF_FIELD_BACKUP_ENCRYPT: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to encrypt or not",
			},
			TF_FIELD_BACKUP_HOST: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Where there are multiple hosts, which host to choose to create backup from.",
			},
			TF_FIELD_BACKUP_COMPRESSION: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to compress backups or not",
			},
			TF_FIELD_BACKUP_COMPRESSION_LEVEL: {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Compression level",
			},
			TF_FIELD_BACKUP_RETENTION: {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Backup retention period in days",
			},
			TF_FIELD_BACKUP_ON_CONTROLLER: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to store the backup on CMON controller host or not",
			},
			TF_FIELD_BACKUP_FAILOVER: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to enable backup failover to another host in case the host crashes",
			},
			TF_FIELD_BACKUP_FAILOVER_HOST: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "If backup failover is enabled, which host to use as backup host in the event of failure of first choice host.",
			},
			TF_FIELD_BACKUP_STORAGE_HOST: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Which host to store the backup on. Typically, used with mongodump backup method.",
			},
			TF_FIELD_BACKUP_SYSTEM_DB: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to enable backup failover to another host in case the host crashes",
			},
			TF_FIELD_CLUSTER_SNAPSHOT_REPO: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Elasticsearch snapshot repository",
			},
		},
	}
}

func resourceCreateDbClusterBackupSched(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	funcName := "resourceCreateDbClusterBackupSched"
	slog.Debug(funcName)

	var diags diag.Diagnostics
	var err error

	providerDetails := m.(*ProviderDetails)

	//newCtx := context.WithValue(ctx, "cookie", gNewCtx.Value("cookie"))
	//newCtx := context.WithValue(ctx, "cookie", providerDetails.SessionIdCtx.Value("cookie"))
	newCtx := context.WithValue(ctx, "cookie", providerDetails.SessionCookie)

	//apiClient := m.(*openapi.APIClient)
	apiClient := providerDetails.ApiClient

	createBackupSched := NewCCJob(CMON_JOB_CREATE_JOB)
	job := createBackupSched.GetJob()
	job.SetStatus(JOB_STATUS_SCHEDULED)
	jobSpec := job.GetJobSpec()
	jobSpec.SetCommand(CMON_JOB_CREATE_BACKUP_COMMAND)
	jobData := jobSpec.GetJobData()

	clusterId := d.Get(TF_FIELD_CLUSTER_ID).(string)
	slog.Debug(funcName, "ClusterId", clusterId)

	var clusterInfo *openapi.ClusterResponse
	if clusterInfo, err = GetClusterByClusterStrId(newCtx, apiClient, clusterId); err != nil {
		slog.Error(err.Error())
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  err.Error(),
		})
		return diags
	}
	createBackupSched.SetClusterId(clusterInfo.GetClusterId())

	title := d.Get(TF_FIELD_BACKUP_SCHED_TITLE).(string)
	if title != "" {
		job.SetTitle(title)
	}
	sched := d.Get(TF_FIELD_BACKUP_SCHED_TIME).(string)
	// sched will never be empty string as it is a required param
	job.SetRecurrence(sched)

	clusterType := clusterInfo.GetClusterType()
	vendor := clusterInfo.GetVendor()
	slog.Debug(funcName, "ClusterType", clusterType)
	clusterType = strings.ToLower(clusterType)
	vendor = strings.ToLower(vendor)

	var backupHandler DbClusterBackupInterface = nil

	switch clusterType {
	case CLUSTER_TYPE_REPLICATION:
		backupHandler = NewMySQLMaria()
	case CLUSTER_TYPE_GALERA:
		backupHandler = NewMySQLMaria()
	case CLUSTER_TYPE_PG_SINGLE:
		backupHandler = NewPostgres()
	case CLUSTER_TYPE_MOGNODB:
		backupHandler = NewMongo()
	case CLUSTER_TYPE_REDIS:
		backupHandler = NewRedis()
	case CLUSTER_TYPE_VALKEY:
		backupHandler = NewRedis()
	case CLUSTER_TYPE_REDIS_SHARDED:
		backupHandler = NewRedisSharded()
	case CLUSTER_TYPE_VALKEY_SHARDED:
		backupHandler = NewRedisSharded()
	case CLUSTER_TYPE_MSSQL_SINGLE:
		backupHandler = NewMsSql()
	case CLUSTER_TYPE_MSSQL_AO_ASYNC:
		backupHandler = NewMsSql()
	case CLUSTER_TYPE_ELASTIC:
		backupHandler = NewElastic()
	default:
		str := fmt.Sprintf("%s - Unknown cluster type: %s", funcName, clusterType)
		slog.Warn(str)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  str,
		})
		return diags
	}

	if err = backupHandler.GetBackupInputs(d, &jobData); err != nil {
		slog.Error(err.Error())
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  err.Error(),
		})
		return diags
	}

	if err = backupHandler.IsValidBackupOptions(vendor, clusterType, &jobData); err != nil {
		slog.Error(err.Error())
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  err.Error(),
		})
		return diags
	}

	_ = backupHandler.SetBackupJobData(&jobData)

	jobSpec.SetJobData(jobData)
	job.SetJobSpec(jobSpec)
	createBackupSched.SetJob(job)

	var resp *http.Response
	if resp, err = apiClient.JobsAPI.JobsPost(newCtx).Jobs(*createBackupSched).Execute(); err != nil {
		slog.Error(err.Error())
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  err.Error(),
		})
		return diags
	}
	slog.Debug(funcName, "Resp `Job`", resp)

	var respBytes []byte
	if respBytes, err = io.ReadAll(resp.Body); err != nil {
		slog.Error(err.Error())
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  err.Error(),
		})
		return diags
	}

	var jobResp JobResponseFields
	if err = json.Unmarshal(respBytes, &jobResp); err != nil {
		slog.Error(err.Error())
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  err.Error(),
		})
		return diags
	}
	slog.Debug(funcName, "Job completed", jobResp.Job.Job_Id)

	backupIdStr := strconv.Itoa(int(jobResp.Job.Job_Id))
	d.SetId(backupIdStr)
	d.Set(TF_FIELD_RESOURCE_ID, backupIdStr)
	d.Set(TF_FIELD_LAST_UPDATED, time.Now().Format(time.RFC850))

	return diags
}

func resourceReadDbClusterBackupSched(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	funcName := "resourceReadDbClusterBackupSched"
	slog.Debug(funcName)

	var diags diag.Diagnostics
	//var err error

	return diags
}

func resourceUpdateDbClusterBackupSched(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	funcName := "resourceUpdateDbClusterBackupSched"
	slog.Debug(funcName)

	var diags diag.Diagnostics
	//var err error

	return diags
}

func resourceDeleteDbClusterBackupSched(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	funcName := "resourceDeleteDbClusterBackupSched"
	slog.Debug(funcName)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	var err error

	providerDetails := m.(*ProviderDetails)

	//newCtx := context.WithValue(ctx, "cookie", gNewCtx.Value("cookie"))
	//newCtx := context.WithValue(ctx, "cookie", providerDetails.SessionIdCtx.Value("cookie"))
	newCtx := context.WithValue(ctx, "cookie", providerDetails.SessionCookie)

	//apiClient := m.(*openapi.APIClient)
	apiClient := providerDetails.ApiClient

	deleteBackupSched := NewCCJob(CMON_JOB_DELETE_JOB)
	backupSchedId := d.Id()
	var schedId int
	if schedId, err = strconv.Atoi(backupSchedId); err != nil {
		slog.Error(err.Error())
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  err.Error(),
		})
		return diags
	}
	deleteBackupSched.SetJobId(int32(schedId))

	if _, err = apiClient.JobsAPI.JobsPost(newCtx).Jobs(*deleteBackupSched).Execute(); err != nil {
		slog.Error(err.Error())
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  err.Error(),
		})
		return diags
	}
	slog.Info(funcName, "Backup schedule successfully deleted", backupSchedId)

	return diags
}
