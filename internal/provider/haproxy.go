package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/severalnines/clustercontrol-client-sdk/go/pkg/openapi"
	"log/slog"
)

type HAProxy struct {
	Common LBCommon
}

func (m *HAProxy) GetInputs(d *schema.ResourceData, jobData *openapi.JobsJobJobSpecJobData) error {
	funcName := "HAProxy::GetInputs"
	slog.Debug(funcName)

	if err := m.Common.GetInputs(d, jobData); err != nil {
		return err
	}

	//clusterId := d.Get(TF_FIELD_CLUSTER_ID)

	return nil
}

func NewHAProxy() *HAProxy {
	return &HAProxy{}
}
