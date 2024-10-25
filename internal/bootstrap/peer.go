package bootstrap

import (
	"fmt"
	"sync"

	"github.com/dgraph-io/badger/v4"
	"google.golang.org/protobuf/proto"
)

// PeerManager mengelola daftar peers dengan thread-safe.
type PeerManager struct {
	peers []Peer
	db    *badger.DB
	mu    sync.Mutex
}

// NewPeerManager membuat instance baru dari PeerManager.
func NewPeerManager(db *badger.DB) *PeerManager {
	return &PeerManager{
		db: db,
	}
}

// RegisterPeer menambahkan peer baru ke dalam daftar.
func (pm *PeerManager) RegisterPeer(peerAddress string) (bool, string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	for _, p := range pm.peers {
		if p.Address == peerAddress {
			return false, "Peer already registered"
		}
	}

	peer := Peer{Address: peerAddress}
	pm.peers = append(pm.peers, peer)

	data, err := proto.Marshal(&peer)
	if err != nil {
		return false, fmt.Sprintf("Failed to serialize peer: %v", err)
	}

	err = pm.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(peer.Address), data)
	})

	if err != nil {
		return false, fmt.Sprintf("Failed to store peer in DB: %v", err)
	}

	return true, "Peer registered successfully"
}

// GetAllPeers mengembalikan daftar semua peers.
func (pm *PeerManager) GetAllPeers() []Peer {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	return pm.peers
}

// RemovePeer menghapus peer dari daftar.
func (pm *PeerManager) RemovePeer(peerAddress string) (bool, string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	for i, p := range pm.peers {
		if p.Address == peerAddress {
			pm.peers = append(pm.peers[:i], pm.peers[i+1:]...)

			err := pm.db.Update(func(txn *badger.Txn) error {
				return txn.Delete([]byte(peerAddress))
			})

			if err != nil {
				return false, fmt.Sprintf("Failed to remove peer from DB: %v", err)
			}

			return true, "Peer removed successfully"
		}
	}

	return false, "Peer not found"
}

func (pm *PeerManager) loadPeers() error {
	pm.peers = []Peer{}
	// Iterasi seluruh record di Badger
	err := pm.db.View(func(txn *badger.Txn) error {
		iter := txn.NewIterator(badger.DefaultIteratorOptions)
		defer iter.Close()

		for iter.Rewind(); iter.Valid(); iter.Next() {
			item := iter.Item()
			err := item.Value(func(val []byte) error {
				var peer Peer
				if err := proto.Unmarshal(val, &peer); err != nil {
					return err
				}
				pm.peers = append(pm.peers, peer)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
