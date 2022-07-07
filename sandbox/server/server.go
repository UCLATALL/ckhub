package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/multierr"

	"github.com/uclatall/ckhub/pkg/logging"
	"github.com/uclatall/ckhub/sandbox"
)

// Server implements a sandbox management server.
type Server struct {
	log     logging.Logger
	manager *sandbox.Manager
	mux     chi.Router

	addr string
}

// NewServer creates a new sandbox management server with the given options.
func NewServer(manager *sandbox.Manager, options ...Option) (*Server, error) {
	server := &Server{
		log:     logging.NopLogger(),
		addr:    ":8080",
		manager: manager,
		mux:     chi.NewRouter(),
	}

	server.mux.Get("/healthz", server.HealthCheck)
	server.mux.Post("/api/v1/execute/{kernel}", server.Execute)

	errs := make([]error, len(options))
	for i, option := range options {
		errs[i] = option.Apply(server)
	}

	err := multierr.Combine(errs...)
	if err != nil {
		return nil, err
	}
	return server, nil
}

// Run kicks off the server. It blocks the execution and interrupts when
// the context is canceled, or any error occurs.
func (srv *Server) Run(ctx context.Context) error {
	log := srv.log

	lis, rerr := net.Listen("tcp", srv.addr)
	if rerr != nil {
		return fmt.Errorf("failed to create listener: %w", rerr)
	}
	defer func() { _ = lis.Close() }()
	log.Debug("starting server", logging.Stringer("address", lis.Addr()))

	server := &http.Server{
		Addr:    srv.addr,
		Handler: srv.mux,
	}

	done := make(chan struct{})
	go func() {
		rerr = server.Serve(lis)
		if errors.Is(rerr, http.ErrServerClosed) {
			rerr = nil
		}
		if rerr != nil {
			log.Error("server interrupted", logging.Error(rerr))
		}
		close(done)
	}()

	select {
	case <-done:
		return rerr
	case <-ctx.Done():
		log.Debug("server shutdown", logging.Error(ctx.Err()))
	}

	log = log.Hooks(logging.Span())

	serr := server.Shutdown(context.Background())
	if serr != nil {
		log.Error("server interrupted", logging.Error(serr))
		return fmt.Errorf("failed to shutdown server: %w", serr)
	}

	return nil
}

// Execute executes the code in the sandbox.
func (srv *Server) Execute(w http.ResponseWriter, req *http.Request) {
	log := srv.log.Hooks(logging.Span())

	id := uuid.New()

	kernel := strings.ToLower(chi.URLParam(req, "kernel"))
	log.Fields(logging.String("kernel", kernel), logging.Stringer("trace", id))

	body, err := io.ReadAll(req.Body)
	defer func() { _ = req.Body.Close() }()
	if err != nil {
		log.Error("failed to read request", logging.Error(err))
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	result, err := srv.manager.ExecuteSnippet(req.Context(), &sandbox.Snippet{
		ID:     id,
		Kernel: kernel,
		Source: string(body),
	})
	if err != nil {
		if errors.Is(err, sandbox.ErrKernelNotFound) {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		log.Error("failed to execute snippet", logging.Error(err))
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	log.Debug("execution complete")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Error("failed to write response", logging.Error(err))
	}
}

// HealthCheck returns a health check status of the service.
func (srv *Server) HealthCheck(w http.ResponseWriter, req *http.Request) {
	_ = req.Body.Close()
	if req.URL.Query().Has("ready") && srv.manager.Kernels() == 0 {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func writeError(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(err.Error())
}
