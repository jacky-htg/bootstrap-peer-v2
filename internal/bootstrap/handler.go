package bootstrap

import (
	"bootstrap-server/pkg"
	"bufio"
	"fmt"
	"net"

	"github.com/bytedance/sonic"
)

type MessageType string

const (
	NewPeerJoined MessageType = "new_peer_joined"
	ShutdownPeer  MessageType = "shutdown_peer"
)

type Message struct {
	Type MessageType `json:"type"`
	Data []byte      `json:"data"`
}

// handleConnection menangani koneksi dan menentukan jenis request.
func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	ip, _, err := net.SplitHostPort(conn.RemoteAddr().String())
	if err != nil {
		fmt.Println("Failed to get IP address:", err)
		return
	}

	if !s.rateLimit.Allow(ip) {
		fmt.Println("Rate limit exceeded for IP:", ip)
		return
	}

	reader := bufio.NewReader(conn)

	// Baca request dari client
	data, err := reader.ReadBytes('\n')
	if err != nil {
		fmt.Println("Error reading from connection:", err)
		return
	}

	var req pkg.Request
	if err := sonic.Unmarshal(data, &req); err != nil {
		fmt.Println("Invalid request format:", err)
		return
	}

	println(req.Type)
	// Handle sesuai tipe request
	switch req.Type {
	case "REGISTER":
		s.registerPeer(req.Payload)
	case "GET_PEERS":
		s.getAllPeers(conn)
	case "REMOVE":
		s.removePeer(req.Payload)
	default:
		fmt.Println("Invalid request type")
	}
}

// registerPeer menambahkan peer baru ke daftar.
func (s *Server) registerPeer(peer []byte) {
	var myPeer Peer
	err := sonic.Unmarshal(peer, &myPeer)
	if err != nil {
		fmt.Println("Failed to unmarshall peer:")
		return
	}

	fmt.Println(myPeer)

	success, _ := s.pm.RegisterPeer(myPeer)

	if !success {
		fmt.Println("Failed to register peer:", peer)
	} else {
		s.broadcastPeers(myPeer, "register")
	}
}

func (s *Server) getAllPeers(conn net.Conn) {
	peers := s.pm.GetAllPeers()

	data, err := sonic.Marshal(peers)
	if err != nil {
		fmt.Println("Error encoding peers:", err)
		return
	}
	conn.Write(append(data, '\n'))
}

func (s *Server) removePeer(peer []byte) {
	var curPeer Peer
	err := sonic.Unmarshal(peer, &curPeer)
	if err != nil {
		fmt.Println("Failed to unmarshall peer:")
		return
	}
	success, _ := s.pm.RemovePeer(curPeer.Address)

	if !success {
		fmt.Println("Failed to remove peer:", peer)
	} else {
		s.broadcastPeers(curPeer, "shutdown")
	}
}

func (s *Server) broadcastPeers(peer Peer, typeMsg string) {
	s.pm.mu.Lock()
	defer s.pm.mu.Unlock()

	for _, existingPeer := range s.pm.peers {
		if existingPeer.Address != peer.Address {
			if typeMsg == "shutdown" {
				go s.notifyShutdownPeer(existingPeer.Address, peer.Address)
			} else {
				go s.notifyNewPeer(existingPeer.Address, peer)
			}
		}
	}
}

func (s *Server) notifyNewPeer(existingPeer string, newPeer Peer) {
	conn, err := net.Dial("tcp", existingPeer)
	if err != nil {
		fmt.Println("Error notifying peer:", err)
		return
	}
	defer conn.Close()

	newPeerdata, err := sonic.Marshal(newPeer)
	if err != nil {
		fmt.Println("Failed to serialize new peer:", err)
		return
	}

	message := Message{
		Type: NewPeerJoined,
		Data: newPeerdata,
	}
	data, _ := sonic.Marshal(message)
	writer := bufio.NewWriter(conn)
	data = append(data, '\n')

	_, err = writer.WriteString(string(data))
	if err != nil {
		fmt.Println("Gagal mengirim request:", err)
		return
	}
	writer.Flush()
}

func (s *Server) notifyShutdownPeer(existingPeer, newPeerAddress string) {
	conn, err := net.Dial("tcp", existingPeer)
	if err != nil {
		fmt.Println("Error notifying peer:", err)
		return
	}
	defer conn.Close()

	message := Message{
		Type: ShutdownPeer,
		Data: []byte(newPeerAddress),
	}
	data, _ := sonic.Marshal(message)

	writer := bufio.NewWriter(conn)
	data = append(data, '\n')

	_, err = writer.WriteString(string(data))
	if err != nil {
		fmt.Println("Gagal mengirim request:", err)
		return
	}
	writer.Flush()
}
