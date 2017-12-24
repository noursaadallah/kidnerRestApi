package controllers

import (
	"net/http"

	"github.com/cloudflare/cfssl/log"
)

func (app *Application) CreatePairHandler(w http.ResponseWriter, r *http.Request) {
	data := &struct {
		TransactionId string
		Success       bool
		Response      bool
		Error         string
	}{
		TransactionId: "",
		Success:       false,
		Response:      false,
		Error:         "",
	}
	if r.FormValue("submitted") == "true" {
		// parse and retrieve form elements
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

		DrID := r.FormValue("DrID")
		DrSig := r.FormValue("DrSig")
		RecipientSig := r.FormValue("RecipientSig")
		DonorSig := r.FormValue("DonorSig")

		// Parameters of createPair
		var args []string
		args = []string{ageD, bloodTypeD, medicalUrgencyD, HLAsD, PRAD, eligibilityD, typeD,
			ageR, bloodTypeR, medicalUrgencyR, HLAsR, PRAR, eligibilityR, typeR,
			DrID, RecipientSig, DonorSig, DrSig}

		txid, err := app.Fabric.CreatePair(args)
		if err != nil {
			log.Error(err.Error())
			//http.Error(w, "Unable to write state in the blockchain"+err.Error(), 500)
			data.Error = "Unable to invoke function in the blockchain : " + renderError(err)
			renderTemplate(w, r, "createPair.html", data)
			return
		}
		data.TransactionId = txid
		data.Success = true
		data.Response = true

		// After the pair is created look for a match cycle automatically
		_, err = app.Fabric.FindMatchCycle()
		if err != nil {
			http.Error(w, "Error while looking for a match cycle "+err.Error(), 500)
		}
	}
	renderTemplate(w, r, "createPair.html", data)
}
