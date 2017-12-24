package blockchain

// DeletePair : method will be called from the DeletePairHandler
func (setup *FabricSetup) DeletePair(pairID string) (string, error) {

	// Prepare args : function and parameters
	var args []string
	args = append(args, "deletePair")
	args = append(args, pairID)

	return setup.invoke(args)
}
