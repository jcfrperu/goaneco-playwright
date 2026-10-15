package playwright

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"
)

// harFile is the top-level structure of an HTTP Archive (HAR) file.
type harFile struct {
	Log struct {
		Entries []harEntry `json:"entries"`
	} `json:"log"`
}

type harEntry struct {
	Request struct {
		URL    string `json:"url"`
		Method string `json:"method"`
	} `json:"request"`
	Response struct {
		Status     int         `json:"status"`
		StatusText string      `json:"statusText"`
		Headers    []harHeader `json:"headers"`
		Content    struct {
			MimeType string `json:"mimeType"`
			Text     string `json:"text"`
			Encoding string `json:"encoding,omitempty"`
		} `json:"content"`
	} `json:"response"`
}

type harHeader struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// HarNotFound controls how requests that have no matching HAR entry are handled.
type HarNotFound string

const (
	HarNotFoundAbort    HarNotFound = "abort"
	HarNotFoundFallback HarNotFound = "fallback"
)

// RouteFromHAROptions configures BrowserContext.RouteFromHAR and Page.RouteFromHAR.
type RouteFromHAROptions struct {
	URL      *string
	NotFound HarNotFound
}

// RouteFromHAR reads a HAR file and serves matching requests from its recorded responses.
func (c *BrowserContext) RouteFromHAR(ctx context.Context, harPath string, opts ...*RouteFromHAROptions) error {
	return c.routeFromHAR(ctx, harPath, opts...)
}

// RouteFromHAR reads a HAR file and serves matching requests from its recorded responses
// for this page only.
func (p *Page) RouteFromHAR(ctx context.Context, harPath string, opts ...*RouteFromHAROptions) error {
	return routeFromHAR(ctx, harPath, func(pattern string, handler RouteHandler) error {
		return p.Route(ctx, pattern, handler)
	}, opts...)
}

func (c *BrowserContext) routeFromHAR(ctx context.Context, harPath string, opts ...*RouteFromHAROptions) error {
	return routeFromHAR(ctx, harPath, func(pattern string, handler RouteHandler) error {
		return c.Route(ctx, pattern, handler)
	}, opts...)
}

func routeFromHAR(ctx context.Context, harPath string, register func(string, RouteHandler) error, opts ...*RouteFromHAROptions) error {
	data, err := os.ReadFile(harPath)
	if err != nil {
		return fmt.Errorf("routeFromHAR: failed to read %q: %w", harPath, err)
	}
	var har harFile
	if err := json.Unmarshal(data, &har); err != nil {
		return fmt.Errorf("routeFromHAR: invalid HAR file: %w", err)
	}

	urlPattern := "**"
	notFound := HarNotFoundAbort
	if len(opts) > 0 && opts[0] != nil && opts[0].URL != nil {
		urlPattern = *opts[0].URL
	}
	if len(opts) > 0 && opts[0] != nil && opts[0].NotFound != "" {
		notFound = opts[0].NotFound
	}

	return register(urlPattern, func(route *Route) {
		callCtx, callCancel := context.WithTimeout(context.Background(), time.Duration(fallbackTimeoutMs*float64(time.Millisecond)))
		defer callCancel()
		req := route.Request()
		if req == nil {
			_ = route.Continue(callCtx, nil)
			return
		}
		reqURL := strings.TrimRight(req.URL(), "/")
		for i := range har.Log.Entries {
			entry := &har.Log.Entries[i]
			if strings.TrimRight(entry.Request.URL, "/") != reqURL {
				continue
			}
			if entry.Request.Method != "" && req.Method() != "" &&
				!strings.EqualFold(entry.Request.Method, req.Method()) {
				continue
			}
			headers := make(map[string]string, len(entry.Response.Headers))
			for _, h := range entry.Response.Headers {
				headers[h.Name] = h.Value
			}
			status := entry.Response.Status
			ct := entry.Response.Content.MimeType

			var bodyStr string
			var bodyBytes []byte
			if entry.Response.Content.Encoding == "base64" {
				decoded, decErr := base64.StdEncoding.DecodeString(entry.Response.Content.Text)
				if decErr == nil {
					bodyBytes = decoded
				}
			} else {
				bodyStr = entry.Response.Content.Text
			}

			fulfillOpts := &RouteFulfillOptions{
				Status:      &status,
				ContentType: &ct,
				Headers:     headers,
			}
			if bodyBytes != nil {
				fulfillOpts.BodyBytes = bodyBytes
			} else {
				fulfillOpts.Body = &bodyStr
			}
			if err := route.Fulfill(callCtx, fulfillOpts); err != nil {
				slog.Default().Warn("routeFromHAR: failed to fulfill request", "url", reqURL, "error", err)
			}
			return
		}
		if notFound == HarNotFoundFallback {
			if err := route.Continue(callCtx, nil); err != nil {
				slog.Default().Warn("routeFromHAR: failed to continue request", "url", reqURL, "error", err)
			}
			return
		}
		if err := route.Abort(callCtx); err != nil {
			slog.Default().Warn("routeFromHAR: failed to abort request", "url", reqURL, "error", err)
		}
	})
}
