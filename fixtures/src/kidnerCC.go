package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// KidnerCC : simple Chaincode implementation
type KidnerCC struct {
}

var logger = shim.NewLogger("kidner")

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(KidnerCC))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// ============================================================================================================================
// Init resets everything
// ============================================================================================================================
func (t *KidnerCC) Init(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Debug("Init() : calling method -")
	return shim.Success(nil)
}

// ============================================================================================================================
// Invoke is our entry point to invoke a chaincode function
// ============================================================================================================================
func (t *KidnerCC) Invoke(stub shim.ChaincodeStubInterface) pb.Response {

	function, args := stub.GetFunctionAndParameters()
	logger.Debug("Invoke(" + function + ") : calling method -")

	// Handle different functions
	switch function {
	case "createPair":
		return createPair(stub, args)
	case "deletePair":
		return deletePair(stub, args)
	case "getPair":
		return getPair(stub, args)
	case "updatePair":
		return updatePair(stub, args)
	case "deactivatePair":
		return deactivatePair(stub, args)
	case "listActivePairs":
		return listActivePairs(stub, args)
	case "findPairedMatch":
		return findPairedMatch(stub, args)
	case "findMatchCycle":
		return findMatchCycle(stub, args)
	case "getMatch":
		return getMatch(stub, args)
	case "getListMatches":
		return getListMatches(stub, args)
	case "createDoctor":
		return createDoctor(stub, args)
	case "getDoctor":
		return getDoctor(stub, args)
	case "approveMatch":
		return approveMatch(stub, args)
	case "refuseMatch":
		return refuseMatch(stub, args)
	}

	return shim.Error("Received unknown function invocation: " + function)
}

// ============================================================================================================================
//	createPair- create a donor-recipient pair of health records
// ============================================================================================================================
func createPair(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	logger.Debug("running the function createPair()")

	if len(args) != 18 {
		return shim.Error(errNumArgs + "7 for each health record and the Dr ID and the 3 signatures")
	}

	pairID := stub.GetTxID()
	recipient, err := validateHealthRecordInput(args[:7]) // recipient first : indices 0 -> 6
	if err != nil {
		return shim.Error(err.Error())
	}

	donor, err := validateHealthRecordInput(args[7:14]) // donor second : indices 7 -> 13
	if err != nil {
		return shim.Error(err.Error())
	}

	// get the doctor ID
	drId := args[14]
	recipientSig := args[15]
	donorSig := args[16]
	DrSig := args[17]

	if strings.EqualFold(recipient.Type, donor.Type) {
		return shim.Error("2 health records of the same type")
	}

	if strings.EqualFold(recipient.Type, "donor") && strings.EqualFold(donor.Type, "recipient") {
		donor, recipient = recipient, donor
	}

	recipient.Signature = recipientSig
	donor.Signature = donorSig

	var pair Pair
	pair = *recipient.Match(*donor)
	pair.ID = pairID
	pair.DrID = drId
	pair.DrSig = DrSig

	//drSign := args[15]
	// get dr by ID
	// dr, err := getDoctorByID(stub, drId)
	// if err != nil {
	// 	logger.Error(err.Error())
	// 	return shim.Error(errGetState + " for doctor ID ; error: " + err.Error())
	// } else if dr == nil {
	// 	logger.Error("Doctor does not exist; ID : " + drId)
	// 	return shim.Error("Doctor does not exist; ID : " + drId)
	// }

	// // compare the signatures
	// if dr.Signature != drSign {
	// 	logger.Error("Bad Dr signature")
	// 	return shim.Error("Bad Dr signature")
	// }

	jsonP, err := json.Marshal(pair)
	if err != nil {
		logger.Error(err.Error())
		return shim.Error(errJsonMarshall + " : " + err.Error())
	}

	err = stub.PutState(pairID, jsonP)
	if err != nil {
		return shim.Error("createPair() : " + errPutState)
	}

	err = createPairIndex(stub, pair)
	if err != nil {
		return shim.Error("Error creating index for pair")
	}

	return shim.Success([]byte(pairID))
}

