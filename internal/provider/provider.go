package provider

import (
	"context"
	"crypto/tls"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/severalnines/clustercontrol-client-sdk/go/pkg/openapi"
	"log/slog"
	"net/http"
	"strings"
)

func Provider() *schema.Provider {
	funcName := "Provider"
	slog.Debug(funcName)

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			API_USER: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				//Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(API_USER, nil),
			},
			API_USER_PW: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				//Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc(API_USER_PW, nil),
			},
			CONTROLLER_URL: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				//Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(CONTROLLER_URL, nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			RESOURCE_DB_CLUSTER:             resourceDbCluster(),
			RESOURCE_DB_LOAD_BALANCER:       resourceDbLoadBalancer(),
			RESOURCE_DB_CLUSTER_MAINTENANCE: resourceDbClusterMaintenance(),
			RESOURCE_DB_CLUSTER_BACKUP:      resourceDbClusterBackup(),
		},
		DataSourcesMap:       map[string]*schema.Resource{},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	funcName := "providerConfigure"
	slog.Debug(funcName)

	apiUser := d.Get(API_USER).(string)
	apiUserPw := d.Get(API_USER_PW).(string)
	apiUrl := d.Get(CONTROLLER_URL).(string)

	gCfg = newConfiguration(apiUrl)
	gApiClient = openapi.NewAPIClient(gCfg)
	authenticate := openapi.NewAuthenticate(CMON_OP_AUTHENTICATE_WITH_PW)
	authenticate.SetUserName(apiUser /*os.Getenv("API_USER")*/)
	authenticate.SetPassword(apiUserPw /*os.Getenv("API_USER_PW")*/)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	if strings.EqualFold(apiUser, CMON_MOCK_USER) {
		return nil, nil
	}

	resp, err := gApiClient.AuthAPI.AuthPost(ctx).Authenticate(*authenticate).Execute()
	if err != nil {
		PrintError(err, resp)
		return nil, diag.FromErr(err)
	}
	//slog.Debug("providerConfigure", "Resp `AuthApi.AuthPost`", resp)
	slog.Info("providerConfigure", "Resp `AuthApi.AuthPost`", resp)

	// fmt.Println("#Cookies: ", len(resp.Cookies()))
	slog.Debug("providerConfigure", "Num cookies", len(resp.Cookies()))
	for _, cookie := range resp.Cookies() {
		slog.Debug("providerConfigure", "Cookie", cookie)
		gNewCtx = context.WithValue(ctx, "cookie", cookie)
		break
	}

	return gApiClient, diags
}

func newConfiguration(url string) *openapi.Configuration {
	funcName := "newConfiguration"
	slog.Debug(funcName)

	cfg := &openapi.Configuration{
		DefaultHeader: make(map[string]string),
		UserAgent:     "OpenAPI-Generator/1.0.0/go",
		//Debug:         true,
		Debug: false,
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

func PrintError(err error, resp *http.Response) {
	slog.Error(err.Error())
	if resp != nil {
		slog.Error("", "Full HTTP response", resp)
	}
}
