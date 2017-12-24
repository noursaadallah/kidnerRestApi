package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/cloudflare/cfssl/log"
	"github.com/noursaadallah/kidner/web/models"
)

func (app *Application) GetPairHandler(w http.ResponseWriter, r *http.Request) {
	data := &struct {
		Pair     models.Pair
		Success  bool
		Response bool
		Error    string
	}{
		Pair:     *new(models.Pair),
		Success:  false,
		Response: false,
		Error:    "",
	}

	if r.FormValue("submitted") == "true" {
		ID := r.FormValue("ID")
		valAsBytes, err := app.Fabric.GetPair(ID)
		if err != nil {
			//http.Error(w, "Unable to query the ID in the blockchain", 500)
			log.Error(err.Error())
			data.Error = "Unable to invoke function in the blockchain : " + renderError(err)
			renderTemplate(w, r, "getPair.html", data)
			return
		}

		err = json.Unmarshal(valAsBytes, &data.Pair)
		if err != nil {
			log.Error(err.Error())
			//http.Error(w, "Get incorrect entity or result is empty", 500)
			data.Error = "Error unmarshalling Pair entity : " + renderError(err)
			renderTemplate(w, r, "getPair.html", data)
			return
		}
		data.Success = true
		data.Response = true
	}

	renderTemplate(w, r, "getPair.html", data)
}