// ============================================================================================================================
//	deletePair - delete a Pair
// ============================================================================================================================
func deletePair(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	logger.Debug("calling method - deletePair()")

	if len(args) != 1 {
		return shim.Error(errNumArgs + "the pairID")
	}
	PairID := args[0]

	valAsBytes, err := stub.GetState(PairID)
	if err != nil {
		logger.Error(err.Error())
		return shim.Error(errGetState + "for pairID : " + PairID)
	} else if valAsBytes == nil {
		errStr := "Pair does not exist :" + PairID
		logger.Error(errStr)
		return shim.Error(errStr)
	}

	// unmarshall pair so that we can delete its index
	var pair Pair
	err = json.Unmarshal(valAsBytes, &pair)
	if err != nil {
		errStr := (errJsonUnmarshall + " : " + err.Error())
		logger.Error(errStr)
		return shim.Error(errStr)
	}

	// delete pair
	err = stub.DelState(PairID)
	if err != nil {
		errStr := ("Failed to delete pair:" + err.Error())
		logger.Error(errStr)
		return shim.Error(errStr)
	}

	// delete index
	err = deletePairIndex(stub, pair)
	if err != nil {
		logger.Error(err.Error())
		return shim.Error("failed to delete index for pair")
	}

	return shim.Success([]byte(PairID))
}

// ============================================================================================================================
// getPair : query to get a Pair by its key
// ============================================================================================================================
func getPair(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var pairID string
	var err error

	if len(args) != 1 {
		return shim.Error(errNumArgs + " pairID")
	}
	pairID = args[0]
	valAsbytes, err := stub.GetState(pairID)
	if err != nil {
		logger.Error(err.Error())
		jsonErr := buildError(errGetState + " for pairID : " + pairID)
		return shim.Error(jsonErr)
	} else if valAsbytes == nil {
		jsonErr := buildError("Pair does not exist : " + pairID)
		logger.Error(jsonErr)
		return shim.Error(jsonErr)
	}

	return shim.Success(valAsbytes)
}

// ============================================================================================================================
//	updatePair- update a given donor-recipient pair of health records
// ============================================================================================================================
func updatePair(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// TODO : add control on DrID
	var err error
	logger.Debug("running the function updatePair()")

	if len(args) != 15 {
		return shim.Error(errNumArgs + "15")
	}

	recipient, err := validateHealthRecordInput(args[1:8]) // recipient first : indices 0 -> 6
	if err != nil {
		return shim.Error(err.Error())
	}

	donor, err := validateHealthRecordInput(args[8:]) // donor second : indices 7 -> 13
	if err != nil {
		return shim.Error(err.Error())
	}

	if strings.EqualFold(recipient.Type, donor.Type) {
		return shim.Error("2 health records of the same type")
	}

	// if strings.EqualFold(recipient.Type, "donor") && strings.EqualFold(donor.Type, "recipient") {
	// 	donor, recipient = recipient, donor
	// }

	if strings.EqualFold(recipient.Type, "donor") && strings.EqualFold(donor.Type, "recipient") {
		return shim.Error("the 2 health records are inverted")
	}

	pairID := args[0]
	pair, err := getPairByID(stub, pairID)
	if err != nil {
		return shim.Error(err.Error())
	}

	// keep fields that don't change
	id := pair.ID
	drId := pair.DrID
	drSig := pair.DrSig
	// update pair
	pair = recipient.Match(*donor)
	pair.ID = id
	pair.DrID = drId
	pair.DrSig = drSig

	jsonP, err := json.Marshal(pair)
	if err != nil {
		return shim.Error(errJsonMarshall + " : " + err.Error())
	}

	err = stub.PutState(pairID, jsonP)
	if err != nil {
		return shim.Error("updatePair() : " + errPutState)
	}

	return shim.Success([]byte(pairID))
}

// ============================================================================================================================
//	deactivatePair- set active field to false
// ============================================================================================================================
func deactivatePair(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	logger.Debug("running the function updatePair()")

	if len(args) != 1 {
		return shim.Error(errNumArgs + "pairID")
	}

	pairID := args[0]

	pair, err := getPairByID(stub, pairID)
	if err != nil {
		return shim.Error(err.Error())
	}

	// update fields
	pair.Active = false

	jsonP, err := json.Marshal(pair)
	if err != nil {
		return shim.Error(errJsonMarshall + " : " + err.Error())
	}

	err = stub.PutState(pairID, jsonP)
	if err != nil {
		return shim.Error("deactivatePair() : " + errPutState)
	}

	return shim.Success([]byte(pairID))
}

