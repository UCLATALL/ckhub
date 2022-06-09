package jupyter

// Option is a generic interface of the jupyter configuration option.
type Option interface {
	// Apply applies the option to the given jupyter.
	Apply(client *Client) error
}

// OptionFunc is an adapter to allow the use of ordinary functions as options.
type OptionFunc func(client *Client) error

// Apply applies the option to the given jupyter.
func (o OptionFunc) Apply(client *Client) error {
	return o(client)
}
