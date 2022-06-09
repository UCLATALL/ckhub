package server

import "github.com/uclatall/ckhub/pkg/logging"

// Option is a generic interface of the server configuration option.
type Option interface {
	// Apply applies the option to the given server.
	Apply(srv *Server) error
}

// OptionFunc is an adapter to allow the use of ordinary functions as options.
type OptionFunc func(srv *Server) error

// Apply applies the option to the given server.
func (o OptionFunc) Apply(srv *Server) error {
	return o(srv)
}

// Logger creates a new option that sets the logger for the server.
func Logger(log logging.Logger) OptionFunc {
	return func(srv *Server) error {
		srv.log = log
		return nil
	}
}
