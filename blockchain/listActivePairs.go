package blockchain

// ListActivePairs : get the list of acive pairs
func (setup *FabricSetup) ListActivePairs() ([]byte, error) {

	// Prepare arguments : function and parameters
	var args []string
	args = append(args, "listActivePairs")

	return setup.query(args)
}
