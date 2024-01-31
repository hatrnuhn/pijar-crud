package database

import "sync"

/*
Host: "localhost"
User: "root"
Password : ""
Database: "pijarcamp"
Produk : "nama_produk", "keterangan", "harga", "jumlah"
*/

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Host     string            `json:"host"`
	User     string            `json:"user"`
	Password string            `json:"password"`
	Database string            `json:"database"`
	Produks  map[string]Produk `json:"produks"`
}

type Produk struct {
	NamaProduk string `json:"nama_produk"`
	Keterangan string `json:"keterangan"`
	Harga      int    `json:"harga"`
	Jumlah     int    `json:"jumlah"`
}
