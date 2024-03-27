package provider

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/severalnines/clustercontrol-client-sdk/go/pkg/openapi"
	"log/slog"
	"slices"
)

type DbBackupCommon struct{}

func (c *DbBackupCommon) GetBackupInputs(d *schema.ResourceData, jobData *openapi.JobsJobJobSpecJobData) error {
	funcName := "DbBackupCommon::GetBackupInputs"
	slog.Info(funcName)

	var err error

	// NOTE: backup_method is/should-be "" for Redis
	backupMethod := d.Get(TF_FIELD_BACKUP_METHOD).(string)
	if backupMethod != "" {
		jobData.SetBackupMethod(backupMethod)
	}

	backupDir := d.Get(TF_FIELD_BACKUP_DIR).(string)
	if backupMethod != "" {
		jobData.SetBackupDir(backupDir)
	}

	backupSubdir := d.Get(TF_FIELD_BACKUP_SUBDIR).(string)
	if backupSubdir != "" {
		jobData.SetBackupsubdir(backupSubdir)
	}

	isEncrypt := d.Get(TF_FIELD_BACKUP_ENCRYPT).(bool)
	jobData.SetEncryptBackup(isEncrypt)

	backupHost := d.Get(TF_FIELD_BACKUP_HOST).(string)
	if backupHost != "" {
		jobData.SetHostname(backupHost)
	}

	isCompressBackup := d.Get(TF_FIELD_BACKUP_COMPRESSION).(bool)
	jobData.SetCompression(isCompressBackup)

	compressLevel := d.Get(TF_FIELD_BACKUP_COMPRESSION_LEVEL).(int)
	jobData.SetCompressionLevel(int32(compressLevel))

	retention := d.Get(TF_FIELD_BACKUP_RETENTION).(int)
	jobData.SetBackupRetention(int32(retention))

	isStoreOnCtlr := d.Get(TF_FIELD_BACKUP_ON_CONTROLLER).(bool)
	jobData.SetCcStorage(isStoreOnCtlr)

	return err
}

func (c *DbBackupCommon) IsValidBackupOptions(vendor string, clusterType string, jobData *openapi.JobsJobJobSpecJobData) error {

	clusterTypeMap, ok := gDbAvailableBackupMethods[clusterType]
	if !ok {
		return errors.New(fmt.Sprintf("Backup method map doesn't support DB cluster-type: %s", clusterTypeMap))
	}
	vendorMap, ok := clusterTypeMap[vendor]
	if !ok {
		return errors.New(fmt.Sprintf("Backup method map doesn't support DB vendor: %s, ClusterType: %s", vendor, clusterType))
	}
	if !slices.Contains(vendorMap, jobData.GetBackupMethod()) {
		return errors.New(fmt.Sprintf("Backup method map doesn't support DB vendor: %s, ClusterType: %s, BackupMethod: %s", vendor, clusterType, jobData.GetBackupMethod()))
	}

	return nil
}
