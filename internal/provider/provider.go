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

type ProviderDetails struct {
	//SessionIdCtx context.Context
	SessionCookie *http.Cookie
	Cfg           *openapi.Configuration
	ApiClient     *openapi.APIClient
}

func NewProviderDetails(cookie *http.Cookie, configuration *openapi.Configuration, client *openapi.APIClient) *ProviderDetails {
	return &ProviderDetails{
		//SessionIdCtx: ctx,
		SessionCookie: cookie,
		Cfg:           configuration,
		ApiClient:     client,
	}
}

func Provider() *schema.Provider {
	funcName := "Provider"
	slog.Debug(funcName)

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			API_USER: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "ClusterControl API user",
				DefaultFunc: schema.EnvDefaultFunc(API_USER, nil),
			},
			API_USER_PW: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "ClusterControl API user's password",
				DefaultFunc: schema.EnvDefaultFunc(API_USER_PW, nil),
			},
			CONTROLLER_URL: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "ClusterControl controller url e.g. (https://cc-host:9501/v2)",
				DefaultFunc: schema.EnvDefaultFunc(CONTROLLER_URL, nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			RESOURCE_DB_CLUSTER:                 resourceDbCluster(),
			RESOURCE_DB_CLUSTER_MAINTENANCE:     resourceDbClusterMaintenance(),
			RESOURCE_DB_CLUSTER_BACKUP:          resourceDbClusterBackup(),
			RESOURCE_DB_CLUSTER_BACKUP_SCHEDULE: resourceDbClusterBackupSchedule(),
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

	cfg := newConfiguration(apiUrl)
	apiClient := openapi.NewAPIClient(cfg)
	authenticate := openapi.NewAuthenticate(CMON_OP_AUTHENTICATE_WITH_PW)
	authenticate.SetUserName(apiUser /*os.Getenv("API_USER")*/)
	authenticate.SetPassword(apiUserPw /*os.Getenv("API_USER_PW")*/)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	if strings.EqualFold(apiUser, CMON_MOCK_USER) {
		return nil, nil
	}

	resp, err := apiClient.AuthAPI.AuthPost(ctx).Authenticate(*authenticate).Execute()
	if err != nil {
		PrintError(err, resp)
		return nil, diag.FromErr(err)
	}
	slog.Debug("providerConfigure", "Resp `AuthApi.AuthPost`", resp)

	// fmt.Println("#Cookies: ", len(resp.Cookies()))
	slog.Debug("providerConfigure", "Num cookies", len(resp.Cookies()))
	//var ccSessionIdCtx context.Context
	var sessionCookie http.Cookie
	for _, cookie := range resp.Cookies() {
		slog.Debug("providerConfigure", "Cookie", cookie)
		//ccSessionIdCtx = context.WithValue(ctx, "cookie", cookie)
		sessionCookie = *cookie
		break
	}

	prividerDetails := NewProviderDetails(&sessionCookie, cfg, apiClient)

	return prividerDetails, diags
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
