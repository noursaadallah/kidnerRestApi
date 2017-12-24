package blockchain

// CreatePair : method will be called from the CreatePairHandler
func (setup *FabricSetup) CreatePair(param []string) (string, error) {

	// Prepare args : function and parameters
	var args []string
	args = append(args, "createPair")
	args = append(args, param...)

	return setup.invoke(args)
}
