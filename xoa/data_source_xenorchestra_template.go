package xoa

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terra-farm/terraform-provider-xenorchestra/client"
)

func dataSourceXoaTemplate() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTemplateRead,
		Schema: map[string]*schema.Schema{
			"name_label": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"pool_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"uuid": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceTemplateRead(d *schema.ResourceData, m interface{}) error {
	config := m.(client.Config)
	c, err := client.NewClient(config)

	if err != nil {
		return err
	}
	fmt.Printf("%#v\n\n", d)
	templateName := d.Get("name_label").(string)
	hostUUID := d.Get("pool_id").(string)

	tmpl, err := c.GetTemplate(templateName, hostUUID)

	if err != nil {
		return err
	}

	if _, ok := err.(client.NotFound); ok {
		d.SetId("")
		return nil
	}

	d.SetId(tmpl.Id)
	d.Set("uuid", tmpl.Uuid)
	d.Set("name_label", tmpl.NameLabel)
	return nil
}
