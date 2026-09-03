package playwright

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jcfrperu/goaneco-playwright/internal/connection"
)

// ChannelOwner is the base struct embedded by all remote Playwright protocol objects.
// It provides generic IPC routing capabilities.
type ChannelOwner struct {
	conn *connection.Connection
	guid string
}

// GUID returns the unique identifier of the remote object.
func (c *ChannelOwner) GUID() string {
	return c.guid
}

// child creates a derived ChannelOwner sharing the same connection.
func (c *ChannelOwner) child(guid string) ChannelOwner {
	return ChannelOwner{
		conn: c.conn,
		guid: guid,
	}
}

// Initializer returns the raw JSON initializer for the given GUID from the connection.
func (c *ChannelOwner) Initializer(guid string) json.RawMessage {
	if obj, ok := c.conn.GetObject(guid); ok {
		return obj.Initializer()
	}
	return nil
}

type requestMetadata struct {
	WallTime int64    `json:"wallTime"`
	Timeout  *float64 `json:"timeout,omitempty"`
}

type requestEnvelope struct {
	ID       *int            `json:"id"`
	GUID     string          `json:"guid"`
	Method   string          `json:"method"`
	Params   any             `json:"params"`
	Metadata requestMetadata `json:"metadata"`
}

// SendMessageRequest wraps an IPC request to the underlying connection,
// serializing the given parameters and waiting for the response.
// Returns json.RawMessage containing valid JSON ready for unmarshaling.
func (c *ChannelOwner) SendMessageRequest(ctx context.Context, method string, params any) (json.RawMessage, error) {
	return c.sendMessage(ctx, method, params, nil)
}

// sendWithTimeout sends an IPC request with an explicit timeout in the metadata.
// The Playwright server uses this timeout to bound the operation (e.g., an HTTP fetch).
func (c *ChannelOwner) sendWithTimeout(ctx context.Context, method string, params any, timeoutMs float64) (json.RawMessage, error) {
	t := timeoutMs
	return c.sendMessage(ctx, method, params, &t)
}

func (c *ChannelOwner) sendMessage(ctx context.Context, method string, params any, timeoutMs *float64) (json.RawMessage, error) {
	if params == nil {
		params = struct{}{}
	}
	id := c.conn.NextID()

	msgStruct := requestEnvelope{
		ID:     &id,
		GUID:   c.guid,
		Method: method,
		Params: params,
		Metadata: requestMetadata{
			WallTime: time.Now().UnixMilli(),
			Timeout:  timeoutMs,
		},
	}

	msgBytes, err := json.Marshal(msgStruct)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal IPC message: %w", err)
	}

	respMsg, err := c.conn.SendRequest(ctx, id, msgBytes)
	if err != nil {
		return nil, err
	}

	return respMsg.Result, nil
}
