package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ChandanGupta31/student-api/internal/config"
	"github.com/ChandanGupta31/student-api/internal/http/handlers/student"
	"github.com/ChandanGupta31/student-api/internal/storage/sqlite"
)

func main() {
	// load config
	cfg := config.MustLoad()

	// database setup
	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("Storage Initialize")

	// setup routers
	router := http.NewServeMux()
	// -> Creating End Point
	router.HandleFunc("POST /api/students", student.New(storage))

	// setup server
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	// here we are applying channel and goroutine so that server shutdown gracefully
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		slog.Info("Server ....")
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("Failed to Start Server")
		}
	}()

	<-done

	slog.Info("Shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Failed to Shutdown Server", slog.String("Error", err.Error()))
	}

	slog.Info("Server Shutdown Successfully")
}
