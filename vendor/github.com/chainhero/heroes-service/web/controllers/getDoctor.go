package controllers

import (
	"net/http"

	"github.com/cloudflare/cfssl/log"
)

func (app *Application) GetDoctorHandler(w http.ResponseWriter, r *http.Request) {
	data := &struct {
		Doctor   string
		Success  bool
		Response bool
	}{
		Doctor:   "",
		Success:  false,
		Response: false,
	}

	if r.FormValue("submitted") == "true" {
		ID := r.FormValue("ID")
		valAsBytes, err := app.Fabric.GetDoctor(ID)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, "Unable to query the ID in the blockchain", 500)
		}

		// err = json.Unmarshal(valAsBytes, &data.Doctor)
		// if err != nil {
		// 	log.Error(err.Error())
		// 	http.Error(w, "error unmarshalling Doctor entity", 500)
		// }
		data.Doctor = string(valAsBytes)
		data.Success = true
		data.Response = true
	}

	renderTemplate(w, r, "getDoctor.html", data)
}
