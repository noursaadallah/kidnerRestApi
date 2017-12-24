package models

type Pair struct {
	ID        string       `json:"id"`
	DocType   string       `json:"docType"`
	Donor     HealthRecord `json:"donor"`
	Recipient HealthRecord `json:"recipient"`
	Score     float64      `json:"score"`
	Match     bool         `json:"match"`
	Active    bool         `json:"active"`
	DrID      string       `json:"drId"`
}
