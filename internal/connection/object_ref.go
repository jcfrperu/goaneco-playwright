package connection

import (
	"encoding/json"
	"sync"
)

type ObjectRef struct {
	mu          sync.RWMutex
	guid        string
	typeName    string
	initializer json.RawMessage
}

// NewObjectRef creates a new object reference.
func NewObjectRef(guid string, args ...any) *ObjectRef {
	ref := &ObjectRef{guid: guid}
	if len(args) == 1 {
		if init, ok := args[0].(json.RawMessage); ok {
			ref.initializer = init
		}
	} else if len(args) >= 2 {
		if t, ok := args[0].(string); ok {
			ref.typeName = t
		}
		if init, ok := args[1].(json.RawMessage); ok {
			ref.initializer = init
		}
	}
	return ref
}

func (c *ObjectRef) GUID() string { return c.guid }
func (c *ObjectRef) Type() string { return c.typeName }

func (c *ObjectRef) Initializer() json.RawMessage {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.initializer
}

// UpdateFrameNavigated updates the cached URL and name in the frame initializer.
func (c *ObjectRef) UpdateFrameNavigated(params json.RawMessage) {
	var nav struct {
		URL  string `json:"url"`
		Name string `json:"name"`
	}
	if err := json.Unmarshal(params, &nav); err != nil || nav.URL == "" {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	var current map[string]any
	if len(c.initializer) > 0 {
		_ = json.Unmarshal(c.initializer, &current)
	}
	if current == nil {
		current = make(map[string]any)
	}
	current["url"] = nav.URL
	if nav.Name != "" {
		current["name"] = nav.Name
	}
	if updated, err := json.Marshal(current); err == nil {
		c.initializer = updated
	}
}

// UpdateFrameLoadState updates the cached loadStates in the frame initializer.
func (c *ObjectRef) UpdateFrameLoadState(params json.RawMessage) {
	var ev struct {
		Add    string `json:"add"`
		Remove string `json:"remove"`
	}
	if err := json.Unmarshal(params, &ev); err != nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	var current map[string]any
	if len(c.initializer) > 0 {
		_ = json.Unmarshal(c.initializer, &current)
	}
	if current == nil {
		current = make(map[string]any)
	}
	var states []any
	if rawStates, ok := current["loadStates"].([]any); ok {
		states = rawStates
	}
	if ev.Add != "" {
		found := false
		for _, s := range states {
			if s == ev.Add {
				found = true
				break
			}
		}
		if !found {
			states = append(states, ev.Add)
		}
	}
	if ev.Remove != "" {
		newStates := make([]any, 0, len(states))
		for _, s := range states {
			if s != ev.Remove {
				newStates = append(newStates, s)
			}
		}
		states = newStates
	}
	current["loadStates"] = states
	if updated, err := json.Marshal(current); err == nil {
		c.initializer = updated
	}
}
