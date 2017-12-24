package blockchain

// FindMatchCycle : method will be called from the FindMatchCycleHandler
// TODO make automatic (call from CreatePairHandler)
func (setup *FabricSetup) FindMatchCycle() (string, error) {

	// Prepare args : function and parameters
	var args []string
	args = append(args, "findMatchCycle")

	return setup.invoke(args)
}
