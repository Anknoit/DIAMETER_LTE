package api

import (
	"encoding/json"
	"net/http"
)

// simple management API used by FastAPI control plane for demo
func RegisterHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(`{"status":"ok"}`))
	})

	mux.HandleFunc("/peers", func(w http.ResponseWriter, r *http.Request) {
		// TODO: implement persistent peers store and runtime add/remove
		if r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode([]map[string]interface{}{})
			return
		}
		w.WriteHeader(http.StatusNotImplemented)
	})

	mux.HandleFunc("/simulate", func(w http.ResponseWriter, r *http.Request) {
		// For demo: accept posted JSON and echo a fake answer
		var body map[string]interface{}
		_ = json.NewDecoder(r.Body).Decode(&body)
		resp := map[string]interface{}{
			"result": "ok (simulated)", "request": body,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})
}
