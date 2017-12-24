package controllers

import (
	"net/http"

	"github.com/cloudflare/cfssl/log"
)

func (app *Application) RefuseMatchHandler(w http.ResponseWriter, r *http.Request) {
	data := &struct {
		TxID     string
		Success  bool
		Response bool
		MatchID  string
		Error    string
	}{
		TxID:     "",
		Success:  false,
		Response: false,
		MatchID:  "",
		Error:    "",
	}

	// request came from listMatches.html
	if r.FormValue("hiddenMatchIDSent") == "true" {
		data.MatchID = r.FormValue("hiddenMatchID")
	}

	// request came from refuseMatch.html
	if r.FormValue("submitted") == "true" {
		matchID := r.FormValue("MatchID")
		drID := r.FormValue("DrID")
		drSig := r.FormValue("DrSig")
		var param []string
		param = append(param, drID, drSig, matchID)
		txID, err := app.Fabric.RefuseMatch(param)
		if err != nil {
			log.Error(err.Error())
			//http.Error(w, "Unable to invoke refuseMatch in the blockchain", 500)
			data.Error = "Unable to invoke function in the blockchain : " + renderError(err)
			renderTemplate(w, r, "refuseMatch.html", data)
			return
		}

		data.TxID = txID
		data.Success = true
		data.Response = true
	}
	renderTemplate(w, r, "refuseMatch.html", data)
}
