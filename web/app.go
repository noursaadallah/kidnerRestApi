package web

import (
	"fmt"
	"net/http"

	"github.com/noursaadallah/kidner/settings"
	"github.com/noursaadallah/kidner/web/controllers"
)

func Serve(app *controllers.Application) {
	fs := http.FileServer(http.Dir("web/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.HandleFunc("/listActivePairs.html", app.ListActivePairsHandler)
	http.HandleFunc("/createDoctor.html", app.CreateDoctorHandler)
	http.HandleFunc("/createPair.html", app.CreatePairHandler)
	http.HandleFunc("/getPair.html", app.GetPairHandler)
	http.HandleFunc("/getDoctor.html", app.GetDoctorHandler)
	http.HandleFunc("/updatePair.html", app.UpdatePairHandler)
	http.HandleFunc("/deactivatePair.html", app.DeactivatePairHandler)
	http.HandleFunc("/deletePair.html", app.DeletePairHandler)
	http.HandleFunc("/findPairedMatch.html", app.FindPairedMatchHandler)
	http.HandleFunc("/listMatches.html", app.ListMatchesHandler)
	http.HandleFunc("/approveMatch.html", app.ApproveMatchHandler)
	http.HandleFunc("/refuseMatch.html", app.RefuseMatchHandler)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/listActivePairs.html", http.StatusTemporaryRedirect)
	})

	var ws settings.WebSettings

	ws, err := settings.GetWebSettings()
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Error loading web server config - web server not starting ")
		return
	}

	address := ws.Address + ":" + ws.Port
	fmt.Println("Listening (http://" + address + "/) ...")
	http.ListenAndServe(address, nil)
}
