package server

// Config represents a configuration of the the sandbox management server.
type Config struct {
	Address string `json:"http" yaml:"http"`
}

// Apply applies the configuration to the given server.
func (cfg Config) Apply(srv *Server) error {
	srv.addr = cfg.Address
	return nil
}
