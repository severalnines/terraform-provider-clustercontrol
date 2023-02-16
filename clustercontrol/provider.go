package clustercontrol

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/severalnines/clustercontrol-client-sdk/go/pkg/openapi"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	gNewCtx    context.Context
	gCfg       *openapi.Configuration
	gApiClient *openapi.APIClient
)

const (
	API_USER    = "user"
	API_USER_PW = "password"
	CNTLR_URL   = "controller_url"
)

const (
	RESOURCE_CLUSTER = "cc_mysql_maria_cluster"
)

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

func Provider() *schema.Provider {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			API_USER: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc(API_USER, nil),
			},
			API_USER_PW: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc(API_USER_PW, nil),
			},
			CNTLR_URL: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			RESOURCE_CLUSTER: resourceMySqlMariaCluster(),
		},
		DataSourcesMap:       map[string]*schema.Resource{},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	username := d.Get(API_USER).(string)
	password := d.Get(API_USER_PW).(string)
	url := d.Get(CNTLR_URL).(string)

	gCfg = newConfiguration(url)
	gApiClient = openapi.NewAPIClient(gCfg)
	authenticate := *openapi.NewAuthenticate("authenticateWithPassword")
	authenticate.SetUserName(username /*os.Getenv("API_USER")*/)
	authenticate.SetPassword(password /*os.Getenv("API_USER_PW")*/)
	// ctx := context.Background()

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	resp, err := gApiClient.AuthApi.AuthPost(ctx).Authenticate(authenticate).Execute()
	if err != nil {
		printError(err, resp)
		return nil, diag.FromErr(err)
	}
	fmt.Fprintf(os.Stderr, "Resp `AuthApi.AuthPost`: %v\n", resp)

	// fmt.Println("#Cookies: ", len(resp.Cookies()))
	for _, cookie := range resp.Cookies() {
		// fmt.Fprintf(os.Stderr, "Found cookie %v\n", cookie)
		gNewCtx = context.WithValue(ctx, "cookie", cookie)
		break
	}

	return gApiClient, diags
}

func resourceMySqlMariaCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreateMySqlMariaCluster,
		ReadContext:   resourceReadMySqlMariaCluster,
		UpdateContext: resourceUpdateMySqlMariaCluster,
		DeleteContext: resourceDeleteCluster,
		Importer:      &schema.ResourceImporter{},
		Schema: map[string]*schema.Schema{
			"cluster_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the resource, also acts as it's unique ID",
				ForceNew:    true,
				// ValidateFunc: validateName,
			},
			"database_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A description of an item",
			},
			"database_vendor": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A description of an item",
			},
			"database_version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A description of an item",
			},
			"database_topology": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A description of an item",
			},
			"primary_database_host": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A description of an item",
			},
			"hostname_internal": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The internal hostnames.",
			},
			"tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "An optional list of tags, represented as a key, value pair.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"ssh_key_file": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "SSH Key file.",
			},
			"ssh_user": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "SSH user.",
			},
			"install_software": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "If CMON should install db software or not.",
			},
		},
	}
}

