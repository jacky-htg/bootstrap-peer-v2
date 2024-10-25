package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

type Request struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

// Fungsi utama client untuk memanggil perintah registrasi, get peers, dan remove peer
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <command>")
		fmt.Println("Commands: register | get_peers | remove")
		return
	}

	command := os.Args[1]

	if command != "get_peers" && len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <command> <peer_address>")
		fmt.Println("Commands: register | remove")
		fmt.Println("Example: go run main.go register localhost:3000")
		return
	}

	var peerAddress string
	if len(os.Args) >= 3 {
		peerAddress = os.Args[2]
	}

	// Membuka koneksi TCP ke bootstrap server
	conn, err := net.DialTimeout("tcp", "localhost:4000", 10*time.Second)
	if err != nil {
		fmt.Println("Gagal terhubung ke server:", err)
		return
	}
	defer conn.Close()

	switch strings.ToLower(command) {
	case "register":
		sendRequest(conn, "REGISTER", peerAddress)
	case "get_peers":
		sendRequest(conn, "GET_PEERS", peerAddress)
	case "remove":
		sendRequest(conn, "REMOVE", peerAddress)
	default:
		fmt.Println("Perintah tidak dikenali:", command)
	}
}

// Fungsi untuk mengirim request dan membaca response dari server
func sendRequest(conn net.Conn, command, peerAddress string) {
	reader := bufio.NewReader(conn)

	req := Request{Type: command, Payload: peerAddress}
	data, err := json.Marshal(req)
	if err != nil {
		fmt.Println("gagal membuat request JSON: ", err)
		return
	}
	writer := bufio.NewWriter(conn)
	data = append(data, '\n')

	_, err = writer.WriteString(string(data))
	if err != nil {
		fmt.Println("Gagal mengirim request:", err)
		return
	}
	writer.Flush()

	// Membaca response dari server
	response, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Gagal membaca response:", err)
		return
	}

	fmt.Println("Response dari server:", strings.TrimSpace(response))
}
