package jupyter

// Config represents a configuration of the jupyter server.
type Config struct {
	Address string `json:"url" yaml:"url"`
	Token   string `json:"token" yaml:"token"`
}

// Apply applies the given configuration to the client.
func (cfg Config) Apply(client *Client) error {
	client.http.SetBaseURL(cfg.Address)
	client.token = cfg.Token
	return nil
}
