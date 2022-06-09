package sandbox

import (
	"errors"
	"fmt"

	"go.uber.org/multierr"

	"github.com/uclatall/ckhub/pkg/jupyter"
)

// Config represents a sandbox configuration.
type Config struct {
	Kernels []struct {
		Name    string         `json:"name" yaml:"name"`
		Init    string         `json:"init,omitempty" yaml:"init,omitempty"`
		Jupyter jupyter.Config `json:"jupyter" yaml:"jupyter"`
		Kernel  string         `json:"kernel" yaml:"kernel"`
		Limit   uint           `json:"limit" yaml:"limit"`
	} `json:"kernels" yaml:"kernels"`
}

// ErrDuplicateKernel is returned when a kernel with the same name is already
// exists.
var ErrDuplicateKernel = errors.New("duplicate kernel name")

// Apply applies the given configuration to the manager.
func (cfg Config) Apply(manager *Manager) error {
	errs := make([]error, len(cfg.Kernels))

	for i, config := range cfg.Kernels {
		if _, ok := manager.kernels[config.Name]; ok {
			errs[i] = fmt.Errorf("%s: %w", config.Name, ErrDuplicateKernel)
			continue
		}

		client, err := jupyter.NewClient(config.Jupyter)
		if err != nil {
			errs[i] = fmt.Errorf("failed to create client for %s: %w", config.Name, err)
			continue
		}

		manager.kernels[config.Name] = &Kernel{
			client: client,
			name:   config.Name,
			init:   config.Init,
			limit:  int64(config.Limit),
		}
	}

	err := multierr.Combine(errs...)
	if err != nil {
		return err
	}

	return nil
}