// =====================================================================================================================
// listActivePairs : Get the list of active Pairs in the chaincode
// =====================================================================================================================
func listActivePairs(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 0 {
		errStr := errNumArgs + "0"
		return shim.Error(buildError(errStr))
	}
	logger.Debug("listActivePairs() : calling method -")
	activePairs, err := getListActivePairs(stub)

	jsonVal, err := json.Marshal(activePairs)
	if err != nil {
		return shim.Error(errJsonMarshall + " : " + err.Error())
	}

	return shim.Success(jsonVal)
}

// ============================================================================================================================
// FindPairedMatch : find a paired match for the spcified pairID
// ============================================================================================================================
func findPairedMatch(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var pairID string
	var err error

	if len(args) != 1 {
		return shim.Error(errNumArgs + " pairID")
	}
	pairID = args[0]
	pair, err := getPairByID(stub, pairID)
	if err != nil {
		logger.Error(err.Error())
		return shim.Error("Failed to get pair")
	}

	activePairs, err := getListActivePairs(stub)
	if err != nil {
		logger.Error(err.Error())
		return shim.Error("Failed to get the list of active pairs")
	} else if activePairs == nil {
		logger.Error("list of active pairs is nil")
		return shim.Error("list of active pairs is nil")
	}

	matchID := stub.GetTxID()
	res := pair.FindBestPairedMatch(activePairs)

	var match Match
	match.ID = matchID
	match.DocType = docTypeMatch
	match.Approved = false
	match.CreateDate = time.Now()

	if res == nil {
		match.MatchedPairs = [][]string{{pair.ID}} // no cross match => put only the pair ID to know that there is no cross-match for that pair
	} else {
		match.MatchedPairs = [][]string{{pair.ID, res.ID}}
	}

	valAsBytes, err := json.Marshal(match)
	if err != nil {
		logger.Error(err.Error())
		return shim.Error(errJsonMarshall)
	}

	err = stub.PutState(matchID, valAsBytes)
	if err != nil {
		logger.Error(err.Error())
		return shim.Error(errPutState + "for match " + matchID)
	}

	err = createMatchIndex(stub, match)
	if err != nil {
		logger.Error(err.Error())
		return shim.Error("fail to create index for match")
	}

	return shim.Success([]byte(matchID))
}

// ============================================================================================================================
// FindMatchCycle : find a chain of matches
// ============================================================================================================================
func findMatchCycle(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	if len(args) != 0 {
		return shim.Error(errNumArgs + " 0 arguments")
	}

	activePairs, err := getListActivePairs(stub)
	if err != nil {
		logger.Error(err.Error())
		return shim.Error("Failed to get the list of active pairs")
	} else if activePairs == nil {
		logger.Error("list of active pairs is nil")
		return shim.Error("list of active pairs is nil")
	}

	matchID := stub.GetTxID()

	res := FindCycles(activePairs)

	var match Match
	match.ID = matchID
	match.DocType = docTypeMatch
	match.Approved = false
	match.MatchedPairs = res
	match.CreateDate = time.Now()

	valAsBytes, err := json.Marshal(match)
	if err != nil {
		logger.Error(err.Error())
		return shim.Error(errJsonMarshall)
	}

	err = stub.PutState(matchID, valAsBytes)
	if err != nil {
		logger.Error(err.Error())
		return shim.Error(errPutState + "for match " + matchID)
	}

	err = createMatchIndex(stub, match)
	if err != nil {
		logger.Error(err.Error())
		return shim.Error("fail to create index for match")
	}

	return shim.Success([]byte(matchID))
}

