package main

import "net/http"

func NewMux() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		// for escaping error, trash return value
		_, _ = w.Write([]byte(`{"status": "ok"}`))
	})
	return mux
}
