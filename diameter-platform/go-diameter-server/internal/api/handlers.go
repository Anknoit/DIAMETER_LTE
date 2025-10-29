package api

import "net/http"

// RegisterHandlers registers the management http handlers.
func RegisterHandlers(mux *http.ServeMux) {
    mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("ok"))
    })
    mux.HandleFunc("/peers", func(w http.ResponseWriter, r *http.Request) {
        // TODO: list/add peers
        w.Write([]byte("peers endpoint"))
    })
}
