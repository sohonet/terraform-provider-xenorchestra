package client

type Host struct {
	// TODO: Not sure the difference between these two
	Id        string
	Uuid      string
	NameLabel string
	PoolId    string
}

func (t Host) Compare(obj map[string]interface{}) bool {
	name_label := obj["name_label"].(string)
	if t.NameLabel == name_label {
		return true
	}
	return false
}

func (p Host) New(obj map[string]interface{}) XoObject {
	id := obj["id"].(string)
	uuid := obj["uuid"].(string)
	name_label := obj["name_label"].(string)
	poolId := obj["$poolId"].(string)
	return Host{
		Id:        id,
		NameLabel: name_label,
		Uuid:      uuid,
		PoolId:    poolId,
	}
}

func (c *Client) GetHost(name string) (Host, error) {
	obj, err := c.FindFromGetAllObjects(Host{NameLabel: name})
	host := obj.(Host)

	if err != nil {
		return host, err
	}

	return host, nil
}
