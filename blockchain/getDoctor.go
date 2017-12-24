package blockchain

// GetDoctor : get dr information by ID
func (setup *FabricSetup) GetDoctor(ID string) ([]byte, error) {

	// Prepare arguments : function and parameters
	var args []string
	args = append(args, "getDoctor")
	args = append(args, ID)

	return setup.query(args)
}
