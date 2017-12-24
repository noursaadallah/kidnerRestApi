package blockchain

// GetPair : get pair information by ID
func (setup *FabricSetup) GetPair(ID string) ([]byte, error) {

	// Prepare arguments : function and parameters
	var args []string
	args = append(args, "getPair")
	args = append(args, ID)

	return setup.query(args)
}
