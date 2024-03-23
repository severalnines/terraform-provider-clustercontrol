package provider

import "github.com/severalnines/clustercontrol-client-sdk/go/pkg/openapi"

type MinResponseFields struct {
	Controller_Id  string
	Request_Status string
	Debug_Messages []string
}

type MinJobSpecificResponseFields struct {
	Job_Id      int32
	Status      string
	Status_Text string
}

type JobResponseFields struct {
	Request_Status string
	Job            MinJobSpecificResponseFields
}

type ClusterResponseFields struct {
	Request_Status string
	Cluster        openapi.ClusterResponse
}

type MaintenanceOperationResponse struct {
	Request_Status string
	UUID           string
}

type MaintenanceGetOperationResponse struct {
	Request_Status      string
	Total               int32
	Maintenance_Records openapi.MaintenanceResponse
}

type BackupGetOperationResponseTop struct {
	Controller_Id  string
	Request_Status string
	Total          int32
}

type BackupGetOperationResponse struct {
	Request_Status string
	Total          int32
	Backup_Records []openapi.BackupResponseBackupRecordsInner
}

//func NewClusterResponse() (*Get)
