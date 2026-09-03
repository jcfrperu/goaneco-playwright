package playwright

import (
	"context"
	"encoding/json"
	"fmt"
)

// JSHandle represents an in-page JavaScript object reference.
type JSHandle struct {
	owner ChannelOwner
}

type jsHandleEvalParams struct {
	Expression string                `json:"expression"`
	Arg        serializedArgumentRaw `json:"arg"`
}

// Evaluate executes a JavaScript expression in the context of the handle and returns the result.
func (h *JSHandle) Evaluate(ctx context.Context, expression string, arg ...any) (any, error) {
	var input any
	if len(arg) > 0 {
		input = arg[0]
	}
	result, err := h.owner.SendMessageRequest(ctx, "evaluateExpression", jsHandleEvalParams{
		Expression: expression,
		Arg:        serializeArgument(input),
	})
	if err != nil {
		return nil, fmt.Errorf("jsHandle.evaluate failed: %w", err)
	}
	var resp struct {
		Value serializedValueRaw `json:"value"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse jsHandle.evaluate response: %w", err)
	}
	return deserializeValue(resp.Value)
}

// EvaluateHandle executes a JavaScript expression and returns the result as a JSHandle.
func (h *JSHandle) EvaluateHandle(ctx context.Context, expression string, arg ...any) (*JSHandle, error) {
	var input any
	if len(arg) > 0 {
		input = arg[0]
	}
	result, err := h.owner.SendMessageRequest(ctx, "evaluateExpressionHandle", jsHandleEvalParams{
		Expression: expression,
		Arg:        serializeArgument(input),
	})
	if err != nil {
		return nil, fmt.Errorf("jsHandle.evaluateHandle failed: %w", err)
	}
	var resp struct {
		Handle struct {
			Guid string `json:"guid"`
		} `json:"handle"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse jsHandle.evaluateHandle response: %w", err)
	}
	return &JSHandle{owner: h.owner.child(resp.Handle.Guid)}, nil
}

// GetProperty returns a JSHandle for the named property of the object.
func (h *JSHandle) GetProperty(ctx context.Context, name string) (*JSHandle, error) {
	result, err := h.owner.SendMessageRequest(ctx, "getProperty", map[string]string{"name": name})
	if err != nil {
		return nil, fmt.Errorf("jsHandle.getProperty(%q) failed: %w", name, err)
	}
	var resp struct {
		Handle struct {
			Guid string `json:"guid"`
		} `json:"handle"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse jsHandle.getProperty response: %w", err)
	}
	return &JSHandle{owner: h.owner.child(resp.Handle.Guid)}, nil
}

// JSONValue returns the JSON-serializable value of the handle.
func (h *JSHandle) JSONValue(ctx context.Context) (any, error) {
	result, err := h.owner.SendMessageRequest(ctx, "jsonValue", struct{}{})
	if err != nil {
		return nil, fmt.Errorf("jsHandle.jsonValue failed: %w", err)
	}
	var resp struct {
		Value serializedValueRaw `json:"value"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse jsHandle.jsonValue response: %w", err)
	}
	return deserializeValue(resp.Value)
}

// Dispose releases the handle, freeing it in the browser process.
func (h *JSHandle) Dispose(ctx context.Context) error {
	_, err := h.owner.SendMessageRequest(ctx, "dispose", struct{}{})
	if err != nil {
		return fmt.Errorf("jsHandle.dispose failed: %w", err)
	}
	return nil
}
