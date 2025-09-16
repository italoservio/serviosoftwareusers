package main

import (
	"context"
	"errors"
	"github.com/italoservio/serviosoftwareusers/pkg/env"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/italoservio/serviosoftwareusers/internal/api"
	"github.com/italoservio/serviosoftwareusers/internal/deps"

	"github.com/gorilla/mux"
)

func main() {
	envVars := env.Load()
	container := deps.NewContainer(envVars)
	router := mux.NewRouter()
	router.MethodNotAllowedHandler = http.HandlerFunc(api.MethodNotAllowed)

	api.RegisterInfraRoutes(router)
	api.RegisterUsersRoutes(router, container)

	wg := sync.WaitGroup{}

	svr := &http.Server{Addr: ":8080", Handler: router}
	wg.Go(func() {
		println("server listening on port :8080")

		err := svr.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}

		println("server stopped")
	})

	exitCode := 0
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT)
	signal.Notify(sigCh, syscall.SIGTERM)

	wg.Go(func() {
		sig := <-sigCh
		println("received signal:", strings.ToUpper(sig.String()))

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := svr.Shutdown(ctx); err != nil {
			println("error shutting down server:", err.Error())
			exitCode = 1
		}

		if err := container.DB.Disconnect(); err != nil {
			println("error disconnecting from database:", err.Error())
			exitCode = 1
		}
	})

	wg.Wait()
	os.Exit(exitCode)
}
