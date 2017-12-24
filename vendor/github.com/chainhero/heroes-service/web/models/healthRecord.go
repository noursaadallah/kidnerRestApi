package models

import "time"

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
