package blockchain

// approve a match
// param : DrID , DrSig , MatchID
func (setup *FabricSetup) ApproveMatch(param []string) (string, error) {

	// Prepare arguments = function + parameters
	var args []string
	args = append(args, "approveMatch")
	args = append(args, param...)

	return setup.invoke(args)
}
