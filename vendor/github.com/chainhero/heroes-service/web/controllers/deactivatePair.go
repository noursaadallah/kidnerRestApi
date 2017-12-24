package controllers

import (
	"net/http"

	"github.com/cloudflare/cfssl/log"
)

func (app *Application) DeactivatePairHandler(w http.ResponseWriter, r *http.Request) {
	data := &struct {
		PairID   string
		Success  bool
		Response bool
	}{
		PairID:   "nil", // PairID != "nil" means Pair is not empty
		Success:  false,
		Response: false,
	}

	// request came from getPair.html or listPairs.html => deactivate the pair
	if r.FormValue("hiddenPairIDSent") == "true" {
		pairID := r.FormValue("hiddenPairID")
		_pairID, err := app.Fabric.DeactivatePair(pairID)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, "Unable to query the ID in the blockchain", 500)
		}

		data.PairID = _pairID
		data.Success = true
		data.Response = true
	}
	renderTemplate(w, r, "deactivatePair.html", data)
}
