package blockchain

import (
	"fmt"

	api "github.com/hyperledger/fabric-sdk-go/api"
	fcutil "github.com/hyperledger/fabric-sdk-go/pkg/util"
)

// GetPair : get pair information by ID
func (setup *FabricSetup) GetPair(ID string) ([]byte, error) {

	// Prepare arguments : function and parameters
	var args []string
	args = append(args, "getPair")
	args = append(args, ID)

	// Make the proposal and submit it to the network (via our primary peer)
	transactionProposalResponses, _, err := fcutil.CreateAndSendTransactionProposal(
		setup.Channel,
		setup.ChaincodeId,
		setup.ChannelId,
		args,
		[]api.Peer{setup.Channel.GetPrimaryPeer()}, // Peer contacted when submitted the proposal
		nil,
	)
	if err != nil {
		return []byte(""), fmt.Errorf("Create and send transaction proposal return error in the query hello: %v", err)
	}
	return transactionProposalResponses[0].ProposalResponse.GetResponse().Payload, nil
}
