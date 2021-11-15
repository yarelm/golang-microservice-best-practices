package api

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/heptiolabs/healthcheck"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type Server struct {
	router *mux.Router
	addr   string
}

func NewServer(addr string) *Server {
	return &Server{
		router: mux.NewRouter(),
		addr:   addr,
	}
}

func (s *Server) Serve(ctx context.Context, health healthcheck.Handler) error {
	s.router.HandleFunc("/", s.rootHandler).Methods(http.MethodGet)

	srv := &http.Server{
		Addr:    s.addr,
		Handler: s.router,
	}

	health.AddReadinessCheck("HTTP", func() error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			return nil
		}
	})

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("failed to server listen and serve")
		}
	}()

	<-ctx.Done()

	log.Debug().Msg("server is shutting down...")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctxShutDown); err != nil {
		return errors.Wrap(err, "failed to shutdown server")
	}

	return nil
}
