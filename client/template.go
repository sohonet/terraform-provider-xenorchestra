package client

type Template struct {
	// TODO: Not sure the difference between these two
	Id        string
	Uuid      string
	NameLabel string
	Container string
}

func (t Template) Compare(obj map[string]interface{}) bool {
	name_label := obj["name_label"].(string)
	host_uuid := obj["$container"].(string)
	//fmt.Printf("%#v%#v\n", name_label, t.NameLabel)
	//fmt.Printf("%#v%#v\n", host_uuid, t.container)
	if (t.NameLabel == name_label) && (t.Container == host_uuid) {
		return true
	}
	return false
}

func (p Template) New(obj map[string]interface{}) XoObject {
	id := obj["id"].(string)
	uuid := obj["uuid"].(string)
	name_label := obj["name_label"].(string)
	poolId := obj["$container"].(string)
	return Template{
		Id:        id,
		NameLabel: name_label,
		Uuid:      uuid,
		Container: poolId,
	}
}

func (c *Client) GetTemplate(name string, uuid string) (Template, error) {
	obj, err := c.FindFromGetAllObjects(Template{NameLabel: name, Container: uuid})
	template := obj.(Template)

	if err != nil {
		return template, err
	}

	return template, nil
}
