package xoa

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terra-farm/terraform-provider-xenorchestra/client"
)

func dataSourceXoaPIF() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePIFRead,
		Schema: map[string]*schema.Schema{
			"attached": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"device": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"host": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"network": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"pool_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"uuid": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"vlan": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"host_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourcePIFRead(d *schema.ResourceData, m interface{}) error {
	config := m.(client.Config)
	c, err := client.NewClient(config)

	if err != nil {
		return err
	}
	//fmt.Printf("%#v\n\n", d)
	device := d.Get("device").(string)
	vlan := d.Get("vlan").(int)
	host_id := d.Get("host_id").(string)

	pif, err := c.GetPIFByDeviceHost(device, vlan, host_id)

	//pif, err := c.GetPIFByDevice(device, vlan)

	if err != nil {
		return err
	}

	if _, ok := err.(client.NotFound); ok {
		d.SetId("")
		return nil
	}

	d.SetId(pif.Id)
	d.Set("uuid", pif.Uuid)
	d.Set("device", pif.Device)
	d.Set("host", pif.Host)
	d.Set("attached", pif.Attached)
	d.Set("pool_id", pif.PoolId)
	d.Set("network", pif.Network)
	d.Set("vlan", pif.Vlan)
	return nil
}
