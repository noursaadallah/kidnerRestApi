package controllers

import (
	"encoding/json"
	"net/http"
	"sort"

	"github.com/cloudflare/cfssl/log"
	"github.com/noursaadallah/kidner/web/models"
)

func (app *Application) ListMatchesHandler(w http.ResponseWriter, r *http.Request) {
	data := &struct {
		Matches  []models.Match
		Success  bool
		Response bool
		Error    string
	}{
		Matches:  make([]models.Match, 0),
		Success:  false,
		Response: false,
		Error:    "",
	}

	matchesAsBytes, err := app.Fabric.ListMatches()
	if err != nil {
		log.Error(err.Error())
		//http.Error(w, "Unable to query the ID in the blockchain", 500)
		data.Error = "Unable to invoke function in the blockchain : " + renderError(err)
		renderTemplate(w, r, "listMatches.html", data)
		return
	}

	err = json.Unmarshal(matchesAsBytes, &data.Matches)
	if err != nil {
		log.Error(err.Error())
		//http.Error(w, "Get incorrect entity", 500)
		data.Error = "Error unmarshalling slice of Matches : " + renderError(err)
		renderTemplate(w, r, "listMatches.html", data)
		return
	}
	data.Success = true
	data.Response = true
	data.Matches = sortMatches(data.Matches)

	renderTemplate(w, r, "listMatches.html", data)
}

// ============================================================================================
// sorting matches by time.Time
// ============================================================================================
type sliceMatches []models.Match

func (p sliceMatches) Len() int {
	return len(p)
}

func (p sliceMatches) Less(i, j int) bool {
	return p[i].CreateDate.After(p[j].CreateDate)
}

func (p sliceMatches) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func sortMatches(matches []models.Match) []models.Match {
	sliceMatches := make(sliceMatches, 0)
	for _, m := range matches {
		sliceMatches = append(sliceMatches, m)
	}
	sort.Sort(sliceMatches) // the matches are now sorted
	sortedMatches := make([]models.Match, 0)
	for _, m := range sliceMatches {
		sortedMatches = append(sortedMatches, m)
	}
	return sortedMatches
}
