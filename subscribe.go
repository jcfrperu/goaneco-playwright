package playwright

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/jcfrperu/goaneco-playwright/protocol"
)

// networkEventPayload is the common shape of the request, response, requestFinished and
// requestFailed BrowserContext events. Only the field relevant to each event is populated.
type networkEventPayload struct {
	Page        *protocol.Page    `json:"page,omitempty"`
	Request     protocol.Request  `json:"request"`
	Response    protocol.Response `json:"response"`
	FailureText *string           `json:"failureText,omitempty"`
}

// subscribeNetwork registers handler for a network event of the BrowserContext and returns an
// idempotent cancel function. If page is non-nil, only events for that page are delivered.
// guidOf selects the request/response guid from the payload; build creates the public object.
func subscribeNetwork[T any](
	c *BrowserContext,
	page *Page,
	event string,
	guidOf func(networkEventPayload) string,
	build func(parent ChannelOwner, guid string, raw json.RawMessage) *T,
	handler func(*T),
) func() {
	subCtx, subCancel := context.WithTimeout(context.Background(), defaultSubscriptionTimeout)
	defer subCancel()
	c.retainSubscription(subCtx, event)

	parent := c.owner
	if page != nil {
		parent = page.owner
	}
	id := c.owner.conn.OnEvent(c.owner.guid, event, func(params json.RawMessage) {
		var payload networkEventPayload
		if err := json.Unmarshal(params, &payload); err != nil {
			return
		}
		if page != nil && payload.Page != nil && payload.Page.Guid != page.owner.guid {
			return
		}
		guid := guidOf(payload)
		var v *T
		if event == "request" || event == "requestFinished" || event == "requestFailed" {
			if c != nil {
				req := c.requestFromGUID(guid, parent.Initializer(guid))
				if anyVal, ok := any(req).(*T); ok {
					v = anyVal
				}
			}
		}
		if v == nil {
			v = build(parent, guid, parent.Initializer(guid))
		}
		if v != nil {
			if req, ok := any(v).(*NetworkRequest); ok {
				if page != nil {
					req.page = page
				} else if payload.Page != nil && c != nil {
					req.page = c.pageByGUID(payload.Page.Guid)
				}
				if payload.FailureText != nil {
					req.failure = payload.FailureText
				}
				if event == "requestFinished" || event == "requestFailed" {
					req.didFailOrFinish = true
				}
			}
			if resp, ok := any(v).(*NetworkResponse); ok {
				if page != nil {
					resp.page = page
				} else if payload.Page != nil && c != nil {
					resp.page = c.pageByGUID(payload.Page.Guid)
				}
				if req := resp.Request(); req != nil {
					req.existingResponse = resp
					if req.page == nil {
						req.page = resp.page
					}
				}
			}
			handler(v)
		}
	})
	var once sync.Once
	return func() {
		once.Do(func() {
			c.owner.conn.OffEvent(c.owner.guid, event, id)
			c.releaseSubscription(event)
		})
	}
}

func requestGuid(p networkEventPayload) string  { return p.Request.Guid }
func responseGuid(p networkEventPayload) string { return p.Response.Guid }
