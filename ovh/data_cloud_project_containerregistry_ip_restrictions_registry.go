package ovh

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCloudProjectContainerRegistryIPRestrictionsRegistry() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudProjectContainerRegistryIPRestrictionsRegistryRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Description: "Service name",
				Required:    true,
				ForceNew:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_CLOUD_PROJECT_SERVICE", nil),
			},
			"registry_id": {
				Type:        schema.TypeString,
				Description: "Registry ID",
				Required:    true,
				ForceNew:    true,
			},

			"ip_restrictions": {
				Type:        schema.TypeList,
				Description: "List your IP restrictions applied on artifact manager component",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"created_at": {
							Type:        schema.TypeString,
							Description: "IP Restrictions creation date",
							Required:    true,
						},
						"description": {
							Type:        schema.TypeString,
							Description: "The Description of Whitelisted IP block",
							Optional:    true,
						},
						"ip_block": {
							Type:        schema.TypeString,
							Description: "Whitelisted IP block",
							Required:    true,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Description: "IP Restrictions update date",
							Required:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceCloudProjectContainerRegistryIPRestrictionsRegistryRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	registryID := d.Get("registry_id").(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/containerRegistry/%s/ipRestrictions/registry", serviceName, registryID)
	ipRestrictions := []CloudProjectContainerRegistryIPRestriction{}

	log.Printf("[DEBUG] Will read Registry IP Restrictions from registry %s and project: %s", registryID, serviceName)
	err := config.OVHClient.Get(endpoint, ipRestrictions)
	if err != nil {
		return fmt.Errorf("calling get %s %w", endpoint, err)
	}

	mapIPRestrictions := make([]map[string]interface{}, len(ipRestrictions))
	for i, ipRestriction := range ipRestrictions {
		mapIPRestrictions[i] = ipRestriction.ToMap()
	}
	d.Set("ip_restrictions", ipRestrictions)
	d.SetId(registryID)

	log.Printf("[DEBUG] Read Registry IP Restrictions %+v", ipRestrictions)

	return nil
}
