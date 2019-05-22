package xoa

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terra-farm/terraform-provider-xenorchestra/client"
)

func dataSourceXoaHost() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceXoaHostRead,
		Schema: map[string]*schema.Schema{
			"host": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"pool_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"uuid": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceXoaHostRead(d *schema.ResourceData, m interface{}) error {
	config := m.(client.Config)
	c, err := client.NewClient(config)

	if err != nil {
		return err
	}

	hostname := d.Get("host").(string)

	host, err := c.GetHost(hostname)

	if err != nil {
		return err
	}

	if _, ok := err.(client.NotFound); ok {
		d.SetId("")
		return nil
	}

	d.SetId(host.Uuid)
	d.Set("host", host.NameLabel)
	d.Set("pool_id", host.PoolId)
	d.Set("uuid", host.Uuid)
	return nil
}
