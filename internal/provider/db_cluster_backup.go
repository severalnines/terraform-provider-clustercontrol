package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/severalnines/clustercontrol-client-sdk/go/pkg/openapi"
	"log/slog"
	"strconv"
	"strings"
	"time"
)

type DbClusterBackupInterface interface {
	GetBackupInputs(d *schema.ResourceData, jobData *openapi.JobsJobJobSpecJobData) error
}

func resourceDbClusterBackup() *schema.Resource {
	funcName := "resourceDbLoadBalancer"
	slog.Debug(funcName)

	return &schema.Resource{
		CreateContext: resourceCreateDbClusterBackup,
		ReadContext:   resourceReadDbClusterBackup,
		UpdateContext: resourceUpdateDbClusterBackup,
		DeleteContext: resourceDeleteDbClusterBackup,
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
			TF_FIELD_CLUSTER_ID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The database cluster ID for which this LB is being deployed to.",
			},
			TF_FIELD_BACKUP_METHOD: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "mariabackup, xtrabackup, ...",
			},
			TF_FIELD_BACKUP_DIR: {
				Type:        schema.TypeString,
				Required:    true,
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
				Required:    true,
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
		},
	}
}

// Prem
func resourceCreateDbClusterBackup(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	funcName := "resourceCreateDbLoadBalancer"
	slog.Debug(funcName)

	var diags diag.Diagnostics
	var err error

	newCtx := context.WithValue(ctx, "cookie", gNewCtx.Value("cookie"))

	apiClient := m.(*openapi.APIClient)

	createBackup := NewCCJob(CMON_JOB_CREATE_JOB)
	job := createBackup.GetJob()
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
	createBackup.SetClusterId(clusterInfo.GetClusterId())

	clusterType := clusterInfo.GetClusterType()
	slog.Debug(funcName, "ClusterType", clusterType)
	clusterType = strings.ToLower(clusterType)

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
	case CLUSTER_TYPE_MSSQL_SINGLE:
		backupHandler = NewMsSql()
	case CLUSTER_TYPE_MSSQL_AO_ASYNC:
		backupHandler = NewMsSql()
	case CLUSTER_TYPE_ELASTIC:
		backupHandler = NewElastic()
	default:
		slog.Warn(funcName, "Unknown cluster type", clusterType)
	}

	if backupHandler != nil {
		if err = backupHandler.GetBackupInputs(d, &jobData); err != nil {
			slog.Error(err.Error())
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error in DB cluster Read handler",
			})
			return diags
		}
	}

	jobSpec.SetJobData(jobData)
	job.SetJobSpec(jobSpec)
	createBackup.SetJob(job)

	var jobId int32
	if jobId, err = SendAndWaitForJobCompletionAndId(newCtx, apiClient, createBackup); err != nil {
		slog.Error(err.Error())
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Job Failed for ClusterCreate",
		})
		return diags
	}
	slog.Debug(funcName, "Job completed", jobId)

	var backupId int32
	if backupId, err = GetBackupIdForCluster(newCtx, apiClient, clusterInfo.GetClusterId(), jobId); err != nil {
		slog.Error(err.Error())
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Job Failed for ClusterCreate",
		})
		return diags
	}

	backupIdStr := strconv.Itoa(int(backupId))
	d.SetId(backupIdStr)
	d.Set(TF_FIELD_RESOURCE_ID, backupIdStr)
	d.Set(TF_FIELD_LAST_UPDATED, time.Now().Format(time.RFC850))

	return diags
}

// Prem
func resourceDeleteDbClusterBackup(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	funcName := "resourceDeleteDbCluster"
	slog.Debug(funcName)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	var err error

	newCtx := context.WithValue(ctx, "cookie", gNewCtx.Value("cookie"))

	apiClient := m.(*openapi.APIClient)

	deleteBackup := NewCCJob(CMON_JOB_CREATE_JOB)
	job := deleteBackup.GetJob()
	jobSpec := job.GetJobSpec()
	jobSpec.SetCommand(CMON_JOB_DELETE_BACKUP_COMMAND)
	jobData := jobSpec.GetJobData()

	backupId := d.Id()
	clusterId := d.Get(TF_FIELD_CLUSTER_ID).(string)
	slog.Info(funcName, "Deleting backup-Id", backupId, "cluster-Id", clusterId)

	var clusterInfo *openapi.ClusterResponse
	if clusterInfo, err = GetClusterByClusterStrId(newCtx, apiClient, clusterId); err != nil {
		slog.Error(err.Error())
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  err.Error(),
		})
		return diags
	}
	deleteBackup.SetClusterId(clusterInfo.GetClusterId())

	var iBkpId int
	if iBkpId, err = strconv.Atoi(backupId); err != nil {
		slog.Error(err.Error())
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  err.Error(),
		})
		return diags
	}
	jobData.SetBackupid(int32(iBkpId))
	jobData.SetClusterid(clusterInfo.GetClusterId())

	jobSpec.SetJobData(jobData)
	job.SetJobSpec(jobSpec)
	deleteBackup.SetJob(job)

	if err = SendAndWaitForJobCompletion(newCtx, apiClient, deleteBackup); err != nil {
		slog.Error(err.Error())
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Job Failed for Delete Backup",
		})
		return diags
	}

	d.SetId("")
	d.Set(TF_FIELD_LAST_UPDATED, time.Now().Format(time.RFC822))

	return diags
}

// Prem
func resourceReadDbClusterBackup(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	funcName := "resourceReadDbClusterBackup"
	slog.Debug(funcName)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	//var err error

	return diags
}

// Prem
func resourceUpdateDbClusterBackup(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	funcName := "resourceUpdateDbClusterBackup"
	slog.Debug(funcName)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	//var err error

	return diags
}
