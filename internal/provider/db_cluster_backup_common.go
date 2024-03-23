package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/severalnines/clustercontrol-client-sdk/go/pkg/openapi"
	"log/slog"
)

type DbBackupCommon struct{}

func (c *DbBackupCommon) GetBackupInputs(d *schema.ResourceData, jobData *openapi.JobsJobJobSpecJobData) error {
	funcName := "DbBackupCommon::GetBackupInputs"
	slog.Info(funcName)

	var err error

	backupMethod := d.Get(TF_FIELD_BACKUP_METHOD).(string)
	jobData.SetBackupMethod(backupMethod)

	backupDir := d.Get(TF_FIELD_BACKUP_DIR).(string)
	jobData.SetBackupDir(backupDir)

	backupSubdir := d.Get(TF_FIELD_BACKUP_SUBDIR).(string)
	jobData.SetBackupsubdir(backupSubdir)

	isEncrypt := d.Get(TF_FIELD_BACKUP_ENCRYPT).(bool)
	jobData.SetEncryptBackup(isEncrypt)

	backupHost := d.Get(TF_FIELD_BACKUP_HOST).(string)
	jobData.SetHostname(backupHost)

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
