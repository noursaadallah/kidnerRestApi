package blockchain

// UpdatePair : method will be called from the UpdatePairHandler
func (setup *FabricSetup) UpdatePair(param []string) (string, error) {

	// Prepare args : function and parameters
	var args []string
	args = append(args, "updatePair")
	args = append(args, param...)

	return setup.invoke(args)
}