func resourceCreateMySqlMariaCluster(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	newCtx := context.WithValue(ctx, "cookie", gNewCtx.Value("cookie"))

	apiClient := m.(*openapi.APIClient)

	tfTags := d.Get("tags").(*schema.Set).List()
	tags := make([]string, len(tfTags))
	for i, tfTag := range tfTags {
		tags[i] = tfTag.(string)
	}

	clusterName := d.Get("cluster_name").(string)
	createCluster := *openapi.NewJobs("createJobInstance")
	sshKey := d.Get("ssh_key_file").(string)
	sshUser := d.Get("ssh_user").(string)
	installSoftwareStr := d.Get("install_software").(string)
	installSoftware := true
	if installSoftwareStr != "" {
		installSoftware, _ = strconv.ParseBool(installSoftwareStr)
	}
	job := *openapi.NewJobsJob()
	job.SetClassName("CmonJobInstance")

	jobSpec := *openapi.NewJobsJobJobSpec()
	jobSpec.SetCommand("create_cluster")

	jobData := *openapi.NewJobsJobJobSpecJobData()
	jobData.SetClusterName(clusterName)
	jobData.SetClusterType(d.Get("database_topology").(string))
	jobData.SetConfigTemplate("my.cnf.repl80")
	jobData.SetDataDir("/var/lib/mysql")
	jobData.SetDbPassword("pA0d7HJuTAb1YDAJSmGD")
	jobData.SetDisableFirewall(true)
	jobData.SetDisableSelinux(true)
	jobData.SetGenerateToken(true)
	jobData.SetInstallSoftware(installSoftware)
	jobData.SetMysqlSemiSync(true)
	jobData.SetPort(3306)
	jobData.SetSshKeyfile(sshKey)
	jobData.SetSshPort("22")
	jobData.SetSshUser(sshUser)
	// jobData.SetUserId(5)
	jobData.SetVendor(d.Get("database_vendor").(string))
	jobData.SetVersion(d.Get("database_version").(string))

	hostname := d.Get("primary_database_host").(string)
	hostname_internal := d.Get("hostname_internal").(string)
	port := 3306
	nodes := generateNodeList(hostname, hostname_internal, port)

	jobData.SetNodes(nodes)
	jobSpec.SetJobData(jobData)
	job.SetJobSpec(jobSpec)
	createCluster.SetJob(job)
	resp, err := apiClient.JobsApi.JobsPost(newCtx).Jobs(createCluster).Execute()
	if err != nil {
		printError(err, resp)
		return diag.FromErr(err)
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		// log.Fatal(err)
		printError(err, resp)
		return diag.FromErr(err)
	}

	var jobResp ResponseJobJson
	err = json.Unmarshal(respBytes, &jobResp)
	if err != nil {
		// log.Fatal(err)
		printError(err, resp)
		return diag.FromErr(err)
	}
	fmt.Fprintf(os.Stderr, "Resp `Job`: %v\n", jobResp)

	// Wait for job to complete
	isCreateSuccess := true
	for true {
		// Calling Sleep method
		time.Sleep(5 * time.Second)

		checkJobStatus := *openapi.NewJobs("getJobInstance")
		//job := *openapi.NewJobsJob()
		//job.SetClassName("CmonJobInstance")
		checkJobStatus.SetJobId(jobResp.Job.Job_Id)
		resp, err = apiClient.JobsApi.JobsPost(newCtx).Jobs(checkJobStatus).Execute()
		if err != nil {
			printError(err, resp)
			return diag.FromErr(err)
		}

		respBytes, err = io.ReadAll(resp.Body)
		if err != nil {
			// log.Fatal(err)
			printError(err, resp)
			return diag.FromErr(err)
		}

		// var jobResp ResponseJobJson
		err = json.Unmarshal(respBytes, &jobResp)
		if err != nil {
			// log.Fatal(err)
			printError(err, resp)
			return diag.FromErr(err)
		}
		fmt.Fprintf(os.Stderr, "Resp `Job`: %v\n", jobResp)

		if jobResp.Job.Status == "FINISHED" {
			break
		}

		if jobResp.Job.Status != "RUNNING" && jobResp.Job.Status != "DEFINED" {
			isCreateSuccess = false
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Resource creation failed",
			})
			break
		}
	}

	if !isCreateSuccess {
		return diags
	}

	/*
	 * Get `cluster_id` to return back to Terraform
	 */
	clusterInfoReq := *openapi.NewClusters("getclusterinfo")
	clusterInfoReq.SetClusterName(clusterName)
	resp, err = apiClient.ClustersApi.ClustersPost(newCtx).Clusters(clusterInfoReq).Execute()
	if err != nil {
		printError(err, resp)
		return diag.FromErr(err)
	}
	fmt.Fprintf(os.Stderr, "Resp `ClustersPost.getallclusterinfo`: %v\n", resp)

	respBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		// log.Fatal(err)
		printError(err, nil)
		return diag.FromErr(err)
	}

	var clusterInfoResp ClusterInfoRespJson
	err = json.Unmarshal(respBytes, &clusterInfoResp)
	if err != nil {
		// log.Fatal(err)
		printError(err, nil)
		return diag.FromErr(err)
	}
	fmt.Fprintf(os.Stderr, "Resp `Job`: %v\n", clusterInfoResp)

	// d.SetId(clusterName)
	d.SetId(strconv.FormatInt(int64(clusterInfoResp.Cluster.Cluster_Id), 10))

	return diags
}

