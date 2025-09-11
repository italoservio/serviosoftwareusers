package main

import (
	"context"
	"errors"
	"github.com/italoservio/serviosoftwareusers/internal/api"
	"github.com/italoservio/serviosoftwareusers/internal/deps"
	"github.com/italoservio/serviosoftwareusers/pkg/db"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	d, err := db.NewDB(os.Getenv("MONGODB_URI"))

	if err != nil {
		panic(err)
	}

	c := deps.NewContainer(d)
	r := mux.NewRouter()
	r.MethodNotAllowedHandler = http.HandlerFunc(api.MethodNotAllowed)

	api.RegisterInfraRoutes(r)
	api.RegisterUsersRoutes(r, c)

	wg := sync.WaitGroup{}

	svr := &http.Server{Addr: ":8080", Handler: r}
	wg.Go(func() {
		println("server listening on port :8080")
		if err := svr.ListenAndServe(); err != nil && !errors.Is(http.ErrServerClosed, err) {
			panic(err)
		}
		println("server stopped")
	})

	exitCode := 0
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	signal.Notify(sigCh, os.Kill)

	wg.Go(func() {
		sig := <-sigCh
		println("received signal:", strings.ToUpper(sig.String()))

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := svr.Shutdown(ctx); err != nil {
			println("error shutting down server:", err.Error())
			exitCode = 1
		}

		if err := d.Disconnect(); err != nil {
			println("error disconnecting from database:", err.Error())
			exitCode = 1
		}
	})

	wg.Wait()
	os.Exit(exitCode)
}
