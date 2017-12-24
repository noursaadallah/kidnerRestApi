package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/chainhero/heroes-service/web/models"
	"github.com/cloudflare/cfssl/log"
)

func (app *Application) ListActivePairsHandler(w http.ResponseWriter, r *http.Request) {
	data := &struct {
		Pairs    []models.Pair
		Success  bool
		Response bool
	}{
		Pairs:    make([]models.Pair, 0),
		Success:  false,
		Response: false,
	}

	pairsAsBytes, err := app.Fabric.ListActivePairs()
	if err != nil {
		log.Error(err.Error())
		http.Error(w, "Unable to query the ID in the blockchain", 500)
	}

	err = json.Unmarshal(pairsAsBytes, &data.Pairs)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, "Get incorrect entity", 500)
	}
	data.Success = true
	data.Response = true

	renderTemplate(w, r, "listActivePairs.html", data)
}
