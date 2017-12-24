package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/cloudflare/cfssl/log"
	"github.com/noursaadallah/kidner/web/models"
)

func (app *Application) GetDoctorHandler(w http.ResponseWriter, r *http.Request) {
	data := &struct {
		Doctor   models.Doctor
		Success  bool
		Response bool
		Error    string
	}{
		Doctor:   *new(models.Doctor),
		Success:  false,
		Response: false,
		Error:    "",
	}

	if r.FormValue("submitted") == "true" {
		ID := r.FormValue("ID")
		valAsBytes, err := app.Fabric.GetDoctor(ID)
		if err != nil {
			log.Error(err.Error())
			//http.Error(w, "Unable to query the ID in the blockchain", 500)
			data.Error = "Unable to invoke function in the blockchain : " + renderError(err)
			renderTemplate(w, r, "getDoctor.html", data)
			return
		}

		err = json.Unmarshal(valAsBytes, &data.Doctor)
		if err != nil {
			log.Error(err.Error())
			//http.Error(w, "error unmarshalling Doctor entity", 500)
			data.Error = "Error unmarshalling Doctor entity" + renderError(err)
			renderTemplate(w, r, "getDoctor.html", data)
			return
		}
		data.Success = true
		data.Response = true
	}

	renderTemplate(w, r, "getDoctor.html", data)
}
