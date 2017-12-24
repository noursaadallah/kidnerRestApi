package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/chainhero/heroes-service/web/models"
	"github.com/cloudflare/cfssl/log"
)

func (app *Application) UpdatePairHandler(w http.ResponseWriter, r *http.Request) {
	data := &struct {
		PairID   string
		Pair     models.Pair
		Success  bool
		Response bool
		GotPair  bool
	}{
		PairID:   "nil", // PairID != "nil" means Pair is not empty
		Pair:     *new(models.Pair),
		Success:  false,
		Response: false,
		GotPair:  false,
	}

	// request came from getPair.html or listPairs.html => get pair and display data
	if r.FormValue("hiddenPairIDSent") == "true" {
		// inject initial pair data in the updatePair.html
		ID := r.FormValue("hiddenPairID")
		valAsBytes, err := app.Fabric.GetPair(ID)
		if err != nil {
			http.Error(w, "Unable to query the ID in the blockchain", 500)
		}

		err = json.Unmarshal(valAsBytes, &data.Pair)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, "Get incorrect entity", 500)
		}
		data.GotPair = true
		data.PairID = data.Pair.ID
	}
	// => the pair is updated // request came from updatePair.html
	if r.FormValue("submitted") == "true" {
		// parse and retrieve form elements
		pairID := r.FormValue("PairID")
		// donor
		ageD := r.FormValue("AgeD")
		bloodTypeD := r.FormValue("BloodTypeD")
		medicalUrgencyD := "0"
		PRAD := r.FormValue("PRAD")
		eligibilityD := r.FormValue("EligibilityD")
		typeD := "donor"

		// donor HLAs :
		hlaA1D := r.FormValue("A1D")
		hlaA2D := r.FormValue("A2D")
		hlaB1D := r.FormValue("B1D")
		hlaB2D := r.FormValue("B2D")
		hlaDR1D := r.FormValue("DR1D")
		hlaDR2D := r.FormValue("DR2D")
		hlaDQ1D := r.FormValue("DQ1D")
		hlaDQ2D := r.FormValue("DQ2D")
		HLAsD := hlaA1D + "," + hlaA2D + "," + hlaB1D + "," + hlaB2D + "," + hlaDR1D + "," + hlaDR2D + "," + hlaDQ1D + "," + hlaDQ2D

		// recipient
		ageR := r.FormValue("AgeR")
		bloodTypeR := r.FormValue("BloodTypeR")
		medicalUrgencyR := r.FormValue("MedicalUrgencyR")
		PRAR := r.FormValue("PRAR")
		eligibilityR := r.FormValue("EligibilityR")
		typeR := "recipient"

		// recipient HLAs :
		hlaA1R := r.FormValue("A1R")
		hlaA2R := r.FormValue("A2R")
		hlaB1R := r.FormValue("B1R")
		hlaB2R := r.FormValue("B2R")
		hlaDR1R := r.FormValue("DR1R")
		hlaDR2R := r.FormValue("DR2R")
		hlaDQ1R := r.FormValue("DQ1R")
		hlaDQ2R := r.FormValue("DQ2R")
		HLAsR := hlaA1R + "," + hlaA2R + "," + hlaB1R + "," + hlaB2R + "," + hlaDR1R + "," + hlaDR2R + "," + hlaDQ1R + "," + hlaDQ2R

		// Parameters of updatePair
		var args []string
		args = []string{pairID, ageR, bloodTypeR, medicalUrgencyR, HLAsR, PRAR, eligibilityR, typeR,
			ageD, bloodTypeD, medicalUrgencyD, HLAsD, PRAD, eligibilityD, typeD}

		txid, err := app.Fabric.UpdatePair(args)
		if err != nil {
			http.Error(w, "Unable to write state in the blockchain"+err.Error(), 500)
		}
		data.PairID = txid
		data.Success = true
		data.Response = true
	}
	renderTemplate(w, r, "updatePair.html", data)
}
