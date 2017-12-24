package blockchain

import (
	"fmt"
	"time"

	api "github.com/hyperledger/fabric-sdk-go/api"
	fcutil "github.com/hyperledger/fabric-sdk-go/pkg/util"
)

// DeactivatePair : method will be called from the DeactivatePairHandler
func (setup *FabricSetup) DeactivatePair(pairID string) (string, error) {

	// Prepare args : function and parameters
	var args []string
	args = append(args, "deactivatePair")
	args = append(args, pairID)

	// Add data that will be visible in the proposal, like a description of the invoke request
	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("Transient data in invoke")

	// Make a nex transaction proposal and send it
	transactionProposalResponse, txID, err := fcutil.CreateAndSendTransactionProposal(
		setup.Channel,
		setup.ChaincodeId,
		setup.ChannelId,
		args,
		[]api.Peer{setup.Channel.GetPrimaryPeer()},
		transientDataMap,
	)
	if err != nil {
		return "", fmt.Errorf("Create and send transaction proposal in the invoke hello return error: %v", err)
	}

	// Register the Fabric SDK to listen to the event that will come back when the transaction will be send
	done, fail := fcutil.RegisterTxEvent(txID, setup.EventHub)

	// Send the final transaction signed by endorser
	if _, err := fcutil.CreateAndSendTransaction(setup.Channel, transactionProposalResponse); err != nil {
		return "", fmt.Errorf("Create and send transaction in the invoke hello return error: %v", err)
	}

	// Wait the result of the submission
	select {
	// Transaction Ok
	case <-done:
		return txID, nil

	// Transaction failed
	case <-fail:
		return "", fmt.Errorf("Error received from eventhub for txid(%s) error(%v)", txID, fail)

	// Transaction timeout
	case <-time.After(time.Second * 30):
		return "", fmt.Errorf("Didn't receive block event for txid(%s)", txID)
	}
}
