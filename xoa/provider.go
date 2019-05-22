package xoa

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terra-farm/terraform-provider-xenorchestra/client"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("XOA_URL", nil),
				Description: "Hostname of the xoa router",
			},
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("XOA_USER", nil),
				Description: "User account for xoa api",
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("XOA_PASSWORD", nil),
				Description: "Password for xoa api",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"xenorchestra_vm":           resourceRecord(),
			"xenorchestra_cloud_config": resourceCloudConfigRecord(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"xenorchestra_template": dataSourceXoaTemplate(),
			"xenorchestra_pif":      dataSourceXoaPIF(),
			"xenorchestra_host":     dataSourceXoaHost(),
			"xenorchestra_disk":     dataSourceXoaDisk(),
		},
		ConfigureFunc: xoaConfigure,
	}
}

func xoaConfigure(d *schema.ResourceData) (c interface{}, err error) {
	url := d.Get("url").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	c = client.Config{
		Url:      url,
		Username: username,
		Password: password,
	}
	return c, nil
}
