package httpserver

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/yourname/go-netintel/internal/events"
	"github.com/yourname/go-netintel/internal/metrics"
	"github.com/yourname/go-netintel/internal/queue"
)

type Server struct {
	mux *http.ServeMux
	srv *http.Server
	pub queue.Publisher
}

func New(addr string, pub queue.Publisher) *Server {
	mux := http.NewServeMux()
	s := &Server{mux: mux, pub: pub}
	mux.HandleFunc("/healthz", s.health)
	mux.Handle("/metrics", metrics.Handler())
	mux.HandleFunc("/v1/collect", s.collectDomain)
	s.srv = &http.Server{Addr: addr, Handler: s.withMetrics(mux)}
	return s
}

func (s *Server) withMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := &statusWriter{ResponseWriter: w, status: 200}
		start := time.Now()
		next.ServeHTTP(ww, r)
		metrics.Requests.WithLabelValues(r.URL.Path, r.Method, http.StatusText(ww.status)).Inc()
		_ = start
	})
}

type statusWriter struct{ http.ResponseWriter; status int }
func (w *statusWriter) WriteHeader(code int) { w.status = code; w.ResponseWriter.WriteHeader(code) }

func (s *Server) health(w http.ResponseWriter, _ *http.Request) { w.Write([]byte("ok")) }

func (s *Server) collectDomain(w http.ResponseWriter, r *http.Request) {
	domain := r.URL.Query().Get("domain")
	if domain == "" { http.Error(w, "missing domain", http.StatusBadRequest); return }
	evt := events.New("dns", domain, map[string]any{"note":"stub"})
	if err := s.pub.Publish(r.Context(), evt); err != nil {
		log.Error().Err(err).Msg("publish failed")
		http.Error(w, "publish failed", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type","application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"status":"queued", "id": evt.ID})
}

func (s *Server) Start(ctx context.Context) error {
	go func() {
		log.Info().Msgf("http listening on %s", s.srv.Addr)
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("http server error")
		}
	}()
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	select { case <-ctx.Done(): case <-sig: }
	ctxShut, cancel := context.WithTimeout(context.Background(), 5*time.Second); defer cancel()
	return s.srv.Shutdown(ctxShut)
}
