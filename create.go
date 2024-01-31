package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/hatrnuhn/pijar-crud/internal/database"
	"github.com/hatrnuhn/pijar-crud/internal/utils"
)

func (cfg *Config) handleCreate(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	dat, err := io.ReadAll(r.Body)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "couldn't read request")
		return
	}

	req := database.Produk{}
	err = json.Unmarshal(dat, &req)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "couldn't unmarshal request")
		return
	}

	ps, err := cfg.db.GetProduks()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "couldn't get users")
		return
	}

	for _, p := range ps {
		if req.NamaProduk == p.NamaProduk {
			utils.RespondWithError(w, http.StatusConflict, "produk already exists")
			return
		}
	}

	newP, err := cfg.db.CreateProduk(string(dat))
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "couldn't create user")
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, newP)
}
