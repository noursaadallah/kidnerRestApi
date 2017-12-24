package web

import (
	"fmt"
	"net/http"

	"github.com/chainhero/heroes-service/web/controllers"
)

func Serve(app *controllers.Application) {
	fs := http.FileServer(http.Dir("web/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.HandleFunc("/listActivePairs.html", app.ListActivePairsHandler)
	http.HandleFunc("/request.html", app.RequestHandler)
	http.HandleFunc("/createDoctor.html", app.CreateDoctorHandler)
	http.HandleFunc("/createPair.html", app.CreatePairHandler)
	http.HandleFunc("/getPair.html", app.GetPairHandler)
	http.HandleFunc("/getDoctor.html", app.GetDoctorHandler)
	http.HandleFunc("/updatePair.html", app.UpdatePairHandler)
	http.HandleFunc("/deactivatePair.html", app.DeactivatePairHandler)
	http.HandleFunc("/deletePair.html", app.DeletePairHandler)
	http.HandleFunc("/findPairedMatch.html", app.FindPairedMatchHandler)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/listActivePairs.html", http.StatusTemporaryRedirect)
	})

	fmt.Println("Listening (http://localhost:3000/) ...")
	http.ListenAndServe(":3000", nil)
}
