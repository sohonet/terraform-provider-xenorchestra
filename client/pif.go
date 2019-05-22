package client

type PIF struct {
	Device   string
	Host     string
	Network  string
	Id       string
	Uuid     string
	PoolId   string
	Attached bool
	Vlan     int
}

func (p PIF) Compare(obj map[string]interface{}) bool {
	device := obj["device"].(string)
	vlan := int(obj["vlan"].(float64))
	hostID := obj["$host"].(string)

	//fmt.Printf("%#v\n\n", obj)

	if p.Vlan == vlan && p.Device == device && len(p.Host) == 0 {
		return true
	} else if p.Vlan == vlan && p.Device == device && hostID == p.Host {
		return true
	}
	return false
}

func (p PIF) New(obj map[string]interface{}) XoObject {
	id := obj["id"].(string)
	device := obj["device"].(string)
	attached := obj["attached"].(bool)
	network := obj["$network"].(string)
	uuid := obj["uuid"].(string)
	poolId := obj["$poolId"].(string)
	host := obj["$host"].(string)
	vlan := int(obj["vlan"].(float64))
	return PIF{
		Device:   device,
		Host:     host,
		Network:  network,
		Id:       id,
		Uuid:     uuid,
		PoolId:   poolId,
		Attached: attached,
		Vlan:     vlan,
	}
}

func (c *Client) GetPIFByDevice(dev string, vlan int) (PIF, error) {
	obj, err := c.FindFromGetAllObjects(PIF{Device: dev, Vlan: vlan})
	pif := obj.(PIF)

	if err != nil {
		return pif, err
	}

	return pif, nil
}
func (c *Client) GetPIFByDeviceHost(dev string, vlan int, host_id string) (PIF, error) {
	obj, err := c.FindFromGetAllObjects(PIF{Device: dev, Vlan: vlan, Host: host_id})
	pif := obj.(PIF)

	if err != nil {
		return pif, err
	}

	return pif, nil
}
