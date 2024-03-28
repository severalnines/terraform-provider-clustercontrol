package provider

import (
	"github.com/severalnines/clustercontrol-client-sdk/go/pkg/openapi"
	"log/slog"
)

type HAProxy struct {
	Common LBCommon
}

func (m *HAProxy) GetInputs(d map[string]any, jobData *openapi.JobsJobJobSpecJobData) error {
	funcName := "HAProxy::GetInputs"
	slog.Debug(funcName)

	var err error

	if err = m.Common.GetInputs(d, jobData); err != nil {
		return err
	}

	jobData.SetAction(JOB_ACTION_SETUP_HAPROXY)
	jobData.SetBuildFromSource(false)

	// TODO - CMON API needs to reconcile the differences in format between PorxySQL and HAProxy for the `node_addresses` field

	return nil
}

func NewHAProxy() *HAProxy {
	return &HAProxy{}
}