// ============================================================================================================================
// getMatch : query to get a Match by its key
// ============================================================================================================================
func getMatch(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var matchID string
	var err error

	if len(args) != 1 {
		return shim.Error(errNumArgs + " matchID")
	}
	matchID = args[0]
	valAsbytes, err := stub.GetState(matchID)
	if err != nil {
		logger.Error(err.Error())
		jsonErr := buildError(errGetState + " for matchID : " + matchID)
		return shim.Error(jsonErr)
	} else if valAsbytes == nil {
		jsonErr := buildError("match does not exist : " + matchID)
		logger.Error(jsonErr)
		return shim.Error(jsonErr)
	}

	return shim.Success(valAsbytes)
}

// =====================================================================================================================
// getListMatches : Get the list of matches in the chaincode
// =====================================================================================================================
func getListMatches(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 0 {
		errStr := errNumArgs + "0"
		return shim.Error(buildError(errStr))
	}
	logger.Debug("getListMatches() : calling method -")

	matches, err := getMatchesByIndex(stub, indexMatch, []string{docTypeMatch})
	if err != nil {
		logger.Error(buildError("Error get list of matches"))
		return shim.Error("Failed to get list of matches")
	} else if matches == nil {
		logger.Error("empty list")
		return shim.Error("empty list")
	}

	jsonVal, err := json.Marshal(matches)
	if err != nil {
		return shim.Error(errJsonMarshall + " : " + err.Error())
	}

	return shim.Success(jsonVal)
}

// ============================================================================================================================
//	approveMatch - set active field to false
// ============================================================================================================================
func approveMatch(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	logger.Debug("running the function approveMatch()")

	if len(args) != 3 {
		return shim.Error(errNumArgs + "3")
	}

	DrID := args[0]
	DrSig := args[1]
	matchID := args[2]

	// doctor, err := getDoctorByID(stub, DrID)
	// if err != nil {
	// 	return shim.Error(err.Error())
	// }

	// if doctor.Signature != DrSig {
	// 	return shim.Error("incorrect signature")
	// }

	match, err := getMatchByID(stub, matchID)
	if err != nil {
		return shim.Error(err.Error())
	}

	match.Approved = true
	match.EndorcingDr = DrID
	match.DrSig = DrSig

	jsonM, err := json.Marshal(match)
	if err != nil {
		logger.Error(err.Error())
		return shim.Error(errJsonMarshall)
	}

	err = stub.PutState(matchID, jsonM)
	if err != nil {
		return shim.Error("approveMatch() : " + errPutState)
	}

	return shim.Success([]byte(matchID))
}

// ============================================================================================================================
//	refuseMatch - set active field to false
// ============================================================================================================================
func refuseMatch(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	logger.Debug("running the function refuseMatch()")

	if len(args) != 3 {
		return shim.Error(errNumArgs + "3")
	}

	DrID := args[0]
	DrSig := args[1]
	matchID := args[2]

	// doctor, err := getDoctorByID(stub, DrID)
	// if err != nil {
	// 	return shim.Error(err.Error())
	// }

	// if doctor.Signature != DrSig {
	// 	return shim.Error("incorrect signature")
	// }

	match, err := getMatchByID(stub, matchID)
	if err != nil {
		return shim.Error(err.Error())
	}

	match.Approved = false
	match.EndorcingDr = DrID
	match.DrSig = DrSig

	jsonM, err := json.Marshal(match)
	if err != nil {
		logger.Error(err.Error())
		return shim.Error(errJsonMarshall)
	}

	err = stub.PutState(matchID, jsonM)
	if err != nil {
		return shim.Error("refuseMatch() : " + errPutState)
	}

	return shim.Success([]byte(matchID))
}

// =====================================================================================================================
// create Doctor
// =====================================================================================================================
func createDoctor(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	logger.Debug("running the function createDoctor()")

	if len(args) != 1 {
		return shim.Error(errNumArgs + " the doctor's signature")
	}

	doctorID := stub.GetTxID()

	var doctor Doctor
	doctor.DocType = docTypeDoctor
	doctor.ID = doctorID
	doctor.Signature = args[0]

	jsonDr, err := json.Marshal(doctor)
	if err != nil {
		logger.Error(err.Error())
		return shim.Error(errJsonMarshall)
	}

	err = stub.PutState(doctorID, jsonDr)
	if err != nil {
		logger.Error(err.Error())
		return shim.Error("createDoctor() : " + errPutState)
	}

	return shim.Success([]byte(doctorID))
}

