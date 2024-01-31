package database

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

func NewDB(path string) (*DB, error) {
	db := DB{
		path: path,
		mux:  &sync.RWMutex{},
	}

	err := db.ensureDB()
	if err != nil {
		return &DB{}, err
	}

	return &db, nil
}

func (db *DB) ensureDB() error {
	_, err := os.Stat(db.path)

	if os.IsNotExist(err) {
		file, err := os.Create(db.path)
		if err != nil {
			return err
		}
		defer file.Close()
	} else if err != nil {
		return err
	}

	return nil
}

func (db *DB) loadDB() (DBStructure, error) {
	// loads database.json
	body, err := os.ReadFile(db.path)
	if err != nil {
		return DBStructure{}, err
	}

	dbS := DBStructure{}

	// unmarshals data from file if it's not empty
	if len(body) != 0 {
		return DBStructure{}, errors.New("unmarshal loadDB error")
	} else {
		dbS.Host = "localhost"
		dbS.Database = "pijarcamp"
		dbS.Produks = make(map[string]Produk)
	}

	return dbS, nil
}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbS DBStructure) error {
	dat, err := json.Marshal(dbS)
	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, dat, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) CreateProduk(body string) (Produk, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbS, err := db.loadDB()
	if err != nil {
		return Produk{}, err
	}

	req := Produk{}

	err = json.Unmarshal([]byte(body), &req)
	if err != nil {
		return Produk{}, errors.New("req body unmarshal error")
	}

	p := Produk{
		NamaProduk: req.NamaProduk,
		Keterangan: req.Keterangan,
		Harga:      req.Harga,
		Jumlah:     req.Jumlah,
	}

	dbS.Produks[p.NamaProduk] = p
	err = db.writeDB(dbS)
	if err != nil {
		return Produk{}, err
	}

	return p, nil
}

func (db *DB) GetProduks() ([]Produk, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	dbS, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	ps := make([]Produk, len(dbS.Produks))
	for _, p := range dbS.Produks {
		ps = append(ps, p)
	}

	return ps, nil
}

func (db *DB) UpdateProduk(produk *Produk) (Produk, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbS, err := db.loadDB()
	if err != nil {
		return Produk{}, err
	}

	dbS.Produks[produk.NamaProduk] = *produk

	err = db.writeDB(dbS)
	if err != nil {
		return Produk{}, errors.New("UpdateProduk: writeDB error")
	}

	return *produk, nil
}

func (db *DB) DeleteProduk(pName string) (Produk, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbS, err := db.loadDB()
	if err != nil {
		return Produk{}, err
	}

	deleted := dbS.Produks[pName]

	delete(dbS.Produks, pName)

	err = db.writeDB(dbS)
	if err != nil {
		return Produk{}, errors.New("DeleteProduk: writeDB error")
	}

	return deleted, nil
}
