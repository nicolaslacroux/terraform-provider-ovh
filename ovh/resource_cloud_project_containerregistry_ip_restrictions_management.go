package ovh

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func resourceCloudProjectContainerRegistryIPRestrictionsManagement() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudProjectContainerRegistryIPRestrictionsManagementPut,
		Delete: resourceCloudProjectContainerRegistryIPRestrictionsManagementDelete,
		Update: resourceCloudProjectContainerRegistryIPRestrictionsManagementPut,
		Read:   resourceCloudProjectContainerIPRestrictionsManagementRead,
		Importer: &schema.ResourceImporter{
			State: resourceCloudProjectContainerRegistryIPRestrictionsManagementImportState,
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
				Description: "List your IP restrictions applied on Harbor UI and API",
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

func resourceCloudProjectContainerRegistryIPRestrictionsManagementImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	givenId := d.Id()
	log.Printf("[DEBUG] Importing cloud project registry IP restrictions of management type %s", givenId)

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

func resourceCloudProjectContainerRegistryIPRestrictionsManagementPut(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)
	registryID := d.Get("registry_id").(string)

	endpoint := fmt.Sprintf(
		"/cloud/project/%s/containerRegistry/%s/ipRestrictions/management",
		url.PathEscape(serviceName),
		url.PathEscape(registryID),
	)

	params := (&CloudProjectContainerRegistryOIDCCreateOpts{}).FromResource(d)
	var res []CloudProjectContainerRegistryIPRestriction

	log.Printf("[DEBUG] Will create management IP restrictions for registry %s in cloud project %s: %+v", registryID, serviceName, params)

	err := config.OVHClient.Put(endpoint, params, res)
	if err != nil {
		return fmt.Errorf("Error calling put %s:\n\t %q", endpoint, err)
	}

	d.SetId(serviceName + "/" + registryID)

	log.Printf("[DEBUG] Registry %s IP restrictions of registry management are created", registryID)

	return resourceCloudProjectContainerRegistryUserRead(d, meta) //TODO check that
}

func resourceCloudProjectContainerIPRestrictionsManagementRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	registryID := d.Get("registry_id").(string)

	log.Printf("[DEBUG] Will read management IP restrictions for registry %s in cloud project: %s", registryID, serviceName)

	endpoint := fmt.Sprintf(
		"/cloud/project/%s/containerRegistry/%s/ipRestrictions/management",
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

	log.Printf("[DEBUG] Read Management IP Restrictions %+v", ipRestrictions)

	return nil
}

func resourceCloudProjectContainerRegistryIPRestrictionsManagementDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	registryID := d.Get("registry_id").(string)

	log.Printf("[DEBUG] Will delete management IP restrictions for registry %s in cloud project: %s", registryID, serviceName)

	var params []CloudProjectContainerRegistryIPRestriction
	var res []CloudProjectContainerRegistryIPRestriction

	endpoint := fmt.Sprintf(
		"/cloud/project/%s/containerRegistry/%s/ipRestrictions/management",
		url.PathEscape(serviceName),
		url.PathEscape(registryID),
	)

	err := config.OVHClient.Put(endpoint, params, res)
	if err != nil {
		return fmt.Errorf("Error calling put %s:\n\t %q", endpoint, err)
	}

	return nil
}
