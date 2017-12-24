package controllers

import (
	"net/http"
)

func (app *Application) FindPairedMatchHandler(w http.ResponseWriter, r *http.Request) {
	data := &struct {
		TxID     string
		Success  bool
		Response bool
	}{
		TxID:     "nil", // TxID != "nil" means Pair is not empty
		Success:  false,
		Response: false,
	}

	// request came from getPair.html or listPairs.html => Find paired match for that pair
	if r.FormValue("hiddenPairIDSent") == "true" {
		ID := r.FormValue("hiddenPairID")
		txId, err := app.Fabric.FindPairedMatch(ID)
		if err != nil {
			http.Error(w, "Unable to invoke FindPairedMatch(ID)", 500)
		}
		data.Success = true
		data.Response = true
		data.TxID = txId
	}
	renderTemplate(w, r, "findPairedMatch.html", data)
}
