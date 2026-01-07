//go:build grpc

package server

import (
	"context"
	"net"
	"time"

	gen "github.com/example/med/gateway/gen"
	"github.com/example/med/gateway/internal/config"
	"github.com/example/med/gateway/internal/revocation"
	"github.com/example/med/gateway/internal/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RPCServer implements the gRPC Gateway service defined in gateway.proto.
type RPCServer struct {
	gen.UnimplementedGatewayServer
	srv *Server
}

// NewRPCServer wraps the HTTP Server dependencies for gRPC handlers.
func NewRPCServer(cfg config.Config) *RPCServer {
	return &RPCServer{
		srv: &Server{
			revCache:  revocation.NewCache(cfg.RevocationCacheTTL),
			revClient: revocation.NewClient(cfg.RevocationRPC),
			ipfs:      storage.NewIPFS(cfg.IPFSURL),
			arweave:   storage.NewArweave(cfg.ArweaveURL),
		},
	}
}

// StartGRPC registers the service and begins serving on cfg.GRPCAddr.
func StartGRPC(cfg config.Config, s *grpc.Server, impl gen.GatewayServer) (net.Listener, error) {
	gen.RegisterGatewayServer(s, impl)
	lis, err := net.Listen("tcp", cfg.GRPCAddr)
	if err != nil {
		return nil, err
	}
	go func() { _ = s.Serve(lis) }()
	return lis, nil
}

func (r *RPCServer) ValidateCapability(ctx context.Context, req *gen.ValidateCapabilityRequest) (*gen.ValidateCapabilityResponse, error) {
	if req == nil || req.Capability == nil || req.Capability.Token == "" || req.Capability.RevHandle == "" {
		return nil, status.Error(codes.InvalidArgument, "token and rev_handle required")
	}
	revoked, err := r.srv.checkRevocation(ctx, req.Capability.RevHandle)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "revocation lookup failed: %v", err)
	}
	revState := "active"
	if revoked {
		revState = "revoked"
	}
	return &gen.ValidateCapabilityResponse{
		Valid:    !revoked,
		RevState: revState,
		Exp:      time.Now().Add(15 * time.Minute).Unix(),
		Scope:    "placeholder",
		Policy:   "policy-ER-v3",
		Reason:   "",
	}, nil
}

func (r *RPCServer) FetchRecord(ctx context.Context, req *gen.FetchRecordRequest) (*gen.FetchRecordResponse, error) {
	if req == nil || req.Capability == nil || req.Cid == "" || req.Capability.RevHandle == "" {
		return nil, status.Error(codes.InvalidArgument, "cid and rev_handle required")
	}
	revoked, err := r.srv.checkRevocation(ctx, req.Capability.RevHandle)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "revocation lookup failed: %v", err)
	}
	if revoked {
		return nil, status.Error(codes.PermissionDenied, "revoked")
	}
	data, err := r.srv.ipfs.Get(ctx, req.Cid)
	if err != nil && req.ArweaveTx != "" {
		data, err = r.srv.arweave.Get(ctx, req.ArweaveTx)
	}
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "storage fetch failed: %v", err)
	}
	return &gen.FetchRecordResponse{
		Ciphertext: data,
		Aad:        nil,
		WrappedKey: nil,
		Cid:        req.Cid,
	}, nil
}

func (r *RPCServer) BreakGlassActivate(ctx context.Context, req *gen.BreakGlassActivateRequest) (*gen.BreakGlassActivateResponse, error) {
	return nil, status.Error(codes.Unimplemented, "break-glass not yet implemented")
}

func (r *RPCServer) Revoke(ctx context.Context, req *gen.RevokeRequest) (*gen.RevokeResponse, error) {
	if req == nil || req.RevHandle == "" {
		return nil, status.Error(codes.InvalidArgument, "rev_handle required")
	}
	if err := r.srv.revClient.Revoke(ctx, req.RevHandle, req.Reason); err != nil {
		return nil, status.Errorf(codes.Internal, "revoke failed: %v", err)
	}
	r.srv.revCache.Set(req.RevHandle, true)
	return &gen.RevokeResponse{TxHash: ""}, nil
}

func (r *RPCServer) GetAttestation(ctx context.Context, req *gen.GetAttestationRequest) (*gen.GetAttestationResponse, error) {
	return nil, status.Error(codes.Unimplemented, "attestation lookup not yet implemented")
}
