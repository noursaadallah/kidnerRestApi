package blockchain

// DeactivatePair : method will be called from the DeactivatePairHandler
func (setup *FabricSetup) DeactivatePair(pairID string) (string, error) {

	// Prepare args : function and parameters
	var args []string
	args = append(args, "deactivatePair")
	args = append(args, pairID)

	return setup.invoke(args)
}
