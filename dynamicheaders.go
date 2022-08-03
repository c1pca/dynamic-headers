package dynamicheaders

import (
	"fmt"
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"go.uber.org/zap"
)

func init() {
	caddy.RegisterModule(DynamicHeaders{})
	httpcaddyfile.RegisterHandlerDirective("dynamic_headers", parseCaddyfile)
}

// DynamicHeaders Middleware implements an HTTP handler that writes
// headers dynamically.
type DynamicHeaders struct {
	// The file or stream to write to. Can be "stdout"
	// or "stderr".
	ToHeader   string `json:"to_header,omitempty"`
	FromHeader string `json:"from_header,omitempty"`
	logger     *zap.SugaredLogger
}

// CaddyModule returns the Caddy module information.
func (DynamicHeaders) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.dynamic_headers",
		New: func() caddy.Module { return new(DynamicHeaders) },
	}
}

// Provision implements caddy.Provisioner.
func (m *DynamicHeaders) Provision(ctx caddy.Context) error {
	m.logger = ctx.Logger(m).Sugar()
	if m.FromHeader != "" {
		m.logger.Debugf("FromHeader: %s", m.FromHeader)
	}
	if m.ToHeader != "" {
		m.logger.Debugf("ToHeader: %s", m.ToHeader)
	}
	return nil
}

// Validate implements caddy.Validator.
func (m *DynamicHeaders) Validate() error {
	if m.ToHeader == "" {
		return fmt.Errorf("provide to_header key to set the copied value")
	}
	if m.ToHeader == "" {
		return fmt.Errorf("provide from_header key to copy its value")
	}
	return nil
}

// ServeHTTP implements caddyhttp.MiddlewareHandler.
func (m DynamicHeaders) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	m.logger.Debugf("Value from %s is copied to %s header", m.FromHeader, m.ToHeader)
	value := r.Header.Get(m.FromHeader)

	if value == "" {
		m.logger.Errorf("header %s has no value", m.FromHeader)
	} else {
		m.logger.Debugf("header %s has value: %s", m.FromHeader, value)
		w.Header().Add(m.ToHeader, value)
	}

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
		case "to_header":
			m.ToHeader = value
		case "from_header":
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
