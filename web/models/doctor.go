package models

type Doctor struct {
	ID        string `json:"id"`
	DocType   string `json:"docType"`
	Signature string `json:"signature"`
}
