// =====================================================================================================================
// This file contains methods that retrive data from the ledger
// =====================================================================================================================
package main

import (
	"encoding/json"
	"errors"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// =====================================================================================================================
// getPairByID : get the pair object by ID
// =====================================================================================================================
func getPairByID(stub shim.ChaincodeStubInterface, ID string) (*Pair, error) {
	valAsbytes, err := stub.GetState(ID)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	} else if valAsbytes == nil {
		return nil, errors.New("pair does not exist")
	}

	var pair Pair
	err = json.Unmarshal(valAsbytes, &pair)
	if err != nil {
		logger.Error(err.Error())
		return nil, errors.New(errJsonUnmarshall)
	}

	return &pair, nil
}

// =====================================================================================================================
// listActivePairs : Get the list of active Pairs
// =====================================================================================================================
func getListActivePairs(stub shim.ChaincodeStubInterface) ([]Pair, error) {
	allPairs, err := getPairsByIndex(stub, indexPair, []string{docTypePair})
	if err != nil {
		logger.Error(buildError("Error get list of active pairs"))
		return nil, err
	} else if allPairs == nil {
		logger.Error("empty list")
		return nil, errors.New("empty list")
	}

	res := make([]Pair, 0)
	for _, v := range allPairs {
		if v.Active {
			res = append(res, v)
		}
	}

	return res, nil
}

// =====================================================================================================================
// use index to retrieve a list of pairs
// =====================================================================================================================
func getPairsByIndex(stub shim.ChaincodeStubInterface, index string, keys []string) ([]Pair, error) {
	logger.Debug("getPairsByIndex(index=" + index + ") : calling method -")
	resultsIterator, err := stub.GetStateByPartialCompositeKey(index, keys)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	defer resultsIterator.Close()

	pairs := make([]Pair, 0)
	for i := 0; resultsIterator.HasNext(); i++ {
		KV, err := resultsIterator.Next()
		if err != nil {
			logger.Error(err)
			return nil, err
		}
		_, compositeKeyParts, err := stub.SplitCompositeKey(KV.Key)
		if err != nil {
			logger.Error(err)
			return nil, err
		}
		pairID := compositeKeyParts[len(compositeKeyParts)-1]
		pairAsBytes, err := stub.GetState(pairID)
		if err != nil {
			logger.Error(err)
			return nil, err
		}
		pair := Pair{}
		err = json.Unmarshal(pairAsBytes, &pair)
		if err == nil {
			pairs = append(pairs, pair)
		}
	}
	return pairs, nil
}

// =====================================================================================================================
// getMatchByID : get the match object by ID
// =====================================================================================================================
func getMatchByID(stub shim.ChaincodeStubInterface, ID string) (*Match, error) {
	valAsbytes, err := stub.GetState(ID)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	} else if valAsbytes == nil {
		return nil, errors.New("match does not exist")
	}

	var match Match
	err = json.Unmarshal(valAsbytes, &match)
	if err != nil {
		logger.Error(err.Error())
		return nil, errors.New(errJsonUnmarshall)
	}

	return &match, nil
}

// =====================================================================================================================
// use index to retrieve a list of matches
// =====================================================================================================================
func getMatchesByIndex(stub shim.ChaincodeStubInterface, index string, keys []string) ([]Match, error) {
	logger.Debug("getMatchesByIndex(index=" + index + ") : calling method -")
	resultsIterator, err := stub.GetStateByPartialCompositeKey(index, keys)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	defer resultsIterator.Close()

	matches := make([]Match, 0)
	for i := 0; resultsIterator.HasNext(); i++ {
		KV, err := resultsIterator.Next()
		if err != nil {
			logger.Error(err)
			return nil, err
		}
		_, compositeKeyParts, err := stub.SplitCompositeKey(KV.Key)
		if err != nil {
			logger.Error(err)
			return nil, err
		}
		matchID := compositeKeyParts[len(compositeKeyParts)-1]
		matchAsBytes, err := stub.GetState(matchID)
		if err != nil {
			logger.Error(err)
			return nil, err
		}
		match := Match{}
		err = json.Unmarshal(matchAsBytes, &match)
		if err == nil {
			matches = append(matches, match)
		}
	}
	return matches, nil
}

// =====================================================================================================================
// getDoctorByID : get the doctor object by ID
// =====================================================================================================================
func getDoctorByID(stub shim.ChaincodeStubInterface, drID string) (*Doctor, error) {
	valAsbytes, err := stub.GetState(drID)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	} else if valAsbytes == nil {
		return nil, errors.New("doctor does not exist")
	}

	var doctor Doctor
	err = json.Unmarshal(valAsbytes, &doctor)
	if err != nil {
		logger.Error(err.Error())
		return nil, errors.New(errJsonUnmarshall)
	}

	return &doctor, nil
}

// ============================================================================================================================
// createIndex - Create all index for a pair
// ============================================================================================================================
func createPairIndex(stub shim.ChaincodeStubInterface, pair Pair) error {
	logger.Debug("createIndex(pairID=" + pair.ID + ") : calling method -")
	indexActiveKey, err := stub.CreateCompositeKey(indexPair, []string{pair.DocType, pair.ID})
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	stub.PutState(indexActiveKey, []byte{0x00})
	return nil
}

// ============================================================================================================================
// deleteIndex - delete all index for a pair
// ============================================================================================================================
func deletePairIndex(stub shim.ChaincodeStubInterface, pair Pair) error {
	logger.Debug("deleteIndex(pairID=" + pair.ID + ") : calling method -")
	indexActiveKey, err := stub.CreateCompositeKey(indexPair, []string{pair.DocType, pair.ID})
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	stub.DelState(indexActiveKey)

	return nil
}

// ============================================================================================================================
// createIndex - Create all index for a match || match is never deleted so no need to delete Index
// ============================================================================================================================
func createMatchIndex(stub shim.ChaincodeStubInterface, match Match) error {
	logger.Debug("createIndex(MatchID=" + match.ID + ") : calling method -")
	indexActiveKey, err := stub.CreateCompositeKey(indexMatch, []string{match.DocType, match.ID})
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	stub.PutState(indexActiveKey, []byte{0x00})
	return nil
}