// ============================================================================================================================
// getDoctor : query to get a Doctor by its key
// ============================================================================================================================
func getDoctor(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var drID string
	var err error

	if len(args) != 1 {
		return shim.Error(errNumArgs + " drID")
	}
	drID = args[0]
	valAsbytes, err := stub.GetState(drID)
	if err != nil {
		logger.Error(err.Error())
		jsonErr := buildError(errGetState + " for drID : " + drID)
		return shim.Error(jsonErr)
	} else if valAsbytes == nil {
		jsonErr := buildError("Dr does not exist : " + drID)
		logger.Error(jsonErr)
		return shim.Error(jsonErr)
	}

	return shim.Success(valAsbytes)
}

// ============================================================================================================================
// validateHealthRecordInput : verify args to correctly create a health record
// ============================================================================================================================
func validateHealthRecordInput(args []string) (*HealthRecord, error) {
	if len(args) != 7 {
		return nil, errors.New(errNumArgs + "7")
	}

	age := args[0]
	bloodtype := args[1]
	medicalUrgency := args[2]
	HLAs := args[3]
	PRA := args[4]
	eligible := args[5]
	_type := args[6]
	createDate := time.Now()

	_age, err := strconv.Atoi(age)
	if err != nil {
		return nil, errors.New("Expecting int for Age value")
	}

	if _age < 0 {
		return nil, errors.New("Age value out of bounds")
	}

	// test bloodtype
	if !strings.EqualFold(bloodtype, "A") && !strings.EqualFold(bloodtype, "B") && !strings.EqualFold(bloodtype, "AB") && !strings.EqualFold(bloodtype, "O") {
		return nil, errors.New("Blood type should be A or B or AB or O")
	}

	_medicalUrg, err := strconv.Atoi(medicalUrgency)
	if err != nil {
		return nil, errors.New("Expecting int for medicalUrgency value")
	}

	if _medicalUrg > 3 || _medicalUrg < 0 {
		return nil, errors.New("Medical urgency value out of bounds")
	}

	_PRA, err := strconv.Atoi(PRA)
	if err != nil {
		return nil, errors.New("Expecting int for PRA value")
	}

	if _PRA > 100 || _PRA < 0 {
		return nil, errors.New("PRA value out of bounds")
	}

	_eligible, err := strconv.ParseBool(eligible)
	if err != nil {
		return nil, errors.New("Expecting true or false for eligible value")
	}

	tmp := strings.SplitN(strings.TrimSpace(HLAs), ",", 8)
	if len(tmp) != 8 {
		return nil, errors.New("Wrong number of HLAs")
	}
	_HLAs := make(map[string]string)
	_HLAs["A1"] = strings.TrimSpace(tmp[0])
	_HLAs["A2"] = strings.TrimSpace(tmp[1])
	_HLAs["B1"] = strings.TrimSpace(tmp[2])
	_HLAs["B2"] = strings.TrimSpace(tmp[3])
	_HLAs["DR1"] = strings.TrimSpace(tmp[4])
	_HLAs["DR2"] = strings.TrimSpace(tmp[5])
	_HLAs["DQ1"] = strings.TrimSpace(tmp[6])
	_HLAs["DQ2"] = strings.TrimSpace(tmp[7])

	// test type donor or recipient
	if !strings.EqualFold(_type, "recipient") && !strings.EqualFold(_type, "donor") {
		return nil, errors.New("type should be either donor or recipient")
	}

	value := &HealthRecord{
		docTypeHealthRecord,
		_age,
		bloodtype,
		_medicalUrg,
		_HLAs,
		_PRA,
		_eligible,
		_type,
		"",
		createDate,
	}
	return value, nil
}

// ============================================================================================================================
// Build a json error to return
// ============================================================================================================================
func buildError(errorStr string) (jsonResp string) {
	jsonResp = "{\"Error\":\"" + errorStr + "\"}"
	logger.Error("return: " + jsonResp)
	return
}
