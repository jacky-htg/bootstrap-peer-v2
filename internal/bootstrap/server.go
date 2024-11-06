package bootstrap

import (
	"bootstrap-server/internal/db"
	"bootstrap-server/pkg"
	"fmt"
	"net"
	"time"
)

// Server menyimpan data peer dan menangani koneksi.
type Server struct {
	port      string
	pm        *PeerManager
	rateLimit *pkg.RateLimiter
}

// NewServer membuat instance baru server.
func NewServer(port, dbPath string) (*Server, error) {
	dbConn, err := db.InitDB(dbPath)
	if err != nil {
		return nil, err
	}

	return &Server{
		port:      port,
		pm:        NewPeerManager(dbConn),
		rateLimit: pkg.NewRateLimiter(100, time.Minute),
	}, nil
}

// ListenAndServe memulai server untuk menerima koneksi dari client.
func (s *Server) ListenAndServe() error {
	if err := s.pm.loadPeers(); err != nil {
		return fmt.Errorf("gagal membaca file peers: %v", err)
	}

	ln, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %v", s.port, err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection:", err)
			continue
		}
		go s.handleConnection(conn)
	}
}
