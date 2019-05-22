package xoa

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/terra-farm/terraform-provider-xenorchestra/client"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
)

var logFile = "/tmp/terraform-provider-xenorchestra"
var XenLog io.Writer
var err error

func init() {
}

func resourceRecord() *schema.Resource {
	duration := 5 * time.Minute
	return &schema.Resource{
		Create: resourceVmCreate,
		Read:   resourceVmRead,
		Update: resourceVmUpdate,
		Delete: resourceVmDelete,
		Importer: &schema.ResourceImporter{
			State: RecordImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: &duration,
			Update: &duration,
		},
		Schema: map[string]*schema.Schema{
			"name_label": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"name_description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"template": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cloud_config": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"core_os": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"cpu_cap": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"cpu_weight": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"cpus": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"memory_max": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"network": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"network_id": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				Set: func(value interface{}) int {
					network := value.(map[string]interface{})

					return hashcode.String(network["network_id"].(string))
				},
			},
			"disk": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sr_id": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"name_label": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"size": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
				Set: func(value interface{}) int {
					var buf bytes.Buffer
					disk := value.(map[string]interface{})

					buf.WriteString(fmt.Sprintf("%s-", disk["sr_id"].(string)))
					buf.WriteString(fmt.Sprintf("%s-", disk["name_label"].(string)))
					buf.WriteString(fmt.Sprintf("%d-", disk["size"]))
					return hashcode.String(buf.String())
				},
			},
		},
	}
}

func resourceVmCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(client.Config)
	c, err := client.NewClient(config)

	if err != nil {
		return err
	}

	network_ids := []string{}
	networks := d.Get("network").(*schema.Set)

	for _, network := range networks.List() {
		net, _ := network.(map[string]interface{})

		network_ids = append(network_ids, net["network_id"].(string))
	}

	vdis := []client.VDI{}

	disks := d.Get("disk").(*schema.Set)

	for _, disk := range disks.List() {
		vdi, _ := disk.(map[string]interface{})

		vdis = append(vdis, client.VDI{
			SrId:      vdi["sr_id"].(string),
			NameLabel: vdi["name_label"].(string),
			Size:      vdi["size"].(int),
		})
	}

	vm, err := c.CreateVm(
		d.Get("name_label").(string),
		d.Get("name_description").(string),
		d.Get("template").(string),
		d.Get("cloud_config").(string),
		d.Get("cpus").(int),
		d.Get("memory_max").(int),
		network_ids,
		vdis,
	)

	if err != nil {
		return err
	}

	d.SetId(vm.Id)
	d.Set("cloud_config", d.Get("cloud_config").(string))
	d.Set("memory_max", d.Get("memory_max").(int))
	return nil
}

func resourceVmRead(d *schema.ResourceData, m interface{}) error {
	xoaId := d.Id()
	config := m.(client.Config)
	c, err := client.NewClient(config)

	if err != nil {
		return err
	}
	vmObj, err := c.GetVm(xoaId)
	if err != nil {
		return err
	}
	recordToData(*vmObj, d)
	return nil
}

func resourceVmUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceVmDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(client.Config)
	c, err := client.NewClient(config)

	if err != nil {
		return err
	}

	err = c.DeleteVm(d.Id())

	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func RecordImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	xoaId := d.Id()

	config := m.(client.Config)
	c, err := client.NewClient(config)

	if err != nil {
		return nil, err
	}

	vmObj, err := c.GetVm(xoaId)
	if err != nil {
		return nil, err
	}
	recordToData(*vmObj, d)
	return []*schema.ResourceData{d}, nil
}

func recordToData(resource client.Vm, d *schema.ResourceData) error {
	d.SetId(resource.Id)
	// d.Set("cloud_config", resource.CloudConfig)
	// d.Set("memory_max", resource.Memory.Size)
	d.Set("cpus", resource.CPUs.Number)
	d.Set("name_label", resource.NameLabel)
	d.Set("name_description", resource.NameDescription)
	return nil
}
