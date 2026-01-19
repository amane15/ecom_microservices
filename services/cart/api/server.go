package main

import (
	"fmt"
	"log/slog"
	"net/http"
)

func (app *application) server() error {
	srv := &http.Server{
		Addr:     fmt.Sprintf(":%d", app.config.port),
		Handler:  app.routes(),
		ErrorLog: slog.NewLogLogger(app.logger.Handler(), slog.LevelError),
	}

	app.logger.Info("starting server", "address", srv.Addr)
	return srv.ListenAndServe()
}
