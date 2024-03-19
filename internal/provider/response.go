package provider

import "github.com/severalnines/clustercontrol-client-sdk/go/pkg/openapi"

type MinResponseFields struct {
	Request_Status string
	Debug_Messages []string
}

type MinJobSpecificResponseFields struct {
	Job_Id      int32
	Status      string
	Status_Text string
}

type JobResponseFields struct {
	TopLevel MinResponseFields
	Job      MinJobSpecificResponseFields
}

type ClusterResponseFields struct {
	TopLevel MinResponseFields
	Cluster  openapi.ClusterResponse
}

//func NewClusterResponse() (*Get)
