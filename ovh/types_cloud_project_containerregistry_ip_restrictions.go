package ovh

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

type CloudProjectContainerRegistryIPRestriction struct {
	CreatedAt   string `json:"createdAt,omitempty"`
	Description string `json:"description,omitempty"`
	IPBlock     string `json:"ipBlock"`
	UpdatedAt   string `json:"updatedAt,omitempty"`
}

type CloudProjectContainerRegistryIPRestrictionCreateOpts struct {
	IPRestrictions []CloudProjectContainerRegistryIPRestriction `json:"ipRestrictions"`
}

type CloudProjectContainerRegistryIPRestrictionUpdateOpts struct {
	IPRestrictions []CloudProjectContainerRegistryIPRestriction `json:"ipRestrictions"`
}

func (opts *CloudProjectContainerRegistryIPRestrictionCreateOpts) FromResource(d *schema.ResourceData) *CloudProjectContainerRegistryIPRestrictionCreateOpts {
	opts.IPRestrictions = loadIPRestrictionsFromResource(d.Get("ip_restrictions"))

	return opts
}

func (opts *CloudProjectContainerRegistryIPRestrictionUpdateOpts) FromResource(d *schema.ResourceData) *CloudProjectContainerRegistryIPRestrictionUpdateOpts {
	opts.IPRestrictions = loadIPRestrictionsFromResource(d.Get("ip_restrictions"))

	return opts
}
func loadIPRestrictionsFromResource(i interface{}) []CloudProjectContainerRegistryIPRestriction {
	ips := make([]CloudProjectContainerRegistryIPRestriction, 0)
	iprestrictionsSet := i.(*schema.Set).List()

	for _, ipSet := range iprestrictionsSet {
		ips = append(ips, CloudProjectContainerRegistryIPRestriction{
			CreatedAt:   ipSet.(map[string]interface{})["created_at"].(string),
			Description: ipSet.(map[string]interface{})["description"].(string),
			IPBlock:     ipSet.(map[string]interface{})["ip_block"].(string),
			UpdatedAt:   ipSet.(map[string]interface{})["updated_at"].(string),
		})
	}
	return ips
}

func (r CloudProjectContainerRegistryIPRestriction) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["created_at"] = r.CreatedAt
	obj["description"] = r.Description
	obj["ip_block"] = r.IPBlock
	obj["updated_at"] = r.UpdatedAt

	return obj
}
