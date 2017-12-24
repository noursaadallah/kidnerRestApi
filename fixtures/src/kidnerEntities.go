// ============================================================================================================================
// This file contains the model :
// Structs and constants
// ============================================================================================================================

package main

import "time"

const (
	// indexes
	indexPair  = "pair~ID"
	indexMatch = "match~ID"

	// document types, i.e object types or classes
	docTypeHealthRecord = "healthRecord"
	docTypePair         = "pair"
	docTypeMatch        = "match"
	docTypeDoctor       = "doctor"

	// bool represented as string
	TRUE  = "true"
	FALSE = "false"

	//errors
	errNumArgs        = "Incorrect number of arguments: Expecting "
	errJsonMarshall   = "Marshall : Error encoding object to Json"
	errJsonUnmarshall = "Unmarshall : Error decoding Json to object"
	errPutState       = "Error Put State"
	errGetState       = "Error Get State"
)

type HealthRecord struct {
	DocType        string            `json:"docType"`
	Age            int               `json:"age"`
	BloodType      string            `json:"bloodType"`
	MedicalUrgency int               `json:"medicalUrgency"`
	HLAs           map[string]string `json:"HLAs"`
	PRA            int               `json:"PRA"`
	Eligible       bool              `json:"eligible"`
	Type           string            `json:"type"`
	Signature      string            `json:"signature"`
	CreateDate     time.Time         `json:"createDate"`
}

type Pair struct {
	ID        string       `json:"id"`
	DocType   string       `json:"docType"`
	Donor     HealthRecord `json:"donor"`
	Recipient HealthRecord `json:"recipient"`
	Score     float64      `json:"score"`
	Match     bool         `json:"match"`
	Active    bool         `json:"active"`
	DrID      string       `json:"drId"`
	DrSig     string       `json:"drSig"`
}

type Match struct {
	ID           string     `json:"id"`
	DocType      string     `json:"docType"`
	MatchedPairs [][]string `json:"matchedPairs"`
	Approved     bool       `json:"approved"`
	EndorcingDr  string     `json:"endorcingDr"`
	DrSig        string     `json:"drSig"`
	CreateDate   time.Time  `json:"createDate"`
}

type Doctor struct {
	ID        string `json:"id"`
	DocType   string `json:"docType"`
	Signature string `json:"signature"`
}
