package config

import (
	"os"
	"time"
)

// Config holds runtime configuration.
type Config struct {
	Addr               string
	GRPCAddr           string
	RevocationRPC      string
	RevocationCacheTTL time.Duration
	IPFSURL            string
	ArweaveURL         string
}

// FromEnv reads configuration with sane defaults.
func FromEnv() Config {
	addr := os.Getenv("GATEWAY_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	grpcAddr := os.Getenv("GATEWAY_GRPC_ADDR")
	if grpcAddr == "" {
		grpcAddr = ":8090"
	}

	revRPC := os.Getenv("REVOCATION_RPC")
	if revRPC == "" {
		revRPC = "http://localhost:8545" // placeholder RPC for on-chain revocation reads
	}

	ipfsURL := os.Getenv("IPFS_URL")
	if ipfsURL == "" {
		ipfsURL = "http://localhost:5001" // Kubo API
	}

	arweaveURL := os.Getenv("ARWEAVE_URL")
	if arweaveURL == "" {
		arweaveURL = "https://arweave.net"
	}

	ttl := 30 * time.Second
	if env := os.Getenv("REVOCATION_CACHE_TTL"); env != "" {
		if d, err := time.ParseDuration(env); err == nil {
			ttl = d
		}
	}

	return Config{
		Addr:               addr,
		GRPCAddr:           grpcAddr,
		RevocationRPC:      revRPC,
		RevocationCacheTTL: ttl,
		IPFSURL:            ipfsURL,
		ArweaveURL:         arweaveURL,
	}
}
