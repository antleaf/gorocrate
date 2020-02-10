package main

type entity struct {
	Properties map[string]interface{}
}

func NewEntity(entityType,id string) entity {
	var e entity
	e.Properties = make(map[string]interface{})
	e.Properties["@id"] = id
	e.Properties["@type"] = entityType
	return e
}

func NewEntityFromMap(m map[string]interface{}) entity {
	var e entity
	e.Properties = make(map[string]interface{})
	for k,v := range m {
		e.Properties[k] = v
	}
	return e
}

func (e *entity) SetProperty(k string, v interface{}) {
	e.Properties[k] = v
}

func (e *entity) GetProperty(k string) interface{} {
	return e.Properties[k]
}

func (e *entity) GetID() string {
	if e.GetProperty("@id") != nil {
		return e.GetProperty("@id").(string)
	} else {
		return ""
	}
}

func (e *entity) GetType() string {
	if e.GetProperty("@type") != nil {
		return e.GetProperty("@type").(string)
	} else {
		return ""
	}
}

