package provider

import (
	"errors"
	"fmt"
	"github.com/severalnines/clustercontrol-client-sdk/go/pkg/openapi"
	"log/slog"
	"strconv"
)

type ProxySql struct {
	Common LBCommon
}

func (m *ProxySql) GetInputs(d map[string]any, jobData *openapi.JobsJobJobSpecJobData) error {
	funcName := "ProxySql::GetInputs"
	slog.Debug(funcName)

	var err error

	if err = m.Common.GetInputs(d, jobData); err != nil {
		return err
	}

	jobData.SetAction(JOB_ACTION_SETUP_PROXYSQL)
	jobData.SetDbDatabase("*.*")

	lbVersion := d[TF_FIELD_LB_VERSION].(string)
	jobData.SetVersion(lbVersion)

	adminUsr := d[TF_FIELD_LB_ADMIN_USER].(string)
	jobData.SetAdminUser(adminUsr)
	adminPw := d[TF_FIELD_LB_ADMIN_USER_PW].(string)
	jobData.SetAdminPassword(adminPw)

	monitorUsr := d[TF_FIELD_LB_MONITOR_USER].(string)
	jobData.SetMonitorUser(monitorUsr)
	monitorPw := d[TF_FIELD_LB_MONITOR_USER_PW].(string)
	jobData.SetMonitorPassword(monitorPw)

	port := d[TF_FIELD_LB_PORT].(string)
	if port == "" {
		port = DEFAULT_PROXYSQL_LISTEN_PORT
	}
	var iPort int
	if iPort, err = strconv.Atoi(port); err != nil {
		iPort, _ = strconv.Atoi(DEFAULT_PROXYSQL_LISTEN_PORT)
	}
	jobData.SetPort(int32(iPort))

	useClustering := d[TF_FIELD_LB_USE_CLUSTERING].(bool)
	jobData.SetUseClustering(useClustering)

	useRWsplitting := d[TF_FIELD_LB_USE_RW_SPLITTING].(bool)
	jobData.SetUseRwSplit(useRWsplitting)

	host := d[TF_FIELD_LB_MY_HOST]
	isAtleastOneNodeDeclared := false
	for _, ff := range host.([]any) {
		f := ff.(map[string]any)
		hostname := f[TF_FIELD_CLUSTER_HOSTNAME].(string)
		myAdminPort := f[TF_FIELD_CLUSTER_HOST_PORT].(string)
		if iPort, err = strconv.Atoi(myAdminPort); err != nil {
			iPort, _ = strconv.Atoi(DEFAULT_PROXYSQL_ADMIN_PORT)
		}
		var node = openapi.JobsJobJobSpecJobDataNode{
			Hostname: &hostname,
		}
		node.SetPort(int32(iPort))
		jobData.SetNode(node)

		isAtleastOneNodeDeclared = true

		// Support only one node (i.e., self)
		break
	}

	if !isAtleastOneNodeDeclared {
		err = errors.New(fmt.Sprintf("ERROR: At lease one %s block must be specified", TF_FIELD_LB_MY_HOST))
		return err
	}

	return nil
}

func NewProxySql() *ProxySql {
	return &ProxySql{}
}
