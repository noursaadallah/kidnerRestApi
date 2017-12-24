// ============================================================================================================================
// This file contains the business logic behind Kidner software
// i.e how to find a match , what is a match, etc.
// ============================================================================================================================
package main

import (
	"math"
	"strings"
	"time"
	//ECS "github.com/noursaadallah/elementary-cycles-search"
)

// ============================================================================================================================
// Match : RECIPIENT calls the method ! ( the logic : we're trying to provide a match for R so the search starts by R)
// checks whether 2 healthRecords match or don't
// ============================================================================================================================
func (recipient HealthRecord) Match(donor HealthRecord) *Pair {
	result := new(Pair)
	result.DocType = docTypePair
	result.Donor = donor
	result.Recipient = recipient
	result.Active = true

	/************************ Surgery Eligibility or unavailability ***********/
	if !recipient.Eligible || !donor.Eligible {
		result.Match = false
		result.Score = 0
		return result
	}

	/*******************    Age Compatibility    **********************/

	ageScore := ageScore(donor.Age, recipient.Age)

	/*******************    Blood Compatibility    **********************/
	bloodScore := 0.0
	if strings.EqualFold(recipient.BloodType, "AB") || strings.EqualFold(donor.BloodType, "O") || strings.EqualFold(donor.BloodType, recipient.BloodType) {
		bloodScore = 1.0
	}

	/*******************    HLA Compatibility    **********************/
	count := 0
	for k := range donor.HLAs {
		if donor.HLAs[k] == recipient.HLAs[k] {
			count++
		}
	}
	HLAscore := float64(count) / 8.0

	/*******************    If Not compatible END   **********************/
	if HLAscore < 0.75 || ageScore == 0 || bloodScore == 0 {
		result.Match = false
		result.Score = 0
		return result
	}

	result.Score = HLAscore + ageScore + bloodScore
	// must be > 1 + 0.1 + 0.75 to be compatible
	if result.Score >= 1.85 {
		result.Match = true
	}

	/*******************   If compatible handle Priorities    **********************/
	/***************************	Seniority	***********************************/
	todayStamp := time.Now().Unix()
	epoch := time.Date(1970, 01, 01, 00, 00, 00, 00, time.UTC).Unix()
	recipientStamp := recipient.CreateDate.Unix()

	seniority := float64(todayStamp-recipientStamp) / float64(todayStamp-epoch)

	/*******************************      %PRA      **********************************/
	PRAscore := 0.0
	if recipient.PRA >= 85 {
		PRAscore = 2.0
	} else if recipient.PRA >= 70 {
		PRAscore = 1.0
	} else if recipient.PRA >= 50 {
		PRAscore = 0.5
	}

	/************************************* Final Score *************************************/
	result.Score += seniority + float64(recipient.MedicalUrgency) + PRAscore
	return result
}

// ============================================================================================================================
// calculates score based on difference of ages
// ============================================================================================================================
func ageScore(donorAge int, recipientAge int) float64 {
	// âge des priorités : 2 mineurs ou 2 seniors
	if (donorAge <= 18 && recipientAge <= 18) || (donorAge >= 65 && recipientAge >= 65) {
		return 1.0
	}

	ageDiff := int(math.Abs(float64(donorAge - recipientAge)))
	if (ageDiff >= 0) && (ageDiff <= 5) {
		return 1.0
	} else if (ageDiff >= 6) && (ageDiff <= 10) {
		return 0.8
	} else if (ageDiff >= 11) && (ageDiff <= 20) {
		return 0.6
	} else if (ageDiff >= 21) && (ageDiff <= 30) {
		return 0.3
	} else if (ageDiff >= 31) && (ageDiff <= 40) {
		return 0.1
	}
	return 0.0
}

// ============================================================================================================================
// find list of pairs that have a cross match with the specified pair
// ============================================================================================================================
func (pair Pair) FindPairedMatches(activePairs []Pair) []Pair {
	var result []Pair
	for _, v := range activePairs {
		tmp := pair.Recipient.Match(v.Donor)
		_tmp := v.Recipient.Match(pair.Donor)
		if tmp.Match && _tmp.Match {
			result = append(result, v)
		}
	}

	return result
}

// ============================================================================================================================
// find the pair that has the best paired match with the specified pair (highest score)
// ============================================================================================================================
func (pair Pair) FindBestPairedMatch(activePairs []Pair) *Pair {
	result := pair.FindPairedMatches(activePairs)
	if len(result) == 0 {
		return nil
	}

	maxVal := 0.0 // keep best score = sum of 2 pairs (D1,R2) and (D2,R1)
	maxInd := 0   // keep index of best pair (D2,R2)
	for k, v := range result {
		// (D1,R1) and (D2,R2) => test scores of (D1,R2) and (D2,R1)
		tmp := pair.Recipient.Match(v.Donor)
		_tmp := v.Recipient.Match(pair.Donor)
		if tmp.Score+_tmp.Score > maxVal {
			maxInd = k
			maxVal = tmp.Score + _tmp.Score
		}
	}
	return &result[maxInd]
}

// ============================================================================================================================
// FindCycles : returns cycles of matches or a chain of matches
// ============================================================================================================================
func FindCycles(activePairs []Pair) [][]string {
	tmpPairs, nodeList, adjMatrix := generateNodeListAndAdjMatrix(activePairs)
	var ecs *ElementaryCyclesSearch
	ecs = NewElementaryCyclesSearch(adjMatrix, nodeList)
	cycles := ecs.GetElementaryCycles()

	var result [][]string // result contains the IDs of the pairs of each cycle
	result = make([][]string, len(cycles))
	for i := 0; i < len(cycles); i++ {
		cycle := cycles[i]
		result[i] = make([]string, len(cycle))
		for j := 0; j < len(cycle); j++ {
			result[i][j] = tmpPairs[cycle[j]].ID // get the id by its index
		}
	}
	return result
}

// ============================================================================================================================
// generateNodeListAndAdjMatrix : takes the list of pairs as arg and returns nodes list + adjacency matrix
// ============================================================================================================================
func generateNodeListAndAdjMatrix(kidnerPairs []Pair) ([]Pair, []int, [][]bool) {

	var nodes []int // list of nodes : contains the indices of the pairs in kidnerPairs
	nodes = make([]int, len(kidnerPairs))
	for k := range kidnerPairs {
		nodes[k] = k
	}

	var adjMatrix [][]bool // adjacency matrix between nodes
	adjMatrix = make([][]bool, len(kidnerPairs))
	for k := range adjMatrix {
		adjMatrix[k] = make([]bool, len(kidnerPairs))
	}

	for k := range kidnerPairs {
		adjMatrix = EnrollPair(kidnerPairs, adjMatrix, k)
	}

	return kidnerPairs, nodes, adjMatrix
}

// ============================================================================================================================
// EnrollPair : enrolls a pair in the list of pairs and returns the updated adjMatrix
// ============================================================================================================================
func EnrollPair(kidnerPairs []Pair, adjMatrix [][]bool, currentPairIndex int) [][]bool {

	currentPair := kidnerPairs[currentPairIndex]

	for k, v := range kidnerPairs {

		proposal := v.Recipient.Match(currentPair.Donor)  // match currentP.donor and Pi.recip
		_proposal := currentPair.Recipient.Match(v.Donor) // match Pi.donor and currentP.recip

		if proposal.Match { // donor of currentPair -> recip of Pi
			adjMatrix[currentPairIndex][k] = true

		} else if _proposal.Match {
			adjMatrix[k][currentPairIndex] = true
		}
	}

	return adjMatrix
}
