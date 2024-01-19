package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccCloudProjectContainerRegistryIPRestrictionsManagementConfig = `
data "ovh_cloud_project_capabilities_containerregistry_filter" "registryCap" {
	service_name = "%s"
	plan_name    = "SMALL"
	region       = "%s"
}

resource "ovh_cloud_project_containerregistry" "registry" {
	service_name = data.ovh_cloud_project_capabilities_containerregistry_filter.registryCap.service_name
	plan_id      = data.ovh_cloud_project_capabilities_containerregistry_filter.registryCap.id
	name         = "%s"
	region       = data.ovh_cloud_project_capabilities_containerregistry_filter.registryCap.region
}

resource "ovh_cloud_project_containerregistry_ip_restrictions_management" "my-mgt-iprestrictions" {
	service_name = ovh_cloud_project_containerregistry.registry.service_name
	registry_id  = ovh_cloud_project_containerregistry.registry.id
	
	ip_restrictions = [
					   { ip_block = "121.121.121.121/32" ,
						 description = "my awesome ip"  
					   }
                      ]
}
`
const testAccCloudProjectContainerRegistryIPRestrictionsManagementConfigUpdated = `
data "ovh_cloud_project_capabilities_containerregistry_filter" "registryCap" {
	service_name = "%s"
	plan_name    = "SMALL"
	region       = "%s"
}

resource "ovh_cloud_project_containerregistry" "registry" {
	service_name = data.ovh_cloud_project_capabilities_containerregistry_filter.registryCap.service_name
	plan_id      = data.ovh_cloud_project_capabilities_containerregistry_filter.registryCap.id
	name         = "%s"
	region       = data.ovh_cloud_project_capabilities_containerregistry_filter.registryCap.region
}

resource "ovh_cloud_project_containerregistry_ip_restrictions_management" "my-mgt-iprestrictions" {
	service_name = ovh_cloud_project_containerregistry.registry.service_name
	registry_id  = ovh_cloud_project_containerregistry.registry.id
	
	ip_restrictions = [
					   { ip_block = "122.122.122.122/32" ,
						 description = "my new awesome ip description"  
					   }
                      ]
}
`

func TestAccCloudProjectContainerRegistryIPRestrictionsManagement_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	registryName := acctest.RandomWithPrefix(test_prefix)
	region := os.Getenv("OVH_CLOUD_PROJECT_CONTAINERREGISTRY_REGION_TEST")
	resourceName := "ovh_cloud_project_containerregistry_ip_restrictions.my-mgt-iprestriction"

	config := fmt.Sprintf(
		testAccCloudProjectContainerRegistryIPRestrictionsManagementConfig,
		serviceName,
		region,
		registryName,
	)

	configUpdated := fmt.Sprintf(
		testAccCloudProjectContainerRegistryIPRestrictionsManagementConfigUpdated,
		serviceName,
		region,
		registryName,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckContainerRegistryIPRestrictionManagement(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_containerregistry_ip_restrictions.my-mgt-iprestriction", "ip_restrictions.0.ip_block", "121.121.121.121/32"),
					resource.TestCheckResourceAttr("ovh_cloud_project_containerregistry_ip_restrictions.my-mgt-iprestriction", "ip_restrictions.0.description", "my new awesome ip"),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_containerregistry_ip_restrictions.my-mgt-iprestriction", "ip_restrictions.0.created_at"),
				),
			},
			{
				Config: configUpdated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_containerregistry_ip_restrictions.my-mgt-iprestriction", "ip_restrictions.0.ip_block", "121.121.121.121/32"),
					resource.TestCheckResourceAttr("ovh_cloud_project_containerregistry_ip_restrictions.my-mgt-iprestriction", "ip_restrictions.0.description", "my new awesome ip description"),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_containerregistry_ip_restrictions.my-mgt-iprestriction", "ip_restrictions.0.created_at"),
					resource.TestCheckResourceAttrSet("ovh_cloud_project_containerregistry_ip_restrictions.my-mgt-iprestriction", "ip_restrictions.0.updated_at"),
				),
			},
			{
				Config:       configUpdated,
				Destroy:      true,
				ResourceName: resourceName,
			},
		},
	})
}
