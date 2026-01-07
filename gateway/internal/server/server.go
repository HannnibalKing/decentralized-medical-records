package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/example/med/gateway/internal/config"
	"github.com/example/med/gateway/internal/middleware"
	"github.com/example/med/gateway/internal/revocation"
	"github.com/example/med/gateway/internal/storage"
)

// Server bundles HTTP handlers.
type Server struct {
	mux       *http.ServeMux
	revCache  *revocation.Cache
	revClient *revocation.Client
	ipfs      storage.IPFSClient
	arweave   storage.ArweaveClient
}

// New creates the HTTP server with routes.
func New(cfg config.Config) *Server {
	mux := http.NewServeMux()
	srv := &Server{
		mux:       mux,
		revCache:  revocation.NewCache(cfg.RevocationCacheTTL),
		revClient: revocation.NewClient(cfg.RevocationRPC),
		ipfs:      storage.NewIPFS(cfg.IPFSURL),
		arweave:   storage.NewArweave(cfg.ArweaveURL),
	}

	mux.Handle("/v1/validate-capability", middleware.DPoP(http.HandlerFunc(srv.handleValidateCapability)))
	mux.Handle("/v1/fetch-record", middleware.DPoP(http.HandlerFunc(srv.handleFetchRecord)))
	mux.Handle("/v1/breakglass/activate", middleware.DPoP(http.HandlerFunc(srv.handleBreakGlassActivate)))
	mux.Handle("/v1/revoke", middleware.DPoP(http.HandlerFunc(srv.handleRevoke)))
	mux.Handle("/v1/attestations/", middleware.DPoP(http.HandlerFunc(srv.handleGetAttestation)))

	return srv
}

// Handler returns the HTTP handler.
func (s *Server) Handler() http.Handler { return s.mux }

func (s *Server) handleValidateCapability(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var req struct {
		Token     string `json:"token"`
		DPoP      string `json:"dpop"`
		RevHandle string `json:"rev_handle"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if req.RevHandle == "" {
		writeError(w, http.StatusBadRequest, "rev_handle required")
		return
	}
	ctx := r.Context()
	revoked, err := s.checkRevocation(ctx, req.RevHandle)
	if err != nil {
		writeError(w, http.StatusBadGateway, "revocation lookup failed")
		return
	}
	revState := "active"
	if revoked {
		revState = "revoked"
	}
	resp := map[string]interface{}{
		"valid":     !revoked,
		"rev_state": revState,
		"exp":       time.Now().Add(15 * time.Minute).Unix(),
		"scope":     "placeholder",
		"policy":    "policy-ER-v3",
		"reason":    "",
	}
	writeJSON(w, http.StatusOK, resp)
}

func (s *Server) handleFetchRecord(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var req struct {
		Token     string `json:"token"`
		DPoP      string `json:"dpop"`
		CID       string `json:"cid"`
		RevHandle string `json:"rev_handle"`
		ArweaveTX string `json:"arweave_tx"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if req.CID == "" {
		writeError(w, http.StatusBadRequest, "cid required")
		return
	}
	if req.RevHandle == "" {
		writeError(w, http.StatusBadRequest, "rev_handle required")
		return
	}
	ctx := r.Context()
	revoked, err := s.checkRevocation(ctx, req.RevHandle)
	if err != nil {
		writeError(w, http.StatusBadGateway, "revocation lookup failed")
		return
	}
	if revoked {
		writeError(w, http.StatusForbidden, "revoked")
		return
	}

	data, err := s.ipfs.Get(ctx, req.CID)
	if err != nil {
		// fallback to arweave if provided
		if req.ArweaveTX != "" {
			data, err = s.arweave.Get(ctx, req.ArweaveTX)
		}
	}
	if err != nil {
		writeError(w, http.StatusBadGateway, "storage fetch failed")
		return
	}
	resp := map[string]interface{}{
		"ciphertext":  data,
		"aad":         "", // placeholder AAD; populate when available
		"wrapped_key": "", // placeholder; real HPKE unwrap not implemented yet
		"cid":         req.CID,
	}
	writeJSON(w, http.StatusOK, resp)
}

func (s *Server) handleBreakGlassActivate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	// TODO: verify multi-sig policy, mint single-use capability, write rev handle on-chain.
	resp := map[string]interface{}{
		"capability":  "",
		"rev_handle":  "",
		"ttl_seconds": 1800,
	}
	writeJSON(w, http.StatusOK, resp)
}

func (s *Server) handleRevoke(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var req struct {
		RevHandle string `json:"rev_handle"`
		Reason    string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if req.RevHandle == "" {
		writeError(w, http.StatusBadRequest, "rev_handle required")
		return
	}
	ctx := r.Context()
	if err := s.revClient.Revoke(ctx, req.RevHandle, req.Reason); err != nil {
		writeError(w, http.StatusBadGateway, "revoke failed")
		return
	}
	s.revCache.Set(req.RevHandle, true)
	resp := map[string]interface{}{
		"tx_hash": "", // placeholder; populate from RPC response if available
	}
	writeJSON(w, http.StatusOK, resp)
}

func (s *Server) handleGetAttestation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	// TODO: query attestation registry.
	resp := map[string]interface{}{
		"cid_md":  "",
		"agg_sig": "",
		"status":  "",
		"ts":      time.Now().Unix(),
	}
	writeJSON(w, http.StatusOK, resp)
}

func writeJSON(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Printf("json encode error: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func (s *Server) checkRevocation(ctx context.Context, handle string) (bool, error) {
	if handle == "" {
		return false, nil
	}
	if revoked, ok := s.revCache.Get(handle); ok {
		return revoked, nil
	}
	revoked, err := s.revClient.Lookup(ctx, handle)
	if err != nil {
		return false, err
	}
	s.revCache.Set(handle, revoked)
	return revoked, nil
}
