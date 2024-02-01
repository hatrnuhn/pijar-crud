package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi"
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

	if req.NamaProduk == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "produk has to have a name")
		return
	}

	ps, err := cfg.db.GetProduks()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("couldn't get users: %s", err.Error()))
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

func (cfg *Config) handleRead(w http.ResponseWriter, r *http.Request) {
	ps, err := cfg.db.GetProduks()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "couldn't get produks")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, ps)
}

func (cfg *Config) handleUpdate(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	param := chi.URLParam(r, "produkName")

	dat, err := io.ReadAll(r.Body)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "couldn't read request")
		return
	}

	ps, err := cfg.db.GetProduks()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "couldn't read produks in database")
		return
	}

	found := false
	for _, p := range ps {
		if p.NamaProduk == param {
			found = true
			break
		}
	}

	if !found {
		utils.RespondWithError(w, http.StatusBadRequest, "produk_name doesn't exist")
		return
	}

	req := database.Produk{}
	err = json.Unmarshal(dat, &req)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "couldn't unmarshal request")
		return
	}

	newP, err := cfg.db.UpdateProduk(&req)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "couldn't update the produk")
		return
	}

	utils.RespondWithJSON(w, http.StatusAccepted, newP)
}

func (cfg *Config) handleDel(w http.ResponseWriter, r *http.Request) {
	pName := chi.URLParam(r, "produkName")

	ps, err := cfg.db.GetProduks()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "couldn't read produks in database")
		return
	}

	found := false
	for _, p := range ps {
		if p.NamaProduk == pName {
			found = true
			break
		}
	}

	if !found {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("produk_name: %s doesn't exist", pName))
		return
	}

	deletedP, err := cfg.db.DeleteProduk(pName)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("couldn't delete produk: %s", err.Error()))
		return
	}

	utils.RespondWithJSON(w, http.StatusAccepted, deletedP)
}
