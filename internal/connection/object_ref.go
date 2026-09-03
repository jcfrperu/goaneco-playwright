package connection

import "encoding/json"

type ObjectRef struct {
	guid        string
	initializer json.RawMessage
}

// NewObjectRef creates a new object reference.
func NewObjectRef(guid string, initializer json.RawMessage) *ObjectRef {
	return &ObjectRef{
		guid:        guid,
		initializer: initializer,
	}
}

func (c *ObjectRef) GUID() string                 { return c.guid }
func (c *ObjectRef) Initializer() json.RawMessage { return c.initializer }
