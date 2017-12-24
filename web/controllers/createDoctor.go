package controllers

import (
	"net/http"

	"github.com/cloudflare/cfssl/log"
)

func (app *Application) CreateDoctorHandler(w http.ResponseWriter, r *http.Request) {
	data := &struct {
		DoctorId string
		Success  bool
		Response bool
		Error    string
	}{
		DoctorId: "",
		Success:  false,
		Response: false,
		Error:    "",
	}
	if r.FormValue("submitted") == "true" {
		signature := r.FormValue("signature")
		drId, err := app.Fabric.CreateDoctor(signature)
		if err != nil {
			//http.Error(w, "Unable to create Doctor in the blockchain : "+err.Error(), 500)
			log.Error(err.Error())
			data.Error = "Unable to invoke function in the blockchain : " + renderError(err)
			renderTemplate(w, r, "createDoctor.html", data)
			return
		}
		data.DoctorId = drId
		data.Success = true
		data.Response = true
	}
	renderTemplate(w, r, "createDoctor.html", data)
}
