package provider

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func resourceDbLoadBalancer() *schema.Resource {
	return &schema.Resource{
		CreateContext: nil,
		ReadContext:   nil,
		UpdateContext: nil,
		DeleteContext: nil,
		Importer:      &schema.ResourceImporter{},
		Schema: map[string]*schema.Schema{
			"cluster_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the resource, also acts as it's unique ID",
				ForceNew:    true,
				// ValidateFunc: validateName,
			},
		},
	}
}
