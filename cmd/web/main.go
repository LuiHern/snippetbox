package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	// config flags
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	// logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// router or "mux"
	mux := http.NewServeMux()

	fileserver := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileserver))

	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /snippet/view/{id}", snippetView)
	mux.HandleFunc("GET /snippet/create", snippetCreate)
	mux.HandleFunc("POST /snippet/create", snippetCreatePost)

	logger.Info("Starting server", slog.Any("addr", *addr))
	// logger.Info("Starting server", "addr", *addr)  OTHER OPTION
	// logger.Info("Starting server", slog.String("addr", *addr))  OTHER OPTION

	err := http.ListenAndServe(*addr, mux)
	logger.Error(err.Error())
	os.Exit(1)
}
