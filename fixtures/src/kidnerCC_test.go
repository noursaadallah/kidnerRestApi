package main

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/op/go-logging"
)

const (
	AGE1            = "50"
	BLOODTYPE1      = "A"
	MEDICALURGENCY1 = "0"
	HLAs1           = "A1,A2,B1,B2,DR1,DR2,DQ1,DQ2"
	PRA1            = "80"
	ELIGIBLE1       = "true"
	TYPE1           = "donor"
	SIGD            = "donorSIG"

	AGE2            = "80"
	BLOODTYPE2      = "B"
	MEDICALURGENCY2 = "1"
	HLAs2           = "A21,A22,B21,B22,DR21,DR22,DQ21,DQ22"
	PRA2            = "0"
	ELIGIBLE2       = "true"
	TYPE2           = "recipient"
	SIGR            = "recipientSIG"

	DRSIG = "signature"
)

func InitLogger() {
	logLevel := logging.INFO
	//consLogLevel := logging.DEBUG
	f := os.Stderr
	//loggerModule := ""
	var format = logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} [%{module}:%{shortfile}] %{level:.4s} : %{color:reset} %{message}`,
	)
	backend := logging.NewLogBackend(f, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	backendLeveled := logging.AddModuleLevel(backendFormatter)
	backendLeveled.SetLevel(logging.Level(logLevel), "mock")

	//backendConsent := logging.NewLogBackend(f, "", 0)
	//backendConsFormatter := logging.NewBackendFormatter(backendConsent, format)
	//backendConsLeveled := logging.AddModuleLevel(backendConsFormatter)
	//backendConsLeveled.SetLevel(consLogLevel,LOGGERMODULE)
	//logging.SetBackend(backendLeveled, backendConsLeveled)
	logging.SetBackend(backendLeveled)
}

// =====================================================================================================================
// check init
// =====================================================================================================================
func TestKidnerCC_CheckInit(t *testing.T) {
	scc := new(KidnerCC)
	stub := shim.NewMockStub("kidner", scc)
	res := stub.MockInit("1", [][]byte{})
	if res.Status != shim.OK {
		t.Log("bad status received, expected: 200; received:" + strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: " + string(res.Message))
		t.FailNow()
	}
}

// =====================================================================================================================
// Create a pair (nominal case)
// =====================================================================================================================
func TestKidnerCC_CreatePairNominal(t *testing.T) {
	scc := new(KidnerCC)
	stub := shim.NewMockStub("kidner", scc)

	res := stub.MockInvoke("1", [][]byte{[]byte("createPair"), []byte(AGE1), []byte(BLOODTYPE1), []byte(MEDICALURGENCY1), []byte(HLAs1), []byte(PRA1), []byte(ELIGIBLE1), []byte(TYPE1), []byte(AGE2), []byte(BLOODTYPE2), []byte(MEDICALURGENCY2), []byte(HLAs2), []byte(PRA2), []byte(ELIGIBLE2), []byte(TYPE2), []byte("DrID"), []byte(SIGR), []byte(SIGD), []byte(DRSIG)})
	if res.Status != shim.OK {
		t.Log("bad status received, expected: 200; received:" + strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: " + string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		t.Log("createPair failed to create a pair and get the pair id")
		t.FailNow()
	}
}

// =====================================================================================================================
// Create a pair wrong input
// =====================================================================================================================
func TestKidnerCC_CreatePairIncorrectNbrArgs(t *testing.T) {
	scc := new(KidnerCC)
	stub := shim.NewMockStub("kidner", scc)
	res := stub.MockInvoke("1", [][]byte{[]byte("createPair"), []byte(BLOODTYPE2), []byte(MEDICALURGENCY2), []byte(HLAs2), []byte(PRA2), []byte(ELIGIBLE2), []byte(TYPE2)})
	if res.Status != shim.ERROR {
		t.Log("bad status received, expected: 500; received:" + strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: " + string(res.Message))
		t.FailNow()
	}
	if res.Payload != nil {
		t.Log("createPair succeeded while it was expected to fail : incorrect number of arguments")
		t.FailNow()
	}
}

func TestKidnerCC_CreatePairNegativeAge(t *testing.T) {
	scc := new(KidnerCC)
	stub := shim.NewMockStub("kidner", scc)
	res := stub.MockInvoke("1", [][]byte{[]byte("createPair"), []byte("-5"), []byte(BLOODTYPE1), []byte(MEDICALURGENCY1), []byte(HLAs1), []byte(PRA1), []byte(ELIGIBLE1), []byte(TYPE1), []byte(AGE2), []byte(BLOODTYPE2), []byte(MEDICALURGENCY2), []byte(HLAs2), []byte(PRA2), []byte(ELIGIBLE2), []byte(TYPE2), []byte("DrID"), []byte(SIGR), []byte(SIGD), []byte(DRSIG)})
	if res.Status != shim.ERROR {
		t.Log("bad status received, expected: 500; received:" + strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: " + string(res.Message))
		t.FailNow()
	}
	if res.Payload != nil {
		t.Log("createPair succeeded while it was expected to fail : negative age value")
		t.FailNow()
	}
}

func TestKidnerCC_CreatePairIncorrectAge(t *testing.T) {
	scc := new(KidnerCC)
	stub := shim.NewMockStub("kidner", scc)
	res := stub.MockInvoke("1", [][]byte{[]byte("createPair"), []byte("a"), []byte(BLOODTYPE1), []byte(MEDICALURGENCY1), []byte(HLAs1), []byte(PRA1), []byte(ELIGIBLE1), []byte(TYPE1), []byte(AGE2), []byte(BLOODTYPE2), []byte(MEDICALURGENCY2), []byte(HLAs2), []byte(PRA2), []byte(ELIGIBLE2), []byte(TYPE2), []byte("DrID"), []byte(SIGR), []byte(SIGD), []byte(DRSIG)})
	if res.Status != shim.ERROR {
		t.Log("bad status received, expected: 500; received:" + strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: " + string(res.Message))
		t.FailNow()
	}
	if res.Payload != nil {
		t.Log("createPair succeeded while it was expected to fail : age is not int")
		t.FailNow()
	}
}

func TestKidnerCC_CreatePairIncorrectBloodType(t *testing.T) {
	scc := new(KidnerCC)
	stub := shim.NewMockStub("kidner", scc)
	res := stub.MockInvoke("1", [][]byte{[]byte("createPair"), []byte(AGE1), []byte("wrongbloodtype"), []byte(MEDICALURGENCY1), []byte(HLAs1), []byte(PRA1), []byte(ELIGIBLE1), []byte(TYPE1), []byte(AGE2), []byte(BLOODTYPE2), []byte(MEDICALURGENCY2), []byte(HLAs2), []byte(PRA2), []byte(ELIGIBLE2), []byte(TYPE2), []byte("DrID"), []byte(SIGR), []byte(SIGD), []byte(DRSIG)})
	if res.Status != shim.ERROR {
		t.Log("bad status received, expected: 500; received:" + strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: " + string(res.Message))
		t.FailNow()
	}
	if res.Payload != nil {
		t.Log("createPair succeeded while it was expected to fail : incorrect blood type")
		t.FailNow()
	}
}

func TestKidnerCC_CreatePairIncorrectMedUrg(t *testing.T) {
	scc := new(KidnerCC)
	stub := shim.NewMockStub("kidner", scc)
	res := stub.MockInvoke("1", [][]byte{[]byte("createPair"), []byte(AGE1), []byte(BLOODTYPE1), []byte("MEDICALURGENCY1"), []byte(HLAs1), []byte(PRA1), []byte(ELIGIBLE1), []byte(TYPE1), []byte(AGE2), []byte(BLOODTYPE2), []byte(MEDICALURGENCY2), []byte(HLAs2), []byte(PRA2), []byte(ELIGIBLE2), []byte(TYPE2), []byte("DrID"), []byte(SIGR), []byte(SIGD), []byte(DRSIG)})
	if res.Status != shim.ERROR {
		t.Log("bad status received, expected: 500; received:" + strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: " + string(res.Message))
		t.FailNow()
	}
	if res.Payload != nil {
		t.Log("createPair succeeded while it was expected to fail : medical urgency is not int")
		t.FailNow()
	}
}

func TestKidnerCC_CreatePairMedUrgOutBounds(t *testing.T) {
	scc := new(KidnerCC)
	stub := shim.NewMockStub("kidner", scc)
	res := stub.MockInvoke("1", [][]byte{[]byte("createPair"), []byte(AGE1), []byte(BLOODTYPE1), []byte("50"), []byte(HLAs1), []byte(PRA1), []byte(ELIGIBLE1), []byte(TYPE1), []byte(AGE2), []byte(BLOODTYPE2), []byte(MEDICALURGENCY2), []byte(HLAs2), []byte(PRA2), []byte(ELIGIBLE2), []byte(TYPE2), []byte("DrID"), []byte(SIGR), []byte(SIGD), []byte(DRSIG)})
	if res.Status != shim.ERROR {
		t.Log("bad status received, expected: 500; received:" + strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: " + string(res.Message))
		t.FailNow()
	}
	if res.Payload != nil {
		t.Log("createPair succeeded while it was expected to fail : medical urgency out of bounds")
		t.FailNow()
	}
}

func TestKidnerCC_CreatePairPRANotInt(t *testing.T) {
	scc := new(KidnerCC)
	stub := shim.NewMockStub("kidner", scc)
	res := stub.MockInvoke("1", [][]byte{[]byte("createPair"), []byte(AGE1), []byte(BLOODTYPE1), []byte(MEDICALURGENCY1), []byte(HLAs1), []byte("PRA1"), []byte(ELIGIBLE1), []byte(TYPE1), []byte(AGE2), []byte(BLOODTYPE2), []byte(MEDICALURGENCY2), []byte(HLAs2), []byte(PRA2), []byte(ELIGIBLE2), []byte(TYPE2), []byte("DrID"), []byte(SIGR), []byte(SIGD), []byte(DRSIG)})
	if res.Status != shim.ERROR {
		t.Log("bad status received, expected: 500; received:" + strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: " + string(res.Message))
		t.FailNow()
	}
	if res.Payload != nil {
		t.Log("createPair succeeded while it was expected to fail : PRA is not int")
		t.FailNow()
	}
}

func TestKidnerCC_CreatePairPRAOutBounds(t *testing.T) {
	scc := new(KidnerCC)
	stub := shim.NewMockStub("kidner", scc)
	res := stub.MockInvoke("1", [][]byte{[]byte("createPair"), []byte(AGE1), []byte(BLOODTYPE1), []byte(MEDICALURGENCY1), []byte(HLAs1), []byte("150"), []byte(ELIGIBLE1), []byte(TYPE1), []byte(AGE2), []byte(BLOODTYPE2), []byte(MEDICALURGENCY2), []byte(HLAs2), []byte(PRA2), []byte(ELIGIBLE2), []byte(TYPE2), []byte("DrID"), []byte(SIGR), []byte(SIGD), []byte(DRSIG)})
	if res.Status != shim.ERROR {
		t.Log("bad status received, expected: 500; received:" + strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: " + string(res.Message))
		t.FailNow()
	}
	if res.Payload != nil {
		t.Log("createPair succeeded while it was expected to fail : PRA value out of bound")
		t.FailNow()
	}
}

func TestKidnerCC_CreatePairIncorrectEligibility(t *testing.T) {
	scc := new(KidnerCC)
	stub := shim.NewMockStub("kidner", scc)
	res := stub.MockInvoke("1", [][]byte{[]byte("createPair"), []byte(AGE1), []byte(BLOODTYPE1), []byte(MEDICALURGENCY1), []byte(HLAs1), []byte(PRA1), []byte("ELIGIBLE1"), []byte(TYPE1), []byte(AGE2), []byte(BLOODTYPE2), []byte(MEDICALURGENCY2), []byte(HLAs2), []byte(PRA2), []byte(ELIGIBLE2), []byte(TYPE2), []byte("DrID"), []byte(SIGR), []byte(SIGD), []byte(DRSIG)})
	if res.Status != shim.ERROR {
		t.Log("bad status received, expected: 500; received:" + strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: " + string(res.Message))
		t.FailNow()
	}
	if res.Payload != nil {
		t.Log("createPair succeeded while it was expected to fail : eligible field is not bool")
		t.FailNow()
	}
}

func TestKidnerCC_CreatePairIncorrectHLAs(t *testing.T) {
	scc := new(KidnerCC)
	stub := shim.NewMockStub("kidner", scc)
	res := stub.MockInvoke("1", [][]byte{[]byte("createPair"), []byte(AGE1), []byte(BLOODTYPE1), []byte(MEDICALURGENCY1), []byte("HLAs1"), []byte(PRA1), []byte(ELIGIBLE1), []byte(TYPE1), []byte(AGE2), []byte(BLOODTYPE2), []byte(MEDICALURGENCY2), []byte(HLAs2), []byte(PRA2), []byte(ELIGIBLE2), []byte(TYPE2), []byte("DrID"), []byte(SIGR), []byte(SIGD), []byte(DRSIG)})
	if res.Status != shim.ERROR {
		t.Log("bad status received, expected: 500; received:" + strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: " + string(res.Message))
		t.FailNow()
	}
	if res.Payload != nil {
		t.Log("createPair succeeded while it was expected to fail : HLAs incorrect")
		t.FailNow()
	}
}

func TestKidnerCC_CreatePairIncorrectType(t *testing.T) {
	scc := new(KidnerCC)
	stub := shim.NewMockStub("kidner", scc)
	res := stub.MockInvoke("1", [][]byte{[]byte("createPair"), []byte(AGE1), []byte(BLOODTYPE1), []byte(MEDICALURGENCY1), []byte(HLAs1), []byte(PRA1), []byte(ELIGIBLE1), []byte("TYPE1"), []byte(AGE2), []byte(BLOODTYPE2), []byte(MEDICALURGENCY2), []byte(HLAs2), []byte(PRA2), []byte(ELIGIBLE2), []byte(TYPE2), []byte("DrID"), []byte(SIGR), []byte(SIGD), []byte(DRSIG)})
	if res.Status != shim.ERROR {
		t.Log("bad status received, expected: 500; received:" + strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: " + string(res.Message))
		t.FailNow()
	}
	if res.Payload != nil {
		t.Log("createPair succeeded while it was expected to fail : type is not donor nor recipient")
		t.FailNow()
	}
}

func TestKidnerCC_CreatePairSameType(t *testing.T) {
	scc := new(KidnerCC)
	stub := shim.NewMockStub("kidner", scc)
	res := stub.MockInvoke("1", [][]byte{[]byte("createPair"), []byte(AGE1), []byte(BLOODTYPE1), []byte(MEDICALURGENCY1), []byte(HLAs1), []byte(PRA1), []byte(ELIGIBLE1), []byte(TYPE1), []byte(AGE2), []byte(BLOODTYPE2), []byte(MEDICALURGENCY2), []byte(HLAs2), []byte(PRA2), []byte(ELIGIBLE2), []byte(TYPE1), []byte("DrID"), []byte(SIGR), []byte(SIGD), []byte(DRSIG)})
	if res.Status != shim.ERROR {
		t.Log("bad status received, expected: 500; received:" + strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: " + string(res.Message))
		t.FailNow()
	}
	if res.Payload != nil {
		t.Log("createPair succeeded while it was expected to fail : same type ")
		t.FailNow()
	}
}

// =====================================================================================================================
// Get a pair (nominal case)
// =====================================================================================================================
func TestKidnerCC_GetPairNominal(t *testing.T) {
	scc := new(KidnerCC)
	stub := shim.NewMockStub("kidner", scc)

	res := stub.MockInvoke("1", [][]byte{[]byte("createPair"), []byte(AGE1), []byte(BLOODTYPE1), []byte(MEDICALURGENCY1), []byte(HLAs1), []byte(PRA1), []byte(ELIGIBLE1), []byte(TYPE1), []byte(AGE2), []byte(BLOODTYPE2), []byte(MEDICALURGENCY2), []byte(HLAs2), []byte(PRA2), []byte(ELIGIBLE2), []byte(TYPE2), []byte("DrID"), []byte(SIGR), []byte(SIGD), []byte(DRSIG)})
	pairID := string(res.Payload)
	res = stub.MockInvoke("2", [][]byte{[]byte("getPair"), []byte(pairID)})
	if res.Status != shim.OK {
		t.Log("bad status received, expected: 200 received:" + strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: " + string(res.Message))
		t.FailNow()
	}
	var pair Pair
	err := json.Unmarshal(res.Payload, &pair)
	if err != nil {
		t.Log("getPair", string(res.Payload))
		t.FailNow()
	}
	if pair.ID != pairID {
		t.Log("getPair bad pairID")
		t.FailNow()
	}

	if pair.Recipient.Signature != SIGR {
		t.Log("bad recipient signature")
		t.FailNow()
	}
	if pair.Donor.Signature != SIGD {
		t.Log("bad donor signature")
		t.FailNow()
	}
}

// =====================================================================================================================
// Get a pair Incorrect args
// =====================================================================================================================
func TestKidnerCC_GetPairIncorrectArgs(t *testing.T) {
	scc := new(KidnerCC)
	stub := shim.NewMockStub("kidner", scc)

	res := stub.MockInvoke("2", [][]byte{[]byte("getPair")})
	if res.Status != shim.ERROR {
		t.Log("bad status received, expected: 500 received:" + strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: " + string(res.Message))
		t.FailNow()
	}
}

// =====================================================================================================================
// Delete a pair (nominal case)
// =====================================================================================================================
func TestKidnerCC_DeletePairNominal(t *testing.T) {
	scc := new(KidnerCC)
	stub := shim.NewMockStub("kidner", scc)

	res := stub.MockInvoke("1", [][]byte{[]byte("createPair"), []byte(AGE1), []byte(BLOODTYPE1), []byte(MEDICALURGENCY1), []byte(HLAs1), []byte(PRA1), []byte(ELIGIBLE1), []byte(TYPE1), []byte(AGE2), []byte(BLOODTYPE2), []byte(MEDICALURGENCY2), []byte(HLAs2), []byte(PRA2), []byte(ELIGIBLE2), []byte(TYPE2), []byte("DrID"), []byte(SIGR), []byte(SIGD), []byte(DRSIG)})
	pairID := string(res.Payload)
	res = stub.MockInvoke("2", [][]byte{[]byte("deletePair"), []byte(pairID)})
	if res.Status != shim.OK {
		t.Log("bad status received, expected: 200 received:" + strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: " + string(res.Message))
		t.FailNow()
	}
	res = stub.MockInvoke("3", [][]byte{[]byte("getPair"), []byte(pairID)})
	if res.Status != shim.ERROR {
		t.Log("bad status received, expected: 500 received:" + strconv.FormatInt(int64(res.Status), 10))
		t.FailNow()
	}
	if !strings.Contains(res.Message, "Pair does not exist") {
		t.Log("Bad return message, expected:" + "Pair does not exist" + " reveived:" + string(res.Message))
		t.FailNow()
	}
}

// =====================================================================================================================
// Delete a pair incorrect number of arguments
// =====================================================================================================================
func TestKidnerCC_DeletePairIncorrectArgs(t *testing.T) {
	scc := new(KidnerCC)
	stub := shim.NewMockStub("kidner", scc)

	res := stub.MockInvoke("1", [][]byte{[]byte("createPair"), []byte(AGE1), []byte(BLOODTYPE1), []byte(MEDICALURGENCY1), []byte(HLAs1), []byte(PRA1), []byte(ELIGIBLE1), []byte(TYPE1), []byte(AGE2), []byte(BLOODTYPE2), []byte(MEDICALURGENCY2), []byte(HLAs2), []byte(PRA2), []byte(ELIGIBLE2), []byte(TYPE2), []byte("DrID"), []byte(SIGR), []byte(SIGD), []byte(DRSIG)})
	pairID := string(res.Payload)
	res = stub.MockInvoke("2", [][]byte{[]byte("deletePair"), []byte(pairID), []byte("additionalArg")})
	if res.Status != shim.ERROR {
		t.Log("bad status received, expected: 500 received:" + strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: " + string(res.Message))
		t.FailNow()
	}
	if !strings.Contains(res.Message, errNumArgs) {
		t.Log("incorrect error message")
		t.Log(res.Message)
		t.FailNow()
	}
}

// =====================================================================================================================
// Get bad function name --> error
// =====================================================================================================================
func TestKidnerCC_GetBadFunction(t *testing.T) {
	scc := new(KidnerCC)
	stub := shim.NewMockStub("kidner", scc)
	res := stub.MockInvoke("1", [][]byte{[]byte("badFunction")})
	if res.Status != shim.ERROR {
		t.Log("bad status received, expected: 500 received:" + strconv.FormatInt(int64(res.Status), 10))
		t.FailNow()
	}
	if !strings.Contains("Received unknown function invocation: ", "Received unknown function invocation: ") {
		t.Log("Bad return message, expected:" + "Received unknown function invocation: " + " ; received:" + string(res.Message))
		t.FailNow()
	}
}

// =====================================================================================================================
// Create a pair with a missing param  --> error
// =====================================================================================================================
func TestKidnerCC_CreatePairWithMissingParams(t *testing.T) {
	scc := new(KidnerCC)
	stub := shim.NewMockStub("kidner", scc)
	res := stub.MockInvoke("1", [][]byte{[]byte("createPair"), []byte(AGE1), []byte(MEDICALURGENCY1), []byte(HLAs1), []byte(PRA1), []byte(ELIGIBLE1), []byte(TYPE1), []byte(AGE2), []byte(BLOODTYPE2), []byte(MEDICALURGENCY2), []byte(HLAs2), []byte(PRA2), []byte(ELIGIBLE2), []byte(TYPE2)})

	if res.Status != shim.ERROR {
		t.Log("bad status received, expected: 500 received:" + strconv.FormatInt(int64(res.Status), 10))
		t.FailNow()
	}
	if !strings.Contains(res.Message, "Incorrect number of arguments") {
		t.Log("Bad return message, expected:" + "Incorrect number of arguments" + " received:" + string(res.Message))
		t.FailNow()
	}
}

// =====================================================================================================================
// Delete a pair with a bad pair ID --> error
// =====================================================================================================================
func TestKidnerCC_DeletePairWithBadPairID(t *testing.T) {
	scc := new(KidnerCC)
	stub := shim.NewMockStub("kidner", scc)
	res := stub.MockInvoke("1", [][]byte{[]byte("createPair"), []byte(AGE1), []byte(BLOODTYPE1), []byte(MEDICALURGENCY1), []byte(HLAs1), []byte(PRA1), []byte(ELIGIBLE1), []byte(TYPE1), []byte(AGE2), []byte(BLOODTYPE2), []byte(MEDICALURGENCY2), []byte(HLAs2), []byte(PRA2), []byte(ELIGIBLE2), []byte(TYPE2), []byte("DrID"), []byte(SIGR), []byte(SIGD), []byte(DRSIG)})
	res = stub.MockInvoke("2", [][]byte{[]byte("deletePair"), []byte("badpairid")})
	if res.Status != shim.ERROR {
		t.Log("bad status received, expected: 500 received:" + strconv.FormatInt(int64(res.Status), 10))
		t.FailNow()
	}
	if !strings.Contains(res.Message, "Pair does not exist") {
		t.Log("Bad return message, expected:" + "Pair does not exist" + " received:" + string(res.Message))
		t.FailNow()
	}
}

// =====================================================================================================================
// Get deleted pair --> error
// =====================================================================================================================
func TestKidner_GetDeletedPair(t *testing.T) {
	scc := new(KidnerCC)
	stub := shim.NewMockStub("kidner", scc)

	res := stub.MockInvoke("1", [][]byte{[]byte("createPair"), []byte(AGE1), []byte(BLOODTYPE1), []byte(MEDICALURGENCY1), []byte(HLAs1), []byte(PRA1), []byte(ELIGIBLE1), []byte(TYPE1), []byte(AGE2), []byte(BLOODTYPE2), []byte(MEDICALURGENCY2), []byte(HLAs2), []byte(PRA2), []byte(ELIGIBLE2), []byte(TYPE2), []byte("DrID"), []byte(SIGR), []byte(SIGD), []byte(DRSIG)})
	pairID := string(res.Payload)
	res = stub.MockInvoke("2", [][]byte{[]byte("deletePair"), []byte(pairID)})
	if res.Status != shim.OK {
		t.Log("deletePair", string(res.Message))
		t.FailNow()
	}
	res = stub.MockInvoke("3", [][]byte{[]byte("getPair"), []byte(pairID)})
	if res.Status == shim.OK {
		t.Log("getPair: this pair should be deleted")
		t.FailNow()
	}
	if res.Status != shim.ERROR {
		t.Log("bad status received, expected: 500 received:" + strconv.FormatInt(int64(res.Status), 10))
		t.FailNow()
	}
	if !strings.Contains(res.Message, "Pair does not exist") {
		t.Log("Bad return message, expected:" + "Pair does not exist" + " reveived:" + string(res.Message))
		t.FailNow()
	}
}

// =====================================================================================================================
// Update a pair (nominal case)
// =====================================================================================================================
func TestKidnerCC_UpdatePairNominal(t *testing.T) {
	scc := new(KidnerCC)
	stub := shim.NewMockStub("kidner", scc)

	res1 := stub.MockInvoke("1", [][]byte{[]byte("createPair"), []byte(AGE1), []byte(BLOODTYPE1), []byte(MEDICALURGENCY1), []byte(HLAs1), []byte(PRA1), []byte(ELIGIBLE1), []byte(TYPE1),
		[]byte(AGE2), []byte(BLOODTYPE2), []byte(MEDICALURGENCY2), []byte(HLAs2), []byte(PRA2), []byte(ELIGIBLE2), []byte(TYPE2), []byte("DrID"), []byte(SIGR), []byte(SIGD), []byte(DRSIG)})
	if res1.Status != shim.OK {
		t.Log("bad status received, expected: 200; received:" + strconv.FormatInt(int64(res1.Status), 10))
		t.Log("response: " + string(res1.Message))
		t.FailNow()
	}
	if res1.Payload == nil {
		t.Log("createPair failed to create a pair and get the pair id")
		t.FailNow()
	}

	res2 := stub.MockInvoke("2", [][]byte{[]byte("updatePair"), []byte(res1.Payload),
		[]byte(AGE2), []byte(BLOODTYPE2), []byte(MEDICALURGENCY2), []byte(HLAs2), []byte(PRA2), []byte(ELIGIBLE2), []byte(TYPE2),
		[]byte(AGE1), []byte(BLOODTYPE1), []byte(MEDICALURGENCY1), []byte(HLAs1), []byte(PRA1), []byte(ELIGIBLE1), []byte(TYPE1)})
	if res2.Status != shim.OK {
		t.Log("bad status received, expected: 200; received:" + strconv.FormatInt(int64(res2.Status), 10))
		t.Log("response: " + string(res2.Message))
		t.FailNow()
	}
	if res2.Payload == nil {
		t.Log("updatePair failed to update the pair and get the pair id")
		t.FailNow()
	}

	if string(res1.Payload) != string(res2.Payload) {
		t.Log("wrong id updated")
		t.FailNow()
	}
}

// =====================================================================================================================
// Update a pair incorrect args
// =====================================================================================================================
func TestKidnerCC_UpdatePairIncorrectArgs(t *testing.T) {
	scc := new(KidnerCC)
	stub := shim.NewMockStub("kidner", scc)

	res2 := stub.MockInvoke("2", [][]byte{[]byte("updatePair"), []byte("ID")})
	if res2.Status != shim.ERROR {
		t.Log("bad status received, expected: 500; received:" + strconv.FormatInt(int64(res2.Status), 10))
		t.Log("response: " + string(res2.Message))
		t.FailNow()
	}
}

// =====================================================================================================================
// Update a pair two health records same type
// =====================================================================================================================
func TestKidnerCC_UpdatePairSameType(t *testing.T) {
	scc := new(KidnerCC)
	stub := shim.NewMockStub("kidner", scc)

	res2 := stub.MockInvoke("2", [][]byte{[]byte("updatePair"), []byte("ID"),
		[]byte(AGE2), []byte(BLOODTYPE2), []byte(MEDICALURGENCY2), []byte(HLAs2), []byte(PRA2), []byte(ELIGIBLE2), []byte(TYPE1),
		[]byte(AGE1), []byte(BLOODTYPE1), []byte(MEDICALURGENCY1), []byte(HLAs1), []byte(PRA1), []byte(ELIGIBLE1), []byte(TYPE1)})
	if res2.Status != shim.ERROR {
		t.Log("bad status received, expected: 500; received:" + strconv.FormatInt(int64(res2.Status), 10))
		t.Log("response: " + string(res2.Message))
		t.FailNow()
	}
	if !strings.Contains(res2.Message, "2 health records of the same type") {
		t.Log("bad error message received")
		t.Log("response: " + string(res2.Message))
		t.FailNow()
	}
}

// =====================================================================================================================
// Deactivate a pair => nominal case
// =====================================================================================================================
func TestKidnerCC_DeactivatePairNominal(t *testing.T) {
	scc := new(KidnerCC)
	stub := shim.NewMockStub("kidner", scc)

	res := stub.MockInvoke("1", [][]byte{[]byte("createPair"), []byte(AGE1), []byte(BLOODTYPE1), []byte(MEDICALURGENCY1), []byte(HLAs1), []byte(PRA1), []byte(ELIGIBLE1), []byte(TYPE1), []byte(AGE2), []byte(BLOODTYPE2), []byte(MEDICALURGENCY2), []byte(HLAs2), []byte(PRA2), []byte(ELIGIBLE2), []byte(TYPE2), []byte("DrID"), []byte(SIGR), []byte(SIGD), []byte(DRSIG)})
	pairID := string(res.Payload)
	res = stub.MockInvoke("2", [][]byte{[]byte("deactivatePair"), []byte(pairID)})
	if res.Status != shim.OK {
		t.Log("bad status received, expected: 200 received:" + strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: " + string(res.Message))
		t.FailNow()
	}

	res = stub.MockInvoke("3", [][]byte{[]byte("getPair"), []byte(pairID)})
	if res.Status != shim.OK {
		t.Log("bad status received, expected: 200 received:" + strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: " + string(res.Message))
		t.FailNow()
	}

	var pair Pair
	err := json.Unmarshal(res.Payload, &pair)
	if err != nil {
		t.Log("Unmarshall getPair", string(res.Payload))
		t.FailNow()
	}
	if pair.Active != false {
		t.Error("expected inactive pair")
		t.FailNow()
	}
}

// =====================================================================================================================
// Deactivate a pair => incorrect args
// =====================================================================================================================
func TestKidnerCC_DeactivatePairIncorrectArgs(t *testing.T) {
	scc := new(KidnerCC)
	stub := shim.NewMockStub("kidner", scc)

	res := stub.MockInvoke("1", [][]byte{[]byte("createPair"), []byte(AGE1), []byte(BLOODTYPE1), []byte(MEDICALURGENCY1), []byte(HLAs1), []byte(PRA1), []byte(ELIGIBLE1), []byte(TYPE1),
		[]byte(AGE2), []byte(BLOODTYPE2), []byte(MEDICALURGENCY2), []byte(HLAs2), []byte(PRA2), []byte(ELIGIBLE2), []byte(TYPE2), []byte("DrID"), []byte(SIGR), []byte(SIGD), []byte(DRSIG)})
	pairID := string(res.Payload)
	res = stub.MockInvoke("2", [][]byte{[]byte("deactivatePair")})
	if res.Status != shim.ERROR {
		t.Log("bad status received, expected: 500 received:" + strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: " + string(res.Message))
		t.FailNow()
	}

	res = stub.MockInvoke("3", [][]byte{[]byte("getPair"), []byte(pairID)})
	if res.Status != shim.OK {
		t.Log("bad status received, expected: 200 received:" + strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: " + string(res.Message))
		t.FailNow()
	}

	var pair Pair
	err := json.Unmarshal(res.Payload, &pair)
	if err != nil {
		t.Log("Unmarshall getPair", string(res.Payload))
		t.FailNow()
	}
	if pair.Active != true {
		t.Error("expected inactive pair")
		t.FailNow()
	}
}

// =====================================================================================================================
// Get list of active pairs empty list
// =====================================================================================================================
func TestKidnerCC_ListActivePairsEmpty(t *testing.T) {
	scc := new(KidnerCC)
	stub := shim.NewMockStub("kidner", scc)
	res := stub.MockInvoke("1", [][]byte{[]byte("listActivePairs")})
	if res.Status != shim.OK {
		t.Log("bad status received, expected: 200 received:" + strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: " + string(res.Message))
		t.FailNow()
	}
	pairs := make([]Pair, 0)
	err := json.Unmarshal(res.Payload, &pairs)
	if err != nil {
		t.Log("listActivePairs", string(res.Payload))
		t.FailNow()
	}
	if len(pairs) != 0 {
		t.Error("empty list expected, but length =", strconv.Itoa(len(pairs)), "received")
		t.FailNow()
	}
}

// =====================================================================================================================
// Get list of active pairs incorrect args
// =====================================================================================================================
func TestKidnerCC_ListActivePairsIncorrectArgs(t *testing.T) {
	scc := new(KidnerCC)
	stub := shim.NewMockStub("kidner", scc)
	res := stub.MockInvoke("1", [][]byte{[]byte("listActivePairs"), []byte("incorrectArg")})
	if res.Status != shim.ERROR {
		t.Log("bad status received, expected: 500 received:" + strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: " + string(res.Message))
		t.FailNow()
	}
}

// =====================================================================================================================
// Get list of active pairs nominal : add 3 pairs, deactivate 1 and delete 1 => return 1
// =====================================================================================================================
func TestKidnerCC_ListActivePairsNominal(t *testing.T) {
	scc := new(KidnerCC)
	stub := shim.NewMockStub("kidner", scc)

	res := stub.MockInvoke("1", [][]byte{[]byte("createPair"), []byte(AGE1), []byte(BLOODTYPE1), []byte(MEDICALURGENCY1), []byte(HLAs1), []byte(PRA1), []byte(ELIGIBLE1), []byte(TYPE1), []byte(AGE2), []byte(BLOODTYPE2), []byte(MEDICALURGENCY2), []byte(HLAs2), []byte(PRA2), []byte(ELIGIBLE2), []byte(TYPE2), []byte("DrID"), []byte(SIGR), []byte(SIGD), []byte(DRSIG)})
	pairID := res.Payload
	res = stub.MockInvoke("2", [][]byte{[]byte("deletePair"), []byte(pairID)})
	if res.Status != shim.OK {
		t.Log("bad status received, expected: 200 received:" + strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: " + string(res.Message))
		t.FailNow()
	}

	res = stub.MockInvoke("3", [][]byte{[]byte("createPair"), []byte(AGE1), []byte(BLOODTYPE1), []byte(MEDICALURGENCY1), []byte(HLAs1), []byte(PRA1), []byte(ELIGIBLE1), []byte(TYPE1), []byte(AGE2), []byte(BLOODTYPE2), []byte(MEDICALURGENCY2), []byte(HLAs2), []byte(PRA2), []byte(ELIGIBLE2), []byte(TYPE2), []byte("DrID"), []byte(SIGR), []byte(SIGD), []byte(DRSIG)})
	pairID = res.Payload
	res = stub.MockInvoke("4", [][]byte{[]byte("deactivatePair"), []byte(pairID)})
	if res.Status != shim.OK {
		t.Log("bad status received, expected: 200 received:" + strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: " + string(res.Message))
		t.FailNow()
	}

	res = stub.MockInvoke("5", [][]byte{[]byte("createPair"), []byte(AGE1), []byte(BLOODTYPE1), []byte(MEDICALURGENCY1), []byte(HLAs1), []byte(PRA1), []byte(ELIGIBLE1), []byte(TYPE1), []byte(AGE2), []byte(BLOODTYPE2), []byte(MEDICALURGENCY2), []byte(HLAs2), []byte(PRA2), []byte(ELIGIBLE2), []byte(TYPE2), []byte("DrID"), []byte(SIGR), []byte(SIGD), []byte(DRSIG)})

	res = stub.MockInvoke("6", [][]byte{[]byte("listActivePairs")})
	if res.Status != shim.OK {
		t.Log("bad status received, expected: 200 received:" + strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: " + string(res.Message))
		t.FailNow()
	}
	pairs := make([]Pair, 0)
	err := json.Unmarshal(res.Payload, &pairs)
	if err != nil {
		t.Log("listActivePairs", string(res.Payload))
		t.FailNow()
	}
	if len(pairs) != 1 {
		t.Error("expected length = 1, received length =", strconv.Itoa(len(pairs)))
		t.FailNow()
	}
}

// =====================================================================================================================
// FindPairedMatch nominal : 2 pairs with a perfect paired match
// then check getMatch then check getListMatches
// =====================================================================================================================
func TestKidnerCC_FindPairedMatchNominal(t *testing.T) {
	scc := new(KidnerCC)
	stub := shim.NewMockStub("kidner", scc)

	// create 1st Pair D1,R1
	res1 := stub.MockInvoke("1", [][]byte{[]byte("createPair"), []byte(AGE1), []byte(BLOODTYPE1), []byte(MEDICALURGENCY1), []byte(HLAs1), []byte(PRA1), []byte(ELIGIBLE1), []byte(TYPE1),
		[]byte(AGE2), []byte(BLOODTYPE2), []byte(MEDICALURGENCY2), []byte(HLAs2), []byte(PRA2), []byte(ELIGIBLE2), []byte(TYPE2), []byte("DrID"), []byte(SIGR), []byte(SIGD), []byte(DRSIG)})
	pair1ID := res1.Payload
	if res1.Status != shim.OK {
		t.Log("bad status received, expected: 200 received:" + strconv.FormatInt(int64(res1.Status), 10))
		t.Log("response: " + string(res1.Message))
		t.FailNow()
	}

	// create 2nd Pair D2,R2
	res2 := stub.MockInvoke("2", [][]byte{[]byte("createPair"), []byte(AGE2), []byte(BLOODTYPE2), []byte(MEDICALURGENCY2), []byte(HLAs2), []byte(PRA2), []byte(ELIGIBLE2), []byte(TYPE1),
		[]byte(AGE1), []byte(BLOODTYPE1), []byte(MEDICALURGENCY1), []byte(HLAs1), []byte(PRA1), []byte(ELIGIBLE1), []byte(TYPE2), []byte("DrID"), []byte(SIGR), []byte(SIGD), []byte(DRSIG)})
	pair2ID := res2.Payload
	if res2.Status != shim.OK {
		t.Log("bad status received, expected: 200 received:" + strconv.FormatInt(int64(res2.Status), 10))
		t.Log("response: " + string(res2.Message))
		t.FailNow()
	}

	// Pair1 search for paired match => a match with pair2 is expected
	res := stub.MockInvoke("3", [][]byte{[]byte("findPairedMatch"), []byte(pair1ID)})
	matchID := res.Payload
	if res.Status != shim.OK {
		t.Log("bad status received, expected: 200 received:" + strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: " + string(res.Message))
		t.FailNow()
	}

	// see if the match exists
	res = stub.MockInvoke("4", [][]byte{[]byte("getMatch"), []byte(matchID)})
	if res.Status != shim.OK {
		t.Log("bad status received, expected: 200 received:" + strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: " + string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		t.Log("match not found")
		t.Log("response : " + string(res.Message))
		t.FailNow()
	}

	// retrieve the match and check its attributes
	var match Match
	err := json.Unmarshal(res.Payload, &match)
	if err != nil {
		t.Log(errJsonUnmarshall)
		t.Log(err.Error())
		t.FailNow()
	}
	if match.Approved {
		t.Log("match should not be approved")
		t.FailNow()
	}

	if match.DocType != docTypeMatch {
		t.Log("incorrect docType")
		t.FailNow()
	}

	if match.MatchedPairs == nil || len(match.MatchedPairs) == 0 {
		t.Log("match does not contain any matched pairs")
		t.FailNow()
	}

	if match.MatchedPairs[0][0] != string(pair1ID) || match.MatchedPairs[0][1] != string(pair2ID) {
		t.Log("matched pairs IDs are incorrect")
		t.FailNow()
	}

	res = stub.MockInvoke("5", [][]byte{[]byte("getListMatches")})
	if res.Status != shim.OK {
		t.Log("bad status received, expected: 200 received:" + strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: " + string(res.Message))
		t.FailNow()
	}
	matches := make([]Match, 0)
	err = json.Unmarshal(res.Payload, &matches)
	if err != nil {
		t.Error("listMatches", string(res.Payload))
		t.FailNow()
	}
	if len(matches) != 1 {
		t.Error("expected length = 1, received length =", strconv.Itoa(len(matches)))
		t.FailNow()
	}

}

// =====================================================================================================================
// FindPairedMatchEmpty : 1 pair and paired matches => empty match
// then check getMatch then check getListMatches
// =====================================================================================================================
func TestKidnerCC_FindPairedMatchEmpty(t *testing.T) {
	scc := new(KidnerCC)
	stub := shim.NewMockStub("kidner", scc)

	// create 1st Pair D1,R1
	res1 := stub.MockInvoke("1", [][]byte{[]byte("createPair"), []byte(AGE1), []byte(BLOODTYPE1), []byte(MEDICALURGENCY1), []byte(HLAs1), []byte(PRA1), []byte(ELIGIBLE1), []byte(TYPE1),
		[]byte(AGE2), []byte(BLOODTYPE2), []byte(MEDICALURGENCY2), []byte(HLAs2), []byte(PRA2), []byte(ELIGIBLE2), []byte(TYPE2), []byte("DrID"), []byte(SIGR), []byte(SIGD), []byte(DRSIG)})
	pair1ID := res1.Payload
	if res1.Status != shim.OK {
		t.Log("bad status received, expected: 200 received:" + strconv.FormatInt(int64(res1.Status), 10))
		t.Log("response: " + string(res1.Message))
		t.FailNow()
	}

	// Pair1 search for paired match => an empty match is expected
	res := stub.MockInvoke("3", [][]byte{[]byte("findPairedMatch"), []byte(pair1ID)})
	matchID := res.Payload
	if res.Status != shim.OK {
		t.Log("bad status received, expected: 200 received:" + strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: " + string(res.Message))
		t.FailNow()
	}

	// see if the match exists
	res = stub.MockInvoke("4", [][]byte{[]byte("getMatch"), []byte(matchID)})
	if res.Status != shim.OK {
		t.Log("bad status received, expected: 200 received:" + strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: " + string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		t.Log("match not found")
		t.Log("response : " + string(res.Message))
		t.FailNow()
	}

	var match Match
	err := json.Unmarshal(res.Payload, &match)
	if err != nil {
		t.Log(errJsonUnmarshall)
		t.Log(err.Error())
		t.FailNow()
	}
	if match.Approved {
		t.Log("match should not be approved")
		t.FailNow()
	}

	if match.DocType != docTypeMatch {
		t.Log("incorrect docType")
		t.FailNow()
	}

	if match.MatchedPairs == nil || len(match.MatchedPairs) == 0 {
		t.Log("match does not contain any matched pairs")
		t.FailNow()
	}

	if match.MatchedPairs[0][0] != string(pair1ID) || len(match.MatchedPairs[0]) != 1 {
		t.Log("match should contain only one ID")
		t.FailNow()
	}

	res = stub.MockInvoke("5", [][]byte{[]byte("getListMatches")})
	if res.Status != shim.OK {
		t.Log("bad status received, expected: 200 received:" + strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: " + string(res.Message))
		t.FailNow()
	}
	matches := make([]Match, 0)
	err = json.Unmarshal(res.Payload, &matches)
	if err != nil {
		t.Error("listMatches", string(res.Payload))
		t.FailNow()
	}
	if len(matches) != 1 {
		t.Error("expected length = 1, received length =", strconv.Itoa(len(matches)))
		t.FailNow()
	}

}

// =====================================================================================================================
// Find paired match incorrect args
// =====================================================================================================================
func TestKidnerCC_FindPairedMatchIncorrectArgs(t *testing.T) {
	scc := new(KidnerCC)
	stub := shim.NewMockStub("kidner", scc)
	res := stub.MockInvoke("1", [][]byte{[]byte("findPairedMatch")})
	if res.Status != shim.ERROR {
		t.Log("bad status received, expected: 500 received:" + strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: " + string(res.Message))
		t.FailNow()
	}
}

// =====================================================================================================================
// FindPairedMatch : No active pairs => Error
// =====================================================================================================================
func TestKidnerCC_FindPairedMatchNoPairs(t *testing.T) {
	scc := new(KidnerCC)
	stub := shim.NewMockStub("kidner", scc)

	res := stub.MockInvoke("3", [][]byte{[]byte("findPairedMatch"), []byte("InexistentID")})
	if res.Status != shim.ERROR {
		t.Log("bad status received, expected: 500 received:" + strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: " + string(res.Message))
		t.FailNow()
	}
}

// =====================================================================================================================
// FindMatchCycle nominal : 4 pairs forming a cycle
// then check getMatch then check getListMatches
// =====================================================================================================================
func TestKidnerCC_FindMatchCycleNominal(t *testing.T) {
	scc := new(KidnerCC)
	stub := shim.NewMockStub("kidner", scc)

	// create 1st Pair D1,R1
	res1 := stub.MockInvoke("1", [][]byte{[]byte("createPair"), []byte("88"), []byte("o"), []byte("0"), []byte("A1,A2,B1,B2,DR1,DR2,DQ1,DQ2"), []byte("0"), []byte("true"), []byte("donor"),
		[]byte("20"), []byte("b"), []byte("2"), []byte("A1,A2,B1,B2,DR1,DR2,DQ1,DQ2"), []byte("0"), []byte("true"), []byte("recipient"), []byte("DrID"), []byte(SIGR), []byte(SIGD), []byte(DRSIG)})
	pair1ID := res1.Payload
	if res1.Status != shim.OK {
		t.Log("bad status received, expected: 200 received:" + strconv.FormatInt(int64(res1.Status), 10))
		t.Log("response: " + string(res1.Message))
		t.FailNow()
	}

	// create 2nd Pair D2,R2
	res2 := stub.MockInvoke("2", [][]byte{[]byte("createPair"), []byte("10"), []byte("b"), []byte("0"), []byte("A1,A2,B1,B2,DR1,DR2,DQ1,DQ2"), []byte("10"), []byte("true"), []byte("donor"),
		[]byte("18"), []byte("ab"), []byte("0"), []byte("A0,A0,B42,B42,DR42,DR42,DQ42,DQ42"), []byte("12"), []byte("true"), []byte("recipient"), []byte("DrID"), []byte(SIGR), []byte(SIGD), []byte(DRSIG)})
	pair2ID := res2.Payload
	if res2.Status != shim.OK {
		t.Log("bad status received, expected: 200 received:" + strconv.FormatInt(int64(res2.Status), 10))
		t.Log("response: " + string(res2.Message))
		t.FailNow()
	}

	// create 3rd Pair D3,R3
	res3 := stub.MockInvoke("3", [][]byte{[]byte("createPair"), []byte("77"), []byte("a"), []byte("0"), []byte("A11,A22,B13,B24,DR15,DR26,DQ17,DQ28"), []byte("25"), []byte("true"), []byte("donor"),
		[]byte("64"), []byte("b"), []byte("1"), []byte("A1,A2,B1,B2,DR1,DR2,DQ1,DQ2"), []byte("20"), []byte("true"), []byte("recipient"), []byte("DrID"), []byte(SIGR), []byte(SIGD), []byte(DRSIG)})
	pair3ID := res3.Payload
	if res3.Status != shim.OK {
		t.Log("bad status received, expected: 200 received:" + strconv.FormatInt(int64(res3.Status), 10))
		t.Log("response: " + string(res3.Message))
		t.FailNow()
	}

	// create 4th Pair D4,R4
	res4 := stub.MockInvoke("4", [][]byte{[]byte("createPair"), []byte("18"), []byte("ab"), []byte("0"), []byte("A0,A0,B42,B42,DR42,DR42,DQ42,DQ42"), []byte("12"), []byte("true"), []byte("donor"),
		[]byte("77"), []byte("a"), []byte("3"), []byte("A11,A22,B13,B24,DR15,DR26,DQ17,DQ28"), []byte("12"), []byte("true"), []byte("recipient"), []byte("DrID"), []byte(SIGR), []byte(SIGD), []byte(DRSIG)})
	pair4ID := res4.Payload
	if res4.Status != shim.OK {
		t.Log("bad status received, expected: 200 received:" + strconv.FormatInt(int64(res4.Status), 10))
		t.Log("response: " + string(res4.Message))
		t.FailNow()
	}

	// Pair1 search for paired match => a match with pair2 is expected
	res := stub.MockInvoke("5", [][]byte{[]byte("findMatchCycle")})
	matchID := res.Payload
	if res.Status != shim.OK {
		t.Log("bad status received, expected: 200 received:" + strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: " + string(res.Message))
		t.FailNow()
	}

	// see if the match exists
	res = stub.MockInvoke("6", [][]byte{[]byte("getMatch"), []byte(matchID)})
	if res.Status != shim.OK {
		t.Log("bad status received, expected: 200 received:" + strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: " + string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		t.Log("match not found")
		t.Log("response : " + string(res.Message))
		t.FailNow()
	}

	var match Match
	err := json.Unmarshal(res.Payload, &match)
	if err != nil {
		t.Log(errJsonUnmarshall)
		t.Log(err.Error())
		t.FailNow()
	}
	if match.Approved {
		t.Log("match should not be approved")
		t.FailNow()
	}

	if match.DocType != docTypeMatch {
		t.Log("incorrect docType")
		t.FailNow()
	}

	if match.MatchedPairs == nil || len(match.MatchedPairs) == 0 {
		t.Log("match does not contain any matched pairs")
		t.FailNow()
	}

	if match.MatchedPairs[0][0] != string(pair1ID) || match.MatchedPairs[0][1] != string(pair3ID) || match.MatchedPairs[0][2] != string(pair4ID) || match.MatchedPairs[0][3] != string(pair2ID) {
		t.Log("matched pairs IDs are incorrect")
		t.FailNow()
	}

	res = stub.MockInvoke("5", [][]byte{[]byte("getListMatches")})
	if res.Status != shim.OK {
		t.Log("bad status received, expected: 200 received:" + strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: " + string(res.Message))
		t.FailNow()
	}
	matches := make([]Match, 0)
	err = json.Unmarshal(res.Payload, &matches)
	if err != nil {
		t.Error("listMatches", string(res.Payload))
		t.FailNow()
	}
	if len(matches) != 1 {
		t.Error("expected length = 1, received length =", strconv.Itoa(len(matches)))
		t.FailNow()
	}

}

// =====================================================================================================================
// Find match cycle incorrect args
// =====================================================================================================================
func TestKidnerCC_FindMatchCycleIncorrectArgs(t *testing.T) {
	scc := new(KidnerCC)
	stub := shim.NewMockStub("kidner", scc)
	res := stub.MockInvoke("1", [][]byte{[]byte("findMatchCycle"), []byte("arg")})
	if res.Status != shim.ERROR {
		t.Log("bad status received, expected: 500 received:" + strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: " + string(res.Message))
		t.FailNow()
	}
}

// =====================================================================================================================
// FindMatchCycle : No active pairs => empty match
// =====================================================================================================================
func TestKidnerCC_FindMatchCycleNoPairs(t *testing.T) {
	scc := new(KidnerCC)
	stub := shim.NewMockStub("kidner", scc)

	res := stub.MockInvoke("1", [][]byte{[]byte("findMatchCycle")})
	if res.Payload == nil {
		t.Log("Payload is nil " + string(res.Message))
		t.FailNow()
	}
	matchID := res.Payload

	// see if the match exists
	res = stub.MockInvoke("2", [][]byte{[]byte("getMatch"), []byte(matchID)})
	if res.Status != shim.OK {
		t.Log("bad status received, expected: 200 received:" + strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: " + string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		t.Log("match not found")
		t.Log("response : " + string(res.Message))
		t.FailNow()
	}

	var match Match
	err := json.Unmarshal(res.Payload, &match)
	if err != nil {
		t.Log(errJsonUnmarshall)
		t.Log(err.Error())
		t.FailNow()
	}
	if match.Approved {
		t.Log("match should not be approved")
		t.FailNow()
	}

	if match.DocType != docTypeMatch {
		t.Log("incorrect docType")
		t.FailNow()
	}

	if len(match.MatchedPairs) != 0 {
		t.Log("match should not contain any matched pairs!")
		t.FailNow()
	}
}

// =====================================================================================================================
// Get inexistent match : error
// =====================================================================================================================
func TestKidnerCC_GetInexistentMatch(t *testing.T) {
	scc := new(KidnerCC)
	stub := shim.NewMockStub("kidner", scc)

	res := stub.MockInvoke("4", [][]byte{[]byte("getMatch"), []byte("inexistentMatchID")})
	if res.Status != shim.ERROR {
		t.Log("bad status received, expected: 500 received:" + strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: " + string(res.Message))
		t.FailNow()
	}
	if res.Payload != nil {
		t.Log("match found : it should NOT exist!")
		t.Log("response : " + string(res.Message))
		t.FailNow()
	}
}

// =====================================================================================================================
// Get match : incorrect args
// =====================================================================================================================
func TestKidnerCC_GetMatchIncorrectArgs(t *testing.T) {
	scc := new(KidnerCC)
	stub := shim.NewMockStub("kidner", scc)

	res := stub.MockInvoke("4", [][]byte{[]byte("getMatch")})
	if res.Status != shim.ERROR {
		t.Log("bad status received, expected: 500 received:" + strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: " + string(res.Message))
		t.FailNow()
	}
	if res.Payload != nil {
		t.Log("match found : it should NOT exist!")
		t.Log("response : " + string(res.Message))
		t.FailNow()
	}
}

// =====================================================================================================================
// Get list matches : incorrect args
// =====================================================================================================================
func TestKidnerCC_GetListMatchesIncorrectArgs(t *testing.T) {
	scc := new(KidnerCC)
	stub := shim.NewMockStub("kidner", scc)

	res := stub.MockInvoke("4", [][]byte{[]byte("getListMatches"), []byte("incorrectarg")})
	if res.Status != shim.ERROR {
		t.Log("bad status received, expected: 500 received:" + strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: " + string(res.Message))
		t.FailNow()
	}
	if res.Payload != nil {
		t.Log("match found : it should NOT exist!")
		t.Log("response : " + string(res.Message))
		t.FailNow()
	}
}

// =====================================================================================================================
// Create a doctor (nominal case)
// =====================================================================================================================
func TestKidnerCC_CreateDoctorNominal(t *testing.T) {
	scc := new(KidnerCC)
	stub := shim.NewMockStub("kidner", scc)
	res := stub.MockInvoke("1", [][]byte{[]byte("createDoctor"), []byte("signature")})
	if res.Status != shim.OK {
		t.Log("bad status received, expected: 200; received:" + strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: " + string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		t.Log("createDoctor failed to create a doctor and get the doctor id")
		t.FailNow()
	}

	dr, err := getDoctorByID(stub, string(res.Payload))
	if err != nil {
		t.Log("error when get doctor" + err.Error())
		t.FailNow()
	}
	if dr.ID != string(res.Payload) {
		t.Log("incorrect doctor ID after create then get Dr")
		t.FailNow()
	}
}

// =====================================================================================================================
// Get a doctor (nominal case)
// =====================================================================================================================
func TestKidnerCC_GetDoctorNominal(t *testing.T) {
	scc := new(KidnerCC)
	stub := shim.NewMockStub("kidner", scc)
	res := stub.MockInvoke("1", [][]byte{[]byte("createDoctor"), []byte("signature")})
	if res.Status != shim.OK {
		t.Log("bad status received, expected: 200; received:" + strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: " + string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		t.Log("createDoctor failed to create a doctor and get the doctor id")
		t.FailNow()
	}

	drID := res.Payload

	res = stub.MockInvoke("2", [][]byte{[]byte("getDoctor"), drID})
	if res.Status != shim.OK {
		t.Log("bad status received, expected: 200; received:" + strconv.FormatInt(int64(res.Status), 10))
		t.FailNow()
	}
	if res.Payload == nil {
		t.Log("getDoctor : result should not be nil")
		t.FailNow()
	}

	var doctor Doctor
	err := json.Unmarshal(res.Payload, &doctor)
	if err != nil {
		t.Log(errJsonUnmarshall + "Doctor entity")
		t.Fail()
	}
}

// =====================================================================================================================
// Get a doctor inexistent doctor
// =====================================================================================================================
func TestKidnerCC_GetDoctorInexistent(t *testing.T) {
	scc := new(KidnerCC)
	stub := shim.NewMockStub("kidner", scc)

	res := stub.MockInvoke("2", [][]byte{[]byte("getDoctor"), []byte("inexistentID")})
	if res.Status != shim.ERROR {
		t.Log("bad status received, expected: 500; received:" + strconv.FormatInt(int64(res.Status), 10))
		t.FailNow()
	}
	if res.Payload != nil {
		t.Log("getDoctor : result should be nil")
		t.FailNow()
	}
}

// =====================================================================================================================
// Create a doctor (incorrect args)
// =====================================================================================================================
func TestKidnerCC_CreateDoctorIncorrectArgs(t *testing.T) {
	scc := new(KidnerCC)
	stub := shim.NewMockStub("kidner", scc)

	// test create doctor
	res := stub.MockInvoke("8", [][]byte{[]byte("createDoctor")})
	if res.Status != shim.ERROR {
		t.Log("bad status received, expected: 500; received:" + strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: " + string(res.Message))
		t.FailNow()
	}

	if res.Payload != nil {
		t.Log("createDoctor should fail")
		t.FailNow()
	}
}

// =====================================================================================================================
// Approve match (nominal case)
// =====================================================================================================================
func TestKidnerCC_ApproveMatchNominal(t *testing.T) {
	scc := new(KidnerCC)
	stub := shim.NewMockStub("kidner", scc)

	// create Doctor
	// res := stub.MockInvoke("0", [][]byte{[]byte("createDoctor"), []byte("signature")})
	// if res.Status != shim.OK {
	// 	t.Log("bad status received, expected: 200; received:" + strconv.FormatInt(int64(res.Status), 10))
	// 	t.Log("response: " + string(res.Message))
	// 	t.FailNow()
	// }
	// if res.Payload == nil {
	// 	t.Log("createDoctor failed to create a doctor and get the doctor id")
	// 	t.FailNow()
	// }
	DrID := "DrID"

	// create 1st Pair D1,R1
	res1 := stub.MockInvoke("1", [][]byte{[]byte("createPair"), []byte(AGE1), []byte(BLOODTYPE1), []byte(MEDICALURGENCY1), []byte(HLAs1), []byte(PRA1), []byte(ELIGIBLE1), []byte(TYPE1),
		[]byte(AGE2), []byte(BLOODTYPE2), []byte(MEDICALURGENCY2), []byte(HLAs2), []byte(PRA2), []byte(ELIGIBLE2), []byte(TYPE2), []byte("DrID"), []byte(SIGR), []byte(SIGD), []byte(DRSIG)})
	pair1ID := res1.Payload
	if res1.Status != shim.OK {
		t.Log("bad status received, expected: 200 received:" + strconv.FormatInt(int64(res1.Status), 10))
		t.Log("response: " + string(res1.Message))
		t.FailNow()
	}

	// create 2nd Pair D2,R2
	res2 := stub.MockInvoke("2", [][]byte{[]byte("createPair"), []byte(AGE2), []byte(BLOODTYPE2), []byte(MEDICALURGENCY2), []byte(HLAs2), []byte(PRA2), []byte(ELIGIBLE2), []byte(TYPE1),
		[]byte(AGE1), []byte(BLOODTYPE1), []byte(MEDICALURGENCY1), []byte(HLAs1), []byte(PRA1), []byte(ELIGIBLE1), []byte(TYPE2), []byte("DrID"), []byte(SIGR), []byte(SIGD), []byte(DRSIG)})
	//pair2ID := res2.Payload
	if res2.Status != shim.OK {
		t.Log("bad status received, expected: 200 received:" + strconv.FormatInt(int64(res2.Status), 10))
		t.Log("response: " + string(res2.Message))
		t.FailNow()
	}

	// Pair1 searches for paired match => a match with pair2 is expected
	res := stub.MockInvoke("3", [][]byte{[]byte("findPairedMatch"), []byte(pair1ID)})
	matchID := res.Payload
	if res.Status != shim.OK {
		t.Log("bad status received, expected: 200 received:" + strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: " + string(res.Message))
		t.FailNow()
	}

	// approve match bad signature
	// res = stub.MockInvoke("7", [][]byte{[]byte("approveMatch"), []byte(DrID), []byte("badsignature"), []byte(matchID)})
	// if res.Status != shim.ERROR {
	// 	t.Log("bad status received, expected: 500 received:" + strconv.FormatInt(int64(res.Status), 10))
	// 	t.Log("response: " + string(res.Message))
	// 	t.FailNow()
	// }

	// doctor approves match : DrID and his signature needed
	res = stub.MockInvoke("4", [][]byte{[]byte("approveMatch"), []byte(DrID), []byte("signature"), []byte(matchID)})
	matchID = res.Payload

	// get Match
	res = stub.MockInvoke("5", [][]byte{[]byte("getMatch"), []byte(matchID)})
	if res.Status != shim.OK {
		t.Log("bad status received, expected: 200 received:" + strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: " + string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		t.Log("match not found")
		t.Log("response : " + string(res.Message))
		t.FailNow()
	}

	var match Match
	err := json.Unmarshal(res.Payload, &match)
	if err != nil {
		t.Log(errJsonUnmarshall)
		t.Log(err.Error())
		t.FailNow()
	}

	if !match.Approved {
		t.Error("match was not approved successfully")
		t.FailNow()
	}

	if match.EndorcingDr != string(DrID) {
		t.Error("incorrect endorcing Dr")
		t.FailNow()
	}

	if match.DrSig != "signature" {
		t.Error("incorrect Dr signature")
		t.FailNow()
	}

}

// =====================================================================================================================
// Approve match : Incorrect number of args
// =====================================================================================================================
func TestKidnerCC_ApproveMatchIncorrectArgs(t *testing.T) {
	scc := new(KidnerCC)
	stub := shim.NewMockStub("kidner", scc)

	// approve match wrong number of arguments
	res := stub.MockInvoke("1", [][]byte{[]byte("approveMatch"), []byte("DrID"), []byte("signature"), []byte("signature"), []byte("matchID")})
	if res.Status != shim.ERROR {
		t.Log("bad status received, expected: 500 received:" + strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: " + string(res.Message))
		t.FailNow()
	}
}

// =====================================================================================================================
// Refuse match (nominal case)
// =====================================================================================================================
func TestKidnerCC_RefuseMatchNominal(t *testing.T) {
	scc := new(KidnerCC)
	stub := shim.NewMockStub("kidner", scc)

	// create Doctor
	// res := stub.MockInvoke("0", [][]byte{[]byte("createDoctor"), []byte("signature")})
	// if res.Status != shim.OK {
	// 	t.Log("bad status received, expected: 200; received:" + strconv.FormatInt(int64(res.Status), 10))
	// 	t.Log("response: " + string(res.Message))
	// 	t.FailNow()
	// }
	// if res.Payload == nil {
	// 	t.Log("createDoctor failed to create a doctor and get the doctor id")
	// 	t.FailNow()
	// }
	DrID := "DrID"

	// create 1st Pair D1,R1
	res1 := stub.MockInvoke("1", [][]byte{[]byte("createPair"), []byte(AGE1), []byte(BLOODTYPE1), []byte(MEDICALURGENCY1), []byte(HLAs1), []byte(PRA1), []byte(ELIGIBLE1), []byte(TYPE1),
		[]byte(AGE2), []byte(BLOODTYPE2), []byte(MEDICALURGENCY2), []byte(HLAs2), []byte(PRA2), []byte(ELIGIBLE2), []byte(TYPE2), []byte("DrID"), []byte(SIGR), []byte(SIGD), []byte(DRSIG)})
	pair1ID := res1.Payload
	if res1.Status != shim.OK {
		t.Log("bad status received, expected: 200 received:" + strconv.FormatInt(int64(res1.Status), 10))
		t.Log("response: " + string(res1.Message))
		t.FailNow()
	}

	// create 2nd Pair D2,R2
	res2 := stub.MockInvoke("2", [][]byte{[]byte("createPair"), []byte(AGE2), []byte(BLOODTYPE2), []byte(MEDICALURGENCY2), []byte(HLAs2), []byte(PRA2), []byte(ELIGIBLE2), []byte(TYPE1),
		[]byte(AGE1), []byte(BLOODTYPE1), []byte(MEDICALURGENCY1), []byte(HLAs1), []byte(PRA1), []byte(ELIGIBLE1), []byte(TYPE2), []byte("DrID"), []byte(SIGR), []byte(SIGD), []byte(DRSIG)})
	//pair2ID := res2.Payload
	if res2.Status != shim.OK {
		t.Log("bad status received, expected: 200 received:" + strconv.FormatInt(int64(res2.Status), 10))
		t.Log("response: " + string(res2.Message))
		t.FailNow()
	}

	// Pair1 searches for paired match => a match with pair2 is expected
	res := stub.MockInvoke("3", [][]byte{[]byte("findPairedMatch"), []byte(pair1ID)})
	matchID := res.Payload
	if res.Status != shim.OK {
		t.Log("bad status received, expected: 200 received:" + strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: " + string(res.Message))
		t.FailNow()
	}

	// refuse match bad signature
	// res = stub.MockInvoke("7", [][]byte{[]byte("refuseMatch"), []byte(DrID), []byte("badsignature"), []byte(matchID)})
	// if res.Status != shim.ERROR {
	// 	t.Log("bad status received, expected: 500 received:" + strconv.FormatInt(int64(res.Status), 10))
	// 	t.Log("response: " + string(res.Message))
	// 	t.FailNow()
	// }

	// doctor refuses match
	res = stub.MockInvoke("4", [][]byte{[]byte("refuseMatch"), []byte(DrID), []byte("signature"), []byte(matchID)})
	matchID = res.Payload

	// get Match
	res = stub.MockInvoke("5", [][]byte{[]byte("getMatch"), []byte(matchID)})
	if res.Status != shim.OK {
		t.Log("bad status received, expected: 200 received:" + strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: " + string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		t.Log("match not found")
		t.Log("response : " + string(res.Message))
		t.FailNow()
	}

	var match Match
	err := json.Unmarshal(res.Payload, &match)
	if err != nil {
		t.Log(errJsonUnmarshall)
		t.Log(err.Error())
		t.FailNow()
	}

	if match.Approved {
		t.Error("match was not refused successfully")
		t.FailNow()
	}

	if match.EndorcingDr != string(DrID) {
		t.Error("incorrect endorcing Dr")
		t.FailNow()
	}
}

// =====================================================================================================================
// Refuse match : Incorrect args
// =====================================================================================================================
func TestKidnerCC_RefuseMatchIncorrectArgs(t *testing.T) {
	scc := new(KidnerCC)
	stub := shim.NewMockStub("kidner", scc)

	// refuse match wrong number of arguments
	res := stub.MockInvoke("1", [][]byte{[]byte("refuseMatch"), []byte("DrID"), []byte("signature"), []byte("signature"), []byte("matchID")})
	if res.Status != shim.ERROR {
		t.Log("bad status received, expected: 500 received:" + strconv.FormatInt(int64(res.Status), 10))
		t.Log("response: " + string(res.Message))
		t.FailNow()
	}
}

// =====================================================================================================================
// Main()
// =====================================================================================================================
func TestMain(m *testing.M) {
	InitLogger()
	code := m.Run()
	os.Exit(code)
}
