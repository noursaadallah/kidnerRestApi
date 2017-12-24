package blockchain

// ListMatches : get the list of matches in the blockchain
func (setup *FabricSetup) ListMatches() ([]byte, error) {

	// Prepare arguments : function (no parameters)
	var args []string
	args = append(args, "getListMatches")

	return setup.query(args)
}
