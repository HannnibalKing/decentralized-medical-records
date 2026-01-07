package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/revocations/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"revoked": false})
	})
	mux.HandleFunc("/revocations", func(w http.ResponseWriter, r *http.Request) {
		// accept POST revoke
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"ok": true})
	})
	addr := ":8181"
	log.Printf("mock revocation RPC listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
