package controllers

import (
	"net/http"
)

func (app *Application) CreateDoctorHandler(w http.ResponseWriter, r *http.Request) {
	data := &struct {
		DoctorId string
		Success  bool
		Response bool
	}{
		DoctorId: "",
		Success:  false,
		Response: false,
	}
	if r.FormValue("submitted") == "true" {
		signature := r.FormValue("signature")
		drId, err := app.Fabric.CreateDoctor(signature)
		if err != nil {
			http.Error(w, "Unable to create Doctor in the blockchain : "+err.Error(), 500)
		}
		data.DoctorId = drId
		data.Success = true
		data.Response = true
	}
	renderTemplate(w, r, "createDoctor.html", data)
}
