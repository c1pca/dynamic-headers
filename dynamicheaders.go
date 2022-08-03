package dynamicheaders

import (
	"fmt"
	"io"
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"go.uber.org/zap"
)

func init() {
	caddy.RegisterModule(DynamicHeaders{})
	httpcaddyfile.RegisterHandlerDirective("dymanic_header", parseCaddyfile)
}

// DynamicHeaders Middleware implements an HTTP handler that writes
// writes headers dynamically.
type DynamicHeaders struct {
	// The file or stream to write to. Can be "stdout"
	// or "stderr".
	ToHeader   string `json:"to_header,omitempty"`
	FromHeader string `json:"from_header,omitempty"`
	log        *zap.Logger
	w          io.Writer
}

// CaddyModule returns the Caddy module information.
func (DynamicHeaders) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.dymanic_header",
		New: func() caddy.Module { return new(DynamicHeaders) },
	}
}

// Provision implements caddy.Provisioner.
func (m *DynamicHeaders) Provision(ctx caddy.Context) error {
	if m.FromHeader != "" {

	}
	if m.ToHeader != "" {

	}
	return nil
}

// Validate implements caddy.Validator.
func (m *DynamicHeaders) Validate() error {
	if m.w == nil {
		return fmt.Errorf("no writer")
	}
	return nil
}

// ServeHTTP implements caddyhttp.MiddlewareHandler.
func (m DynamicHeaders) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	m.w.Write([]byte(r.RemoteAddr))
	return next.ServeHTTP(w, r)
}

// UnmarshalCaddyfile implements caddyfile.Unmarshaler.
func (m *DynamicHeaders) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		key := d.Val()
		var value string
		if !d.Args(&value) {
			continue
		}

		switch key {
		case "redis_url":
			m.ToHeader = value
		case "path_prefix":
			m.FromHeader = value

		default:
			return fmt.Errorf("unknown key %s", key)
		}
	}
	return nil
}

// parseCaddyfile unmarshals tokens from h into a new Middleware.
func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var m DynamicHeaders
	err := m.UnmarshalCaddyfile(h.Dispenser)
	return m, err
}

// Interface guards
var (
	_ caddy.Provisioner           = (*DynamicHeaders)(nil)
	_ caddy.Validator             = (*DynamicHeaders)(nil)
	_ caddyhttp.MiddlewareHandler = (*DynamicHeaders)(nil)
	_ caddyfile.Unmarshaler       = (*DynamicHeaders)(nil)
)
