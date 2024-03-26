package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/severalnines/clustercontrol-client-sdk/go/pkg/openapi"
	"log/slog"
	"strconv"
)

type LBCommon struct {
}

func (m *LBCommon) GetInputs(d *schema.ResourceData, jobData *openapi.JobsJobJobSpecJobData) error {
	funcName := "ProxySql::GetInputs"
	slog.Debug(funcName)

	var err error

	clusterId := d.Get(TF_FIELD_CLUSTER_ID).(string)
	var cId int
	if cId, err = strconv.Atoi(clusterId); err != nil {
		slog.Error(funcName, "Non-numeric clusterId", clusterId)
		return err
	}
	jobData.SetClusterid(int32(cId))

	disableFirewall := d.Get(TF_FIELD_CLUSTER_DISABLE_FW).(bool)
	jobData.SetDisableFirewall(disableFirewall)

	disableSelinux := d.Get(TF_FIELD_CLUSTER_DISABLE_SELINUX).(bool)
	jobData.SetDisableSelinux(disableSelinux)

	installSoftware := d.Get(TF_FIELD_LB_INSTALL_SW).(bool)
	jobData.SetInstallSoftware(installSoftware)

	uninstallSoftware := d.Get(TF_FIELD_LB_ENABLE_UNINSTALL).(bool)
	jobData.SetEnableUninstall(uninstallSoftware)

	//jobData.SetGenerateToken(true)

	sshUser := d.Get(TF_FIELD_CLUSTER_SSH_USER).(string)
	jobData.SetSshUser(sshUser)

	sshUserPassword := d.Get(TF_FIELD_CLUSTER_SSH_PW).(string)
	jobData.SetSudoPassword(sshUserPassword)

	sshKeyFile := d.Get(TF_FIELD_CLUSTER_SSH_KEY_FILE).(string)
	jobData.SetSshKeyfile(sshKeyFile)

	sshPort := d.Get(TF_FIELD_CLUSTER_SSH_PORT).(string)
	jobData.SetSshPort(sshPort)

	return nil
}
