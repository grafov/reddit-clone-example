package main

import (
	"net/http"

	"redditclone/internal/config"

	"github.com/grafov/kiwi"
)

func httpServe(h http.Handler) *http.Server {
	var srv = &http.Server{
		Handler:      h,
		Addr:         config.App.ListenAt(),
		WriteTimeout: config.App.WriteTimeout,
		ReadTimeout:  config.App.ReadTimeout,
		IdleTimeout:  config.App.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			kiwi.Log("info", "http server stopped", "err", err)
		} else {
			kiwi.Log("info", "http server stopped")
		}
		if r := recover(); r != nil {
			panic(r)
		}
	}()

	return srv
}
