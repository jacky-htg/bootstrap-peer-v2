# Bootstrap Peer Server

Bootstrap Peer Server adalah aplikasi yang bertugas untuk mengelola daftar peer dalam jaringan P2P. Server ini memungkinkan peer untuk mendaftar, mengambil daftar peer yang terdaftar, dan menghapus peer yang ingin menonaktifkan diri. Data peer disimpan dalam format prototobuff dalam database badgerDB untuk kemudahan akses dan pengelolaan.

## Fitur

- **Register Peer**: Peer baru dapat mendaftar ke server.
- **Get All Peers**: Mengambil daftar semua peer yang terdaftar.
- **Remove Peer**: Menghapus peer dari daftar ketika peer ingin menonaktifkan diri.
- Membaca file json ketika pertama kali server dinyalakan
- Menulis ke file json ketika server di-shutdown.

## Struktur Folder

```console
.
├── cmd
│   └── main.go
├── data
│   └── peers.db
├── internal
│   ├── bootstrap
│   │   ├── handler.go
│   │   ├── peer.go
│   │   ├── peer.pb.go
│   │   ├── peer.proto
│   │   └── server.go
│   └── db
│       └── badger.go
├── pkg
│   └── models.go
└── go.mod
```


## Dependensi

Pastikan untuk menginstal dependensi berikut:

- [Go](https://golang.org/doc/install) (versi 1.23 atau lebih baru)
- [Sonic](https://github.com/bytedance/sonic) untuk serialisasi/deserialisai json.
- [Badger](github.com/dgraph-io/badger/v4) untuk menyimpan database key-value.
- [Protobuf](google.golang.org/protobuf) untuk proses serialisasi dan deserialisasi proto

## Instalasi

1. Clone repository ini:

    ```bash
    git clone https://github.com/jacky-htg/bootstrap-peer-v1.git
    cd bootstrap-peer-v1
    ```

2. Instal dependensi:

    ```bash
    go mod tidy
    ```

## Menjalankan Server

Untuk menjalankan server bootstrap, gunakan perintah berikut:

```bash
go run cmd/main.go
```

Server akan mendengarkan pada port 4000.

## Testing
Untuk pengetesean mendaftar peer baru, ambil daftar peer, atau menghapus peer, Anda dapat menggunakan client yang disediakan.

**Mendaftar Peer**

Untuk mengetes mendaftar sebagai peer, gunakan perintah berikut:

```bash
go run client.go register <peer_address> 
```
Ganti <peer_address> dengan alamat peer yang ingin didaftarkan, misalnya `go run client.go register localhost:3000`

**Mengambil Daftar Peer**

Untuk mengambil semua peer yang terdaftar, gunakan perintah berikut:

```bash
go run client.go get_peers
```

**Menghapus Peer**

Untuk menghapus peer dari daftar, gunakan perintah berikut:

```bash
go run cmd/client/main.go remove <peer_address>
```

Ganti <peer_address> dengan alamat peer yang ingin dihapus.


## Lisensi
GNU GPL License. Silakan lihat file LICENSE untuk rincian lebih lanjut.