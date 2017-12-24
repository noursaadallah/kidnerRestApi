package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/chainhero/heroes-service/web/models"
	"github.com/cloudflare/cfssl/log"
)

func (app *Application) GetPairHandler(w http.ResponseWriter, r *http.Request) {
	data := &struct {
		Pair     models.Pair
		Success  bool
		Response bool
	}{
		Pair:     *new(models.Pair),
		Success:  false,
		Response: false,
	}

	if r.FormValue("submitted") == "true" {
		ID := r.FormValue("ID")
		valAsBytes, err := app.Fabric.GetPair(ID)
		if err != nil {
			http.Error(w, "Unable to query the ID in the blockchain", 500)
		}

		err = json.Unmarshal(valAsBytes, &data.Pair)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, "Get incorrect entity or result is empty", 500)
		}
		data.Success = true
		data.Response = true
	}

	renderTemplate(w, r, "getPair.html", data)
}
