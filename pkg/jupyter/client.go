package jupyter

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/go-resty/resty/v2"
	"go.uber.org/multierr"
)

// Client implements a client for jupyter server.
type Client struct {
	http  *resty.Client
	token string
}

// NewClient creates a new jupyter client with the given options.
func NewClient(options ...Option) (*Client, error) {
	client := &Client{
		http:  resty.New().SetBaseURL("http://jupyter:8888").SetAuthScheme("token"),
		token: "",
	}

	errs := make([]error, len(options))
	for i, option := range options {
		errs[i] = option.Apply(client)
	}

	err := multierr.Combine(errs...)
	if err != nil {
		return nil, err
	}
	return client, nil
}

type createRequest struct {
	Name string `json:"name"`
}

// CreateKernel creates a new jupyter kernel with the given name.
func (client *Client) CreateKernel(ctx context.Context, name string) (*Kernel, error) {
	uri, err := url.Parse(client.http.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid server url: %w", err)
	}

	var result Response[Kernel]

	res, err := client.http.R().
		EnableTrace().
		SetContext(ctx).
		SetAuthToken(client.token).
		SetError(&result.Error).
		SetResult(&result.Result).
		SetBody(createRequest{Name: name}).
		Post("/api/kernels")
	if err != nil {
		return nil, fmt.Errorf("failed to process request: %w", err)
	}
	if !res.IsSuccess() {
		return nil, fmt.Errorf("invalid server response: %w", result.Error)
	}

	kernel := &result.Result

	kernel.Address = res.Request.TraceInfo().RemoteAddr

	uri.Scheme = "ws"
	if strings.EqualFold(uri.Scheme, "https") {
		uri.Scheme = "wss"
	}
	uri.Path = fmt.Sprintf("%s/api/kernels/%s/channels", uri.Path, kernel.ID)
	query := uri.Query()
	query.Add("token", client.token)
	uri.RawQuery = query.Encode()

	kernel.ChanURL = uri.String()

	return kernel, nil
}

// RemoveKernel removes the jupyter kernel with the given identifier.
func (client *Client) RemoveKernel(ctx context.Context, kernel *Kernel) error {
	var result Response[json.RawMessage]

	uri, err := url.Parse(client.http.BaseURL)
	if err != nil {
		return fmt.Errorf("invalid server url: %w", err)
	}
	hostname := uri.Hostname()
	uri.Host = kernel.Address.String()

	res, err := resty.New().SetBaseURL(uri.String()).SetAuthScheme("token").R().
		SetContext(ctx).
		SetAuthToken(client.token).
		SetError(&result.Error).
		SetHeader("host", hostname).
		SetResult(&result.Result).
		SetPathParam("id", kernel.ID.String()).
		Delete("/api/kernels/{id}")
	if err != nil {
		return fmt.Errorf("failed to process request: %w", err)
	}
	if !res.IsSuccess() {
		return fmt.Errorf("invalid server response: %w", result.Error)
	}

	return nil
}

// Response describes a jupyter server response.
type Response[T any] struct {
	Error  Error
	Result T
}

// Error represents an error returned by the jupyter server.
type Error struct {
	Message string `json:"message"`
}

// Error returns the error message.
func (err Error) Error() string {
	return err.Message
}
