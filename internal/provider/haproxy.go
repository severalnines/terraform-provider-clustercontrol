package provider

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/severalnines/clustercontrol-client-sdk/go/pkg/openapi"
	"log/slog"
	"strconv"
)

type HAProxy struct {
	Common LBCommon
}

func (m *HAProxy) GetInputs(d *schema.ResourceData, jobData *openapi.JobsJobJobSpecJobData) error {
	funcName := "HAProxy::GetInputs"
	slog.Debug(funcName)

	var err error

	if err = m.Common.GetInputs(d, jobData); err != nil {
		return err
	}

	jobData.SetAction(JOB_ACTION_SETUP_HAPROXY)
	jobData.SetBuildFromSource(false)
	//jobData.SetNode

	var iPort int
	hosts := d.Get(TF_FIELD_CLUSTER_HOST)
	nodeAddresses := []openapi.JobsJobJobSpecJobDataNodeAdressesInner{}
	isAtleastOneNodeAddressDeclared := false
	for _, ff := range hosts.([]any) {
		f := ff.(map[string]any)
		hostname := f[TF_FIELD_CLUSTER_HOSTNAME].(string)
		port := f[TF_FIELD_CLUSTER_HOST_PORT].(string)
		if port == "" {
			port = DEFAULT_MYSQL_PORT
		}
		if iPort, err = strconv.Atoi(port); err != nil {
			iPort, _ = strconv.Atoi(DEFAULT_MYSQL_PORT)
		}
		var node = openapi.JobsJobJobSpecJobDataNodeAdressesInner{
			Hostname: &hostname,
		}
		node.SetPort(int32(iPort))
		nodeAddresses = append(nodeAddresses, node)
		isAtleastOneNodeAddressDeclared = true

	}
	jobData.SetNodeAdresses(nodeAddresses)
	if !isAtleastOneNodeAddressDeclared {
		err = errors.New(fmt.Sprintf("ERROR: At lease one %s block must be specified", TF_FIELD_CLUSTER_HOST))
		return err
	}

	host := d.Get(TF_FIELD_LB_MY_HOST)
	isAtleastOneNodeDeclared := false
	var node = openapi.JobsJobJobSpecJobDataNode{}
	for _, ff := range host.([]any) {
		f := ff.(map[string]any)
		hostname := f[TF_FIELD_CLUSTER_HOSTNAME].(string)
		node.SetHostname(hostname)

		// Must set hostname even thought it is not used. This is so that the caller of this method
		// can receive it and used it to set the resourece-ID
		jobData.SetHostname(hostname)

		isAtleastOneNodeDeclared = true

		// Support only one node (i.e., self)
		break
	}
	if !isAtleastOneNodeDeclared {
		err = errors.New(fmt.Sprintf("ERROR: At lease one %s block must be specified", TF_FIELD_LB_MY_HOST))
		return err
	}
	//node.SetXXX
	jobData.SetNode(node)

	return nil
}

func NewHAProxy() *HAProxy {
	return &HAProxy{}
}
