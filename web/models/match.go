package models

import "time"

type Match struct {
	ID           string     `json:"id"`
	DocType      string     `json:"docType"`
	MatchedPairs [][]string `json:"matchedPairs"`
	Approved     bool       `json:"approved"`
	EndorcingDr  string     `json:"endorcingDr"`
	DrSig        string     `json:"drSig"`
	CreateDate   time.Time  `json:"createDate"`
}
