package provider

type JobJson struct {
	Job_Id      int32
	Status      string
	Status_Text string
}

type ClusterInfo struct {
	Cluster_Id int32
	Tags       []string
}

type ResponseJobJson struct {
	Request_Status string
	Debug_Messages []string
	Job            JobJson
}

type ClusterInfoRespJson struct {
	Request_Status string
	Cluster        ClusterInfo
}
