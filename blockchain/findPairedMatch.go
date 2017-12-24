package blockchain

// FindPairedMatch : method will be called from the FindPairedMatchHandler
func (setup *FabricSetup) FindPairedMatch(pairID string) (string, error) {

	// Prepare args : function and parameters
	var args []string
	args = append(args, "findPairedMatch")
	args = append(args, pairID)

	return setup.invoke(args)
}
