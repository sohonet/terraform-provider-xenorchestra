package xoa

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terra-farm/terraform-provider-xenorchestra/client"
)

func dataSourceXoaDisk() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDiskRead,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"container": &schema.Schema{
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

func dataSourceDiskRead(d *schema.ResourceData, m interface{}) error {
	config := m.(client.Config)
	c, err := client.NewClient(config)

	if err != nil {
		return err
	}
	name_label := d.Get("name").(string)
	host_id := d.Get("container").(string)

	dsk, err := c.GetDisk(name_label, host_id)

	if err != nil {
		return err
	}

	if _, ok := err.(client.NotFound); ok {
		d.SetId("")
		return nil
	}

	d.SetId(dsk.Id)
	d.Set("uuid", dsk.Uuid)
	d.Set("name", dsk.Name)
	d.Set("container", dsk.Container)
	d.Set("pool_id", dsk.PoolId)
	return nil
}
