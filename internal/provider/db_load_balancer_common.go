package provider

import (
	"github.com/severalnines/clustercontrol-client-sdk/go/pkg/openapi"
	"log/slog"
)

type LBCommon struct {
}

func (m *LBCommon) GetInputs(d map[string]any, jobData *openapi.JobsJobJobSpecJobData) error {
	funcName := "ProxySql::GetInputs"
	slog.Debug(funcName)

	disableFirewall := d[TF_FIELD_CLUSTER_DISABLE_FW].(bool)
	jobData.SetDisableFirewall(disableFirewall)

	disableSelinux := d[TF_FIELD_CLUSTER_DISABLE_SELINUX].(bool)
	jobData.SetDisableSelinux(disableSelinux)

	installSoftware := d[TF_FIELD_LB_INSTALL_SW].(bool)
	jobData.SetInstallSoftware(installSoftware)

	uninstallSoftware := d[TF_FIELD_LB_ENABLE_UNINSTALL].(bool)
	jobData.SetEnableUninstall(uninstallSoftware)

	sshUser := d[TF_FIELD_CLUSTER_SSH_USER].(string)
	jobData.SetSshUser(sshUser)

	sshUserPassword := d[TF_FIELD_CLUSTER_SSH_PW].(string)
	jobData.SetSudoPassword(sshUserPassword)

	sshKeyFile := d[TF_FIELD_CLUSTER_SSH_KEY_FILE].(string)
	jobData.SetSshKeyfile(sshKeyFile)

	sshPort := d[TF_FIELD_CLUSTER_SSH_PORT].(string)
	jobData.SetSshPort(sshPort)

	return nil
}
