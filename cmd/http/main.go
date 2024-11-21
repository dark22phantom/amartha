package main

import (
	"amartha/cmd/initialize"
	"amartha/config"
	"context"
	"fmt"
	"net/http"
	"os"

	handlerHttp "amartha/handler/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	os.Setenv("TZ", "Asia/Jakarta")

	var ctx = context.Background()

	cfg, err := config.New()
	if err != nil {
		return
	}
	_ = cfg

	port := ":5500"

	init, err := initialize.Initialize(ctx, cfg)

	httpHandler := handlerHttp.New(init.Loan)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/", func(r chi.Router) {
		httpHandler.RouteHandler(r)
	})

	fmt.Println("Listening on " + port)
	http.ListenAndServe(port, r)
}
