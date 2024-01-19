package ovh

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func resourceCloudProjectContainerRegistryIPRestrictionsRegistry() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudProjectContainerRegistryIPRestrictionsRegistryPut,
		Delete: resourceCloudProjectContainerRegistryIPRestrictionsRegistryDelete,
		Update: resourceCloudProjectContainerRegistryIPRestrictionsRegistryPut,
		Read:   resourceCloudProjectContainerIPRestrictionsRegistryRead,
		Importer: &schema.ResourceImporter{
			State: resourceCloudProjectContainerRegistryIPRestrictionsRegistryImportState,
		},

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Description: "Service name",
				ForceNew:    true,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_CLOUD_PROJECT_SERVICE", nil),
			},
			"registry_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "RegistryID",
				Required:    true,
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

func resourceCloudProjectContainerRegistryIPRestrictionsRegistryImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	givenId := d.Id()
	log.Printf("[DEBUG] Importing cloud project registry IP restrictions of registry type %s", givenId)

	splitId := strings.SplitN(givenId, "/", 3)

	if len(splitId) != 2 {
		return nil, fmt.Errorf("import Id is not service_name/registryid formatted")
	}

	serviceName := splitId[0]
	registryID := splitId[1]

	d.SetId(serviceName + "/" + registryID)
	d.Set("registry_id", registryID)
	d.Set("service_name", serviceName)

	results := make([]*schema.ResourceData, 1)
	results[0] = d

	return results, nil
}

func resourceCloudProjectContainerRegistryIPRestrictionsRegistryPut(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)
	registryID := d.Get("registry_id").(string)

	endpoint := fmt.Sprintf(
		"/cloud/project/%s/containerRegistry/%s/ipRestrictions/registry",
		url.PathEscape(serviceName),
		url.PathEscape(registryID),
	)

	params := (&CloudProjectContainerRegistryOIDCCreateOpts{}).FromResource(d)
	var res []CloudProjectContainerRegistryIPRestriction

	log.Printf("[DEBUG] Will create registry IP restrictions for registry %s in cloud project %s: %+v", registryID, serviceName, params)

	err := config.OVHClient.Put(endpoint, params, res)
	if err != nil {
		return fmt.Errorf("Error calling put %s:\n\t %q", endpoint, err)
	}

	d.SetId(serviceName + "/" + registryID)

	log.Printf("[DEBUG] Registry %s IP restrictions of registry type are created", registryID)

	return resourceCloudProjectContainerRegistryUserRead(d, meta)
}

func resourceCloudProjectContainerIPRestrictionsRegistryRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	registryID := d.Get("registry_id").(string)

	log.Printf("[DEBUG] Will read cloud project registry IP restrictions of registry type %s for project: %s", registryID, serviceName)

	endpoint := fmt.Sprintf(
		"/cloud/project/%s/containerRegistry/%s/ipRestrictions/registry",
		url.PathEscape(serviceName),
		url.PathEscape(registryID),
	)

	ipRestrictions := []CloudProjectContainerRegistryIPRestriction{}

	if err := config.OVHClient.Get(endpoint, &ipRestrictions); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	mapIPRestrictions := make([]map[string]interface{}, len(ipRestrictions))
	for i, ipRestriction := range ipRestrictions {
		mapIPRestrictions[i] = ipRestriction.ToMap()
	}

	d.Set("ip_restrictions", ipRestrictions)
	d.SetId(serviceName + "/" + registryID)

	log.Printf("[DEBUG] Read Registry IP Restrictions %+v", ipRestrictions)

	return nil
}

func resourceCloudProjectContainerRegistryIPRestrictionsRegistryDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	registryID := d.Get("registry_id").(string)

	log.Printf("[DEBUG] Will delete registry IP restrictions for registry %s in cloud project: %s", registryID, serviceName)

	var params []CloudProjectContainerRegistryIPRestriction
	var res []CloudProjectContainerRegistryIPRestriction

	endpoint := fmt.Sprintf(
		"/cloud/project/%s/containerRegistry/%s/ipRestrictions/registry",
		url.PathEscape(serviceName),
		url.PathEscape(registryID),
	)

	err := config.OVHClient.Put(endpoint, params, res)
	if err != nil {
		return fmt.Errorf("Error calling put %s:\n\t %q", endpoint, err)
	}

	return nil
}