func generateNodeList(hostname string, hostnameInternal string, port int) []openapi.JobsJobJobSpecJobDataNodesInner {
	hostnameArr := strings.Split(hostname, ",")
	hostnameIntArr := strings.Split(hostnameInternal, ",")
	cnt := len(hostnameArr)
	if hostnameInternal != "" && len(hostnameIntArr) > 0 && len(hostnameArr) != len(hostnameIntArr) {
		fmt.Fprintf(os.Stderr,
			"inconsistent number of elements in internal and public hostnames: %d.",
			len(hostnameIntArr))
		return nil
	}
	p := strconv.Itoa(port)
	var nodes = make([]openapi.JobsJobJobSpecJobDataNodesInner, cnt)
	for i := 0; i < cnt; i++ {
		hostnameInt := ""
		if hostnameInternal != "" && len(hostnameIntArr) > 0 {
			hostnameInt = hostnameIntArr[i]
		}
		nodes[i] = openapi.JobsJobJobSpecJobDataNodesInner{
			Hostname:         &hostnameArr[i],
			HostnameData:     &hostnameArr[i],
			HostnameInternal: &hostnameInt,
			Port:             &p,
		}
	}
	return nodes
}

func resourceReadMySqlMariaCluster(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	// apiClient := m.(*openapi.APIClient)

	return diags
}

func resourceUpdateMySqlMariaCluster(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}

func resourceDeleteCluster(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	newCtx := context.WithValue(ctx, "cookie", gNewCtx.Value("cookie"))

	apiClient := m.(*openapi.APIClient)

	delCluster := *openapi.NewJobs("createJobInstance")

	job := *openapi.NewJobsJob()
	job.SetClassName("CmonJobInstance")

	jobSpec := *openapi.NewJobsJobJobSpec()
	jobSpec.SetCommand("remove_cluster")

	jobData := *openapi.NewJobsJobJobSpecJobData()
	clusterId, err := strconv.ParseInt(d.Id(), 10, 32)
	jobData.SetClusterid(int32(clusterId))
	jobSpec.SetJobData(jobData)
	job.SetJobSpec(jobSpec)
	delCluster.SetJob(job)
	resp, err := gApiClient.JobsApi.JobsPost(newCtx).Jobs(delCluster).Execute()
	if err != nil {
		printError(err, resp)
		return diag.FromErr(err)
	}
	fmt.Fprintf(os.Stderr, "Resp `Cluster.Delete`: %v\n", resp)

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		// log.Fatal(err)
		printError(err, resp)
		return diag.FromErr(err)
	}

	var jobResp ResponseJobJson
	err = json.Unmarshal(respBytes, &jobResp)
	if err != nil {
		// log.Fatal(err)
		printError(err, resp)
		return diag.FromErr(err)
	}
	fmt.Fprintf(os.Stderr, "Resp `Job`: %v\n", jobResp)

	// Wait for job to complete
	IsDeleteSuccess := true
	for true {
		// Calling Sleep method
		time.Sleep(5 * time.Second)

		checkJobStatus := *openapi.NewJobs("getJobInstance")
		//job := *openapi.NewJobsJob()
		//job.SetClassName("CmonJobInstance")
		checkJobStatus.SetJobId(jobResp.Job.Job_Id)
		resp, err = apiClient.JobsApi.JobsPost(newCtx).Jobs(checkJobStatus).Execute()
		if err != nil {
			printError(err, resp)
			return diag.FromErr(err)
		}

		respBytes, err = io.ReadAll(resp.Body)
		if err != nil {
			// log.Fatal(err)
			printError(err, resp)
			return diag.FromErr(err)
		}

		// var jobResp ResponseJobJson
		err = json.Unmarshal(respBytes, &jobResp)
		if err != nil {
			// log.Fatal(err)
			printError(err, resp)
			return diag.FromErr(err)
		}
		fmt.Fprintf(os.Stderr, "Resp `Job`: %v\n", jobResp)

		if jobResp.Job.Status == "FINISHED" {
			break
		}

		if jobResp.Job.Status != "RUNNING" && jobResp.Job.Status != "DEFINED" {
			IsDeleteSuccess = false
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to DELETE resource",
			})
			break
		}
	}

	if IsDeleteSuccess {
		d.SetId("")
	}

	return diags
}

func newConfiguration(url string) *openapi.Configuration {
	cfg := &openapi.Configuration{
		DefaultHeader: make(map[string]string),
		UserAgent:     "OpenAPI-Generator/1.0.0/go",
		Debug:         false,
		Servers: openapi.ServerConfigurations{
			{
				URL:         url,
				Description: "No description provided",
			},
		},
		OperationServers: map[string]openapi.ServerConfigurations{},
	}
	return cfg
}

func printError(err error, resp *http.Response) {
	fmt.Fprintf(os.Stderr, "Error : %v\n", err)
	if resp != nil {
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", resp)
	}
}
