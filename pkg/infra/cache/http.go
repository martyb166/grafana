package cache

import (
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/grafana/grafana/pkg/infra/remotecache"
	"github.com/grafana/grafana/pkg/setting"
)

func registerRoutes(cfg *setting.Cfg, c *Cache) error {
	mux := http.NewServeMux()
	mux.Handle("/ring", c.lfc)
	mux.Handle("/kv", c.mlist)

	mux.HandleFunc("GET /cache/{key}", func(w http.ResponseWriter, r *http.Request) {
		key := r.PathValue("key")
		c.logger.Info("get cached item", "key", key)
		value, err := c.Get(r.Context(), key)
		if err != nil {
			if errors.Is(err, remotecache.ErrCacheItemNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			c.logger.Error("failed to get item", "err", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(&getResponse{Value: value})
	})

	mux.HandleFunc("POST /cache/internal", func(w http.ResponseWriter, r *http.Request) {
		c.logger.Info("set new item internal")
		var req setRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			c.logger.Error("failed to parse request internal", "err", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := c.Set(r.Context(), req.Key, req.Value, req.Expr); err != nil {
			c.logger.Error("failed to set item internal", "err", err)
			w.WriteHeader(http.StatusBadRequest)
			return

		}
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("POST /cache", func(w http.ResponseWriter, r *http.Request) {
		c.logger.Info("set new item")
		type request struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		}
		var req request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			c.logger.Error("failed to parse request", "err", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := c.Set(r.Context(), req.Key, []byte(req.Value), 1*time.Hour); err != nil {
			c.logger.Error("failed to set item", "err", err)
			w.WriteHeader(http.StatusBadRequest)
			return

		}
		w.WriteHeader(http.StatusOK)
	})

	listener, err := net.Listen("tcp", net.JoinHostPort(cfg.HTTPAddr, httpPort))
	if err != nil {
		return err
	}

	go func() {
		panic(http.Serve(listener, mux))
	}()

	return nil
}
