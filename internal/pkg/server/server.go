package server

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/zerowidth/gh-shorthand/internal/pkg/config"
)

// Run runs the gh-shorthand RPC server on the configured unix socket path
func Run() {
	path, err := homedir.Expand("~/.gh-shorthand.yml")
	if err != nil {
		log.Fatal("couldn't load config", err)
	}
	cfg, err := config.LoadFromFile(path)
	if err != nil {
		log.Fatal("couldn't load config", err)
	}

	if len(cfg.SocketPath) == 0 {
		log.Fatal("no socket_path configured in ~/.gh-shorthand.yml")
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	rpc := NewRPCHandler(cfg)
	rpc.Mount(r)

	server := &http.Server{
		Handler:      r,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	sock, err := net.Listen("unix", cfg.SocketPath)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		os.Remove(cfg.SocketPath)
	}()

	go func() {
		log.Printf("server started on %s\n", cfg.SocketPath)
		if err := server.Serve(sock); err != nil {
			log.Fatal("server error", err)
		}
	}()

	<-sig

	log.Printf("shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("server shutdown error", err)
	}
}
