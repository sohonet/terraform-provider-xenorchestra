package client

type Disk struct {
	Container string
	Id        string
	Uuid      string
	PoolId    string
	Name      string
}

func (p Disk) Compare(obj map[string]interface{}) bool {
	name_label := obj["name_label"].(string)
	hostID := obj["$container"].(string)

	if p.Container == hostID && p.Name == name_label {
		return true
	}
	return false
}

func (p Disk) New(obj map[string]interface{}) XoObject {
	id := obj["id"].(string)
	container := obj["$container"].(string)
	uuid := obj["uuid"].(string)
	poolId := obj["$poolId"].(string)
	name := obj["name_label"].(string)
	return Disk{
		Container: container,
		Id:        id,
		Uuid:      uuid,
		PoolId:    poolId,
		Name:      name,
	}
}

func (c *Client) GetDisk(name string, host_id string) (Disk, error) {
	obj, err := c.FindFromGetAllObjects(Disk{Name: name, Container: host_id})
	DISK := obj.(Disk)

	if err != nil {
		return DISK, err
	}

	return DISK, nil
}
